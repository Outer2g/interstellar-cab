package user

import (
	"encoding/json"
	"log"
	"net/http"
	"net/mail"

	"github.com/Outer2g/interstellar-cab/pkg/auth"
	user "github.com/Outer2g/interstellar-cab/pkg/user/repository"
)

type AuthController struct {
	user.UserRepository
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
	return &AuthController{user.NewUserInMemoryDatabase()}
}

func (u AuthController) HandleSignupUser(rw http.ResponseWriter, r *http.Request) {
	userObject := user.User{}
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
