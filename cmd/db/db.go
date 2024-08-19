package db

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lib/pq"
)

type Storage interface {
	Registeruser(user User) error
	Loginuser(user User) error
	UpOrderuser(name string, numorder int) error
}

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type DB struct {
	db *sql.DB
}

var (
	ErrDuplicateUser      = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid password or login")
	ErrDuplicateOrder     = errors.New("order belongs to another user")
	ErrAlreadyUpload      = errors.New("order upload before")
)

func NewDataBase(dbname string) (Storage, error) {
	db, err := sql.Open("postgres", dbname)
	if err != nil {
		fmt.Println("1", err)
		return nil, err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		fmt.Println("444", err)
		fmt.Println("2", err)
		return nil, err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./../migrations",
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

func (db *DB) Registeruser(user User) error {
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

func (db *DB) Loginuser(user User) error {
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

func (db *DB) UpOrderuser(name string, numorder int) error {
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
