package middleware

import (
	"context"
	"errors"
	"net/http"
	"x2/cmd/cookies"
)

type ContextKey string
type ToHand struct {
	Value  string
	IsAuth bool
}

// константа норм глобальная переменная нет
const namecookie string = "username"

func AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var th ToHand
		secretKey := []byte("mandarinmandarin")
		namecookie := "username"
		value, err := cookies.ReadEncrypted(r, namecookie, secretKey)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				h.ServeHTTP(w, r)
				return
			} else {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			th.IsAuth = true
			th.Value = value
		}
		var key ContextKey = "Name"
		ctx := context.WithValue(r.Context(), key, th)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NewCookie(w http.ResponseWriter, name string) error {
	secretKey := []byte("mandarinmandarin")
	cookie := http.Cookie{
		Name:  namecookie,
		Value: name,
	}
	err1 := cookies.WriteEncrypted(w, cookie, secretKey)
	// if err1 != nil {
	// 	return err1
	// }
	// return nil
	//можно просто ошибку возвращять
	return err1
}
