package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"x2/cmd/config"
	"x2/cmd/db"
	"x2/cmd/middleware"
)

func Init() start {
	url, addressURL, dbname := config.Flags()
	//dbname = "postgresql://postgres:190603@localhost:5432/postgres?sslmode=disable"
	fmt.Println(url)
	db, err1 := db.NewDataBase(dbname)
	if err1 != nil {
		panic(err1)
	}
	return start{URL: url, addressBonus: addressURL, database: db}
}

type start struct {
	URL          string
	addressBonus string
	database     db.Storage
}

func (st start) Register(res http.ResponseWriter, req *http.Request) {
	var user db.User
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	err := st.database.RegisterUser(user)
	if err != nil {
		if errors.Is(err, db.ErrDuplicateUser) {
			res.WriteHeader(http.StatusConflict)
			return
		}
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = middleware.NewCookie(res, user.Login)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func (st start) Login(res http.ResponseWriter, req *http.Request) {
	var user db.User
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	err := st.database.LoginUser(user)
	if err != nil {
		if errors.Is(err, db.ErrInvalidCredentials) {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = middleware.NewCookie(res, user.Login)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func (st start) UpOrder(res http.ResponseWriter, req *http.Request) {
	name := req.Context().Value(middleware.ContextKey("Name")).(middleware.ToHand)
	if !name.IsAuth {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	numorderbyte, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "unable to read body", http.StatusBadRequest)
		return
	}
	numorder, err1 := strconv.Atoi(string(numorderbyte))
	if err1 != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	err = st.database.UpOrderUser(name.Value, numorder)
	if err != nil {
		if errors.Is(err, db.ErrDuplicateOrder) {
			res.WriteHeader(http.StatusConflict)
			return
		}
		if errors.Is(err, db.ErrAlreadyUpload) {
			res.WriteHeader(http.StatusOK)
			return
		}
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusAccepted)
}

func (st start) GetOrder(res http.ResponseWriter, req *http.Request) {
	name := req.Context().Value(middleware.ContextKey("Name")).(middleware.ToHand)
	if !name.IsAuth {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
}

func (st start) Balance(res http.ResponseWriter, req *http.Request) {
	name := req.Context().Value(middleware.ContextKey("Name")).(middleware.ToHand)
	if !name.IsAuth {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
}

func (st start) LossBonus(res http.ResponseWriter, req *http.Request) {
	name := req.Context().Value(middleware.ContextKey("Name")).(middleware.ToHand)
	if !name.IsAuth {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
}

func (st start) Info(res http.ResponseWriter, req *http.Request) {
	name := req.Context().Value(middleware.ContextKey("Name")).(middleware.ToHand)
	if !name.IsAuth {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
}

func (st start) InfoBonus(res http.ResponseWriter, req *http.Request) {
	name := req.Context().Value(middleware.ContextKey("Name")).(middleware.ToHand)
	if !name.IsAuth {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
}
