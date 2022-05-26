package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"strconv"

	"github.com/Outer2g/interstellar-cab/pkg/auth"
)

type AuthController struct {
	UserRepository
}

type responseOutput struct {
	Email string
	Token string
}

type credentials struct {
	Email        string
	PasswordHash string
}

func NewUserAuth() *AuthController {
	return &AuthController{NewUserInMemoryDatabase()}
}

func (u AuthController) HandleSignupUser(rw http.ResponseWriter, r *http.Request) {
	userObject := User{}
	json.NewDecoder(r.Body).Decode(&userObject)

	_, err := mail.ParseAddress(userObject.Email)
	if err != nil {
		rw.WriteHeader(http.StatusNotAcceptable)
		log.Println("Received a request with a non-valid email")
		return
	}

	exists := u.AddUser(userObject.Email, userObject.Passwordhash, userObject.Vip)
	if exists {
		rw.WriteHeader(http.StatusNotAcceptable)
		log.Println("Received a request to register a user that already exists, discarding")
		return
	}

	token, err := auth.GenerateJwtToken(userObject.Email, userObject.Vip)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println("Failed To Generate New JWT Token!")
		return
	}

	respondJson(rw, responseOutput{
		Token: token,
		Email: userObject.Email,
	})
}

func respondJson(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func (u AuthController) HandleLoginUser(rw http.ResponseWriter, r *http.Request) {
	credentials := credentials{}
	json.NewDecoder(r.Body).Decode(&credentials)

	_, err := mail.ParseAddress(credentials.Email)
	if err != nil {
		rw.WriteHeader(http.StatusNotAcceptable)
		log.Println("Received a request with a non-valid email")
		return
	}

	user := u.GetUser(credentials.Email)
	if user == nil {
		log.Println("Request for login failed because user is not in database")
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	if user.Passwordhash != credentials.PasswordHash {
		rw.WriteHeader(http.StatusNotFound)
		log.Println("Request for login failed because wrong credentials were passed")
		return
	}

	token, err := auth.GenerateJwtToken(user.Email, user.Vip)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println("Failed To Generate New JWT Token!")
		return
	}

	respondJson(rw, responseOutput{
		Token: token,
		Email: user.Email,
	})
}

func (u AuthController) CheckAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Token")

		if len(authToken) < 2 {
			fmt.Errorf("Token not provided!")
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		claims, err := auth.VerifyJwtToken(authToken)
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
