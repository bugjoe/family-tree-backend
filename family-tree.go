package main

import (
	"log"
	"net/http"

	"github.com/bugjoe/family-tree-backend/handlers"
	muxHandler "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	log.Println("Start Family Tree Server")
	router := mux.NewRouter()
	router.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	router.HandleFunc("/accounts", handlers.CreateAccountHandler).Methods("PUT")
	router.HandleFunc("/accounts/{email}", handlers.GetAccountHandler).Methods("GET")
	allowedMethods := muxHandler.AllowedMethods([]string{"GET", "POST", "OPTIONS", "PUT", "HEAD", "DELETE"})
	allowedHeaders := muxHandler.AllowedHeaders([]string{"x-requested-with", "authorization", "content-type"})
	corsEnabledRouter := muxHandler.CORS(allowedMethods, allowedHeaders)(router)
	log.Println(http.ListenAndServe(":8080", corsEnabledRouter))
	log.Println("Server shutdown ...")
}
