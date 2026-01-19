package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

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

	// Logging
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	// Routing
	r := mux.NewRouter()

	// Endpoints
	r.HandleFunc("/users", handlers.HandleGetAllUsers).Methods("GET")
	r.HandleFunc("/createUser", handlers.HandleUserCreate).Methods("POST")
	r.HandleFunc("/login", handlers.HandleUserLogin).Methods("POST")
	r.HandleFunc("/user/{id}", handlers.HandleGetUserById).Methods("GET")
	r.HandleFunc("/user/search/{txt}", handlers.HandleSearchUsers).Methods("GET")
	r.HandleFunc("/user/{id}/messages", handlers.HandleGetMessagesByUserId).Methods("GET")
	r.HandleFunc("/message/{id}", handlers.HandleGetMessageById).Methods("GET")
	r.HandleFunc("/message/search/{txt}", handlers.HandleSearchMessages).Methods("GET")
	r.HandleFunc("/messages", handlers.HandleGetAllMessages).Methods("GET")

	// Secure endpoints
	r.Handle("/user/{id}", handlers.JwtMiddleware(http.HandlerFunc(handlers.HandleDeleteUser))).Methods("DELETE")
	r.Handle("/sendMessage", handlers.JwtMiddleware(http.HandlerFunc(handlers.HandleSendMessage))).Methods("POST")
	r.Handle("/changeUsername", handlers.JwtMiddleware(http.HandlerFunc(handlers.HandleUserChangeName))).Methods("POST")
	r.Handle("/message/{id}/update", handlers.JwtMiddleware(http.HandlerFunc(handlers.HandleEditMessageById))).Methods("POST")
	r.Handle("/message/{id}", handlers.JwtMiddleware(http.HandlerFunc(handlers.HandleDeleteMessageById))).Methods("DELETE")

	// Swagger
	r.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", http.FileServer(http.Dir("docs"))))
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:" + os.Getenv("port") + "/docs/swagger.json"), //The url pointing to API definition
	)).Methods("GET")

	log.Fatal(http.ListenAndServe(":"+os.Getenv("port"), r))
}
