package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bugjoe/family-tree-backend/models"
	"github.com/bugjoe/family-tree-backend/persistence"
	"github.com/bugjoe/family-tree-backend/utils"
	"github.com/gorilla/mux"

	"github.com/dgrijalva/jwt-go"
)

var secret = []byte("OSR6eqMv7N01PlPFZyBS1k508daeP8hC15dwRQ5pzr7hwsIOcAWQuhdZlGUKHIQw")

type claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// CreateAccount handles account create requests
func CreateAccount(response http.ResponseWriter, request *http.Request) {
	account := new(models.Account)
	err := utils.ExtractFromRequest(request, account)
	if err != nil {
		log.Println("ERROR: Can't extract account from request:", err)
		http.Error(response, "Error while parsing request body", 500)
		return
	}

	log.Println("Create account with email:", account.Payload.Email)
	err = persistence.InsertNewAccount(account)
	if err != nil {
		log.Println("ERROR: Can't create new account:", err)
		if err == persistence.ErrAccountAlreadyExists {
			http.Error(response, "User already exists", 403)
		} else {
			http.Error(response, "Error while creating new account", 500)
		}
		return
	}

	response.WriteHeader(200)
}

// GetAccount handles get requests for accounts
func GetAccount(response http.ResponseWriter, request *http.Request) {
	email := mux.Vars(request)["email"]
	account, err := persistence.GetAccountByEmail(email)
	if err != nil {
		log.Printf("ERROR: Can't get account with email=%s from database: %v\n", email, err)
		http.Error(response, "Error while getting account from database", 500)
		return
	}

	account.Payload.Password = ""

	json, err := json.Marshal(account)
	if err != nil {
		log.Printf("ERROR: Can't create with email=%s to a JSON object: %v\n", email, err)
		http.Error(response, "Error while create JSON", 500)
		return
	}

	response.Write(json)
}

// Login handles login requests. If the login is successfull, it will
// respond with a JWT.
func Login(response http.ResponseWriter, request *http.Request) {
	account := new(models.Account)
	err := utils.ExtractFromRequest(request, account)
	if err != nil {
		log.Println("ERROR: Can't extract account from request", err)
		http.Error(response, "Error while parsing request body", 500)
		return
	}

	accountFromDb, err := persistence.GetAccountByEmail(account.Payload.Email)
	if err != nil {
		log.Printf("ERROR: Can't get account with email=%s from database: %v\n", account.Payload.Email, err)
		http.Error(response, "Error while getting account from database", 500)
		return
	}

	if accountFromDb == nil {
		http.Error(response, "Invalid login", 401)
		return
	}

	passwordHash, err := account.GetPasswordHash()
	if err != nil {
		log.Println("ERROR: Can't get password hash from account", err)
		http.Error(response, "Error while creating password hash", 500)
		return
	}

	if passwordHash != accountFromDb.Payload.Password {
		http.Error(response, "Invalid login", 401)
		return
	}

	claims := claims{
		Email: accountFromDb.Payload.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 60000,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		log.Println("ERROR: Can't sign JWT", err)
		http.Error(response, "Error wile signing JWT", 500)
		return
	}

	response.Write([]byte(tokenString))
}
