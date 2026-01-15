package ports

import (
	"errors"
	"strings"

	"github.com/malanak2/nextap-chat/domain"
	"github.com/malanak2/nextap-chat/gen/chatdb/public/model"
	. "github.com/malanak2/nextap-chat/gen/chatdb/public/table"
)

func CreateUser(name string, password string) (model.User, error) {
	// TODO: validate password to match expectation
	stmt := User.INSERT(User.Username).VALUES(
		name,
	).RETURNING(User.AllColumns)
	var dest struct {
		model.User
	}
	err := stmt.Query(domain.Db, &dest)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"User_username_key\"") {
			return model.User{}, errors.New("A user with this username already exists.")
		}

		return model.User{}, errors.New(`Database query error ` + err.Error())
	}
	return dest.User, nil
}
