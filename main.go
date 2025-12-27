package main

import (
	"log"
	"net/http"
	"os"

	ghandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/malanak2/nextap-chat/domain"
	"github.com/malanak2/nextap-chat/handlers"
	"github.com/swaggo/http-swagger/v2"
)

//	@title			Chat app
//	@version		0.1
//	@description	Chat application backend

//	@host		localhost:8080
//	@BasePath	/

//	@securityDefinitions.basic JWTTokenBasic

func main() {
	// Init db
	err := domain.InitDb()
	if err != nil {
		panic(err)
	}

	// Routing
	r := mux.NewRouter()

	// Endpoints
	r.HandleFunc("/users", handlers.HandleGetAllUsers).Methods("GET")
	r.HandleFunc("/createUser", handlers.HandleUserCreate).Methods("POST")
	r.HandleFunc("/login", handlers.HandleUserLogin).Methods("POST")
	r.HandleFunc("/user/{id}", handlers.HandleGetUserById).Methods("GET")
	r.HandleFunc("/user/{id}/messages", handlers.HandleGetMessagesByUserId).Methods("GET")

	// Secure endpoints
	r.Handle("/sendMessage", handlers.JwtMiddleware(http.HandlerFunc(handlers.HandleSendMessage))).Methods("POST")
	r.Handle("/changeUsername", handlers.JwtMiddleware(http.HandlerFunc(handlers.HandleUserChangeName))).Methods("POST")

	// Swagger
	r.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", http.FileServer(http.Dir("docs"))))
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:" + os.Getenv("port") + "/docs/swagger.json"), //The url pointing to API definition
	)).Methods("GET")

	loggedRouter := ghandlers.LoggingHandler(os.Stdout, r)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("port"), loggedRouter))
}
