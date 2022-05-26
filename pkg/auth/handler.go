package auth

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func CheckAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Token")

		if len(authToken) < 2 {
			fmt.Errorf("Token not provided!")
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		claims, err := VerifyJwtToken(authToken)
		if err != nil {
			log.Println(fmt.Errorf(err.Error()))
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		r.Header.Set("Email", claims.Email)
		r.Header.Set("Vip", strconv.FormatBool(claims.Vip))
		next.ServeHTTP(w, r)
	})
}
