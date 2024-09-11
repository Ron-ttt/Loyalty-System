package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"x2/internal/config"
	"x2/internal/db"
	"x2/internal/middleware"
)

func Init() start {
	cfg := config.NewConfig()
	//dbb := "postgresql://postgres:190603@localhost:5432/postgres?sslmode=disable"
	//cfg.DBAddress = &dbb
	fmt.Println(cfg.ServerAddress)
	db, err1 := db.NewDataBase(*cfg.DBAddress)
	if err1 != nil {
		panic(err1)
	}
	return start{URL: cfg.ServerAddress, addressBonus: *cfg.AccrualAddress, database: db}
}
func (st start) Bonus() error {
	var order db.Accrual
	list := st.database.Numorder()
	for _, num := range list {
		address := st.addressBonus + "/api/orders/" + num
		resp, err := http.Get(address)
		if err != nil {
			log.Println("запрос кудато не ушел все наебнулось")
			return err
		}
		if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
			log.Println("джейсон выебывается")
			return err
		}
		st.database.Updateorderdata(order)
	}
	return nil
}

type start struct {
	URL          string
	addressBonus string
	database     db.Storage
}

// регистрация пользователя
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

// аутентификация пользователя
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

// загрузка пользователем номера заказа для расчёта
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
	if !checkLuhn(string(numorderbyte)) {
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

func checkLuhn(number string) bool {
	sum := 0
	alt := false
	// Итерируемся по цифрам числа с конца к началу
	for i := len(number) - 1; i >= 0; i-- {
		digit, err := strconv.Atoi(string(number[i]))
		if err != nil {
			return false // Возврат false в случае недопустимого символа
		}
		if alt {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		alt = !alt // Меняем флаг альтернативной обработки
	}
	return sum%10 == 0 // Проверяем, делится ли сумма на 10
}

// получение списка загруженных пользователем номеров заказов,
// статусов их обработки и информации о начислениях
func (st start) GetOrder(res http.ResponseWriter, req *http.Request) {
	name := req.Context().Value(middleware.ContextKey("Name")).(middleware.ToHand)
	if !name.IsAuth {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	orders, err := st.database.GetOrderUser(name.Value)
	if err != nil {
		if errors.Is(err, db.ErrNoOrders) {
			http.Error(res, err.Error(), http.StatusNoContent)
			return
		}
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(res).Encode(orders); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}

// получение текущего баланса счёта баллов лояльности пользователя
func (st start) Balance(res http.ResponseWriter, req *http.Request) {
	name := req.Context().Value(middleware.ContextKey("Name")).(middleware.ToHand)
	if !name.IsAuth {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	var bonus db.Account
	bonus, err := st.database.BalanceUser(name.Value)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(res).Encode(bonus); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}

// запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа
func (st start) LossBonus(res http.ResponseWriter, req *http.Request) {
	name := req.Context().Value(middleware.ContextKey("Name")).(middleware.ToHand)
	if !name.IsAuth {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
}

// получение информации о выводе средств с накопительного счёта пользователем
func (st start) Info(res http.ResponseWriter, req *http.Request) {
	name := req.Context().Value(middleware.ContextKey("Name")).(middleware.ToHand)
	if !name.IsAuth {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
}
