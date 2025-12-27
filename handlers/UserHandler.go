package handlers

import (
	"net/http"

	. "github.com/malanak2/nextap-chat/.gen/chatdb/public/table"

	"github.com/malanak2/nextap-chat/.gen/chatdb/public/model"
	"github.com/malanak2/nextap-chat/domain"
)

func HandleUserCreate(w http.ResponseWriter, r *http.Request) {
	stmt := User.INSERT(User.Username).VALUES(
		"new_user",
	).RETURNING(User.AllColumns)
	var dest struct {
		model.User
	}
	err := stmt.Query(domain.Db, &dest)
	if err != nil {
		http.Error(w, `Database query error `+err.Error(), http.StatusInternalServerError)
		return
	}
}
