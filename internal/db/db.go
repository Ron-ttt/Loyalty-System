package db

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lib/pq"
)

type Storage interface {
	RegisterUser(user User) error
	LoginUser(user User) error
	UpOrderUser(name string, numorder int) error
	GetOrderUser(name string) ([]Orders, error)
	BalanceUser(name string) (Account, error)
	UpdateOrderData(data Accrual) error
	NumOrder() ([]string, error)
}
type Accrual struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual"`
}
type Orders struct {
	Number  string    `json:"number"`
	Status  string    `json:"status"`
	Accrual int       `json:"accrual"`
	Time    time.Time `json:"-"`
	TimeRfc string    `json:"uploaded_at"`
}
type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
type Account struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}
type DB struct {
	db *sql.DB
}

var (
	ErrDuplicateUser      = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid password or login")
	ErrDuplicateOrder     = errors.New("order belongs to another user")
	ErrAlreadyUpload      = errors.New("order upload before")
	ErrNoOrders           = errors.New("order no upload")
)

func NewDataBase(dbname string) (Storage, error) {
	db, err := sql.Open("postgres", dbname)
	if err != nil {
		fmt.Println("1", err)
		return nil, err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithDatabaseInstance(
		//"file://C:/Users/lolim/x2/x2/db/migrations",
		"file://./db/migrations",
		"postgres", driver)
	if err != nil {
		fmt.Println("3", err)
		return nil, err
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		fmt.Println("5", err)
		return nil, err
	}

	return &DB{db: db}, nil
}

func (db *DB) RegisterUser(user User) error {
	hashPassword := md5.Sum([]byte(user.Password))

	_, err := db.db.Exec("INSERT INTO users(login, password)"+" VALUES($1,$2)", user.Login, hex.EncodeToString(hashPassword[:]))
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return ErrDuplicateUser
			}
		}
		return err
	}
	return nil
}

func (db *DB) LoginUser(user User) error {
	rows := db.db.QueryRow("SELECT password FROM users WHERE login = $1", user.Login)
	var password string
	err := rows.Scan(&password)
	if err != nil {
		return err
	}
	hashPassword := md5.Sum([]byte(user.Password))
	if password != hex.EncodeToString(hashPassword[:]) {
		return ErrInvalidCredentials
	}
	return nil
}

func (db *DB) UpOrderUser(name string, numorder int) error {
	rows := db.db.QueryRow("SELECT id FROM users WHERE login = $1", name)
	var id int
	err := rows.Scan(&id)
	if err != nil {
		return err
	}
	_, err = db.db.Exec("INSERT INTO orders(users_id, order_id)"+" VALUES($1,$2)", id, numorder)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				rows := db.db.QueryRow("SELECT users_id FROM orders WHERE order_id = $1", numorder)
				var id2 int
				err1 := rows.Scan(&id2)
				if err1 != nil {
					return err1
				}
				if id == id2 {
					return ErrAlreadyUpload
				}
				return ErrDuplicateOrder
			}
		}
		return err
	}
	return nil
}

func (db *DB) GetOrderUser(name string) ([]Orders, error) {
	row := db.db.QueryRow("SELECT id FROM users WHERE login = $1", name)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return nil, err
	}

	var listorders []Orders
	//////errrrr v zaprose pochitay tz accural nado tozhe polychat'
	rows, err := db.db.Query("SELECT order_id, status, bonus, created_at FROM orders WHERE users_id=$1 order by created_at desc", id)
	if rows.Err() != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	if rows == nil {
		return nil, ErrNoOrders
	}
	for rows.Next() {
		var order Orders
		err := rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.Time)
		if err != nil {
			return nil, err
		}

		order.TimeRfc = order.Time.Format(time.RFC3339)
		listorders = append(listorders, order)
	}
	return listorders, nil
}

func (db *DB) BalanceUser(name string) (Account, error) {
	rows, err := db.db.Query("SELECT bonus FROM orders WHERE users_id=(SELECT id FROM users WHERE login = $1)", name)
	if rows.Err() != nil {
		return Account{}, err
	}
	if err != nil {
		return Account{}, err
	}
	var wd float32 = 0
	var cur float32 = 0
	for rows.Next() {
		// y tebya v bd int a ne float
		var num float32
		err := rows.Scan(&num)
		if err != nil {
			return Account{}, err
		}
		if num < 0 {
			wd = wd - num
		}
		cur = cur + num
	}
	return Account{Current: cur / 100, Withdrawn: wd / 100}, nil
}

func (db *DB) UpdateOrderData(data Accrual) error {
	_, err := db.db.Exec("UPDATE orders SET status=$1,bonus=$2 WHERE order_id = $3", data.Status, 100*data.Accrual, data.Order)
	if err != nil {
		log.Println("бд решила что может творить хуйню", err)
		return err
	}
	return nil
}

func (db *DB) NumOrder() ([]string, error) {
	rows, err := db.db.Query("SELECT order_id FROM orders WHERE status = 'NEW' OR status = 'PROCESSING' OR status = 'REGISTERED'")
	if rows.Err() != nil {
		log.Println("бд решила что может творить хуйню", err)
		return nil, err
	}
	if err != nil {
		log.Println("бд решила что может творить хуйню", err)
		return nil, err
	}
	var list []string
	for rows.Next() {
		var num string
		err := rows.Scan(&num)
		if err != nil {
			log.Println("скан решила что может творить хуйню", err)
			return nil, err
		}
		list = append(list, num)
	}
	return list, nil
}
