package db

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func NewDataBase(dbname string) (*sql.DB, error) {
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
		"x2/cmd/migrations",
		"postgres", driver)
	if err != nil {
		fmt.Println("3", err)
		return nil, err
	}
	err = m.Up()
	if err != nil {
		fmt.Println("5", err)
		return nil, err
	}

	return db, nil
}
