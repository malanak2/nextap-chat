package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/gorilla/mux"
	. "github.com/malanak2/nextap-chat/.gen/chatdb/public/table"

	"github.com/malanak2/nextap-chat/.gen/chatdb/public/model"

	"github.com/malanak2/nextap-chat/domain"
	"github.com/malanak2/nextap-chat/handlers"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Home Page!")
}

func handlerUsers(w http.ResponseWriter, r *http.Request) {
	stmt := SELECT(User.AllColumns).FROM(User)
	var dest []struct {
		model.User
	}
	err := stmt.Query(domain.Db, &dest)
	if err != nil {
		http.Error(w, `Database query error `+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Users Page")
}

func main() {
	var connectString = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_host"), os.Getenv("DB_port"), os.Getenv("DB_user"), os.Getenv("DB_pass"), os.Getenv("DB_name"))
	err := domain.InitDb(connectString)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handleRoot).Methods("GET")
	r.HandleFunc("/users", handlerUsers).Methods("GET")
	r.HandleFunc("/create-user", handlers.HandleUserCreate).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", r))
}
