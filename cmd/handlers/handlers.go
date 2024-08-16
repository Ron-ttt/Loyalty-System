package handlers

import (
	"database/sql"
	"net/http"
	"x2/cmd/config"
	"x2/cmd/db"
)

func Starts() start {
	_, addressURl, dbname := config.Flags()
	//dbname = "postgresql://postgres:190603@localhost:5432/postgres?sslmode=disable"

	db, err1 := db.NewDataBase(dbname)
	if err1 != nil {
		panic(err1)
	}
	return start{addressBonus: addressURl, database: db}
}

type start struct {
	addressBonus string
	database     *sql.DB
}

func (st start) Register(res http.ResponseWriter, req *http.Request) {

}

func (st start) Login(res http.ResponseWriter, req *http.Request) {

}

func (st start) UpOrder(res http.ResponseWriter, req *http.Request) {

}

func (st start) GetOrder(res http.ResponseWriter, req *http.Request) {

}

func (st start) Balance(res http.ResponseWriter, req *http.Request) {

}

func (st start) LossBonus(res http.ResponseWriter, req *http.Request) {

}

func (st start) Info(res http.ResponseWriter, req *http.Request) {

}

func (st start) InfoBonus(res http.ResponseWriter, req *http.Request) {

}
