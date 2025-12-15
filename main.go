package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/gorilla/mux"
	. "github.com/malanak2/nextap_chat/.gen/chatdb/public/table"

	"github.com/malanak2/nextap_chat/.gen/chatdb/public/model"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Home Page!")
}

var db *sql.DB

func handlerUsers(w http.ResponseWriter, r *http.Request) {
	stmt := SELECT(User.AllColumns).FROM(User)
	var dest []struct {
		model.User
	}
	err := stmt.Query(db, &dest)
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Users Page")
}

func main() {
	var connectString = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_host"), os.Getenv("DB_port"), os.Getenv("DB_user"), os.Getenv("DB_pass"), os.Getenv("DB_name"))
	err := error(nil)
	db, err = sql.Open("postgres", connectString)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handleRoot).Methods("GET")
	r.HandleFunc("/users", handlerUsers).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}
