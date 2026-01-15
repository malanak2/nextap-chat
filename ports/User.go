package ports

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/malanak2/nextap-chat/domain"
	"github.com/malanak2/nextap-chat/gen/chatdb/public/model"
	. "github.com/malanak2/nextap-chat/gen/chatdb/public/table"
)

func CreateUser(name string, password string) (model.User, error) {
	// TODO: validate password to match expectation
	slog.Info("CreateUser", "Username", name)
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
		slog.Error("Error searching the database", "error", err.Error())
		return model.User{}, errors.New(`Database query error ` + err.Error())
	}
	return dest.User, nil
}

func SearchUsers(text string, limit int, pageNo int) ([]struct{ model.User }, error) {
	slog.Info("SearchUsers", "text", text, "limit", limit, "pageNo", pageNo)
	stmtSearch := User.SELECT(User.AllColumns).WHERE(postgres.LOWER(User.Username).LIKE(postgres.LOWER(postgres.String("%" + text + "%")))).LIMIT(int64(limit)).OFFSET(int64((pageNo - 1) * limit))
	var dest []struct {
		model.User
	}
	err := stmtSearch.Query(domain.Db, &dest)
	return dest, err
}

func DeleteUser(id int32) error {
	// Get all messages for user
	destUM, err := SelectUserMessagesByUserId(id)
	if err != nil {
		return err
	}
	// For every message delete its entry in UserMessage and then delete the message
	for i := 0; i < len(destUM); i++ {
		err = DeleteUserMessageById(destUM[i].ID)
		if err != nil {
			return err
		}
		err = DeleteMessageById(destUM[i].Message)
		if err != nil {
			return err
		}
	}
	// Finally, delete user
	stmtDelUser := User.DELETE().WHERE(User.ID.EQ(postgres.Int(int64(id)))).RETURNING(User.AllColumns)
	var destDMSG []struct {
		model.User
	}
	err = stmtDelUser.Query(domain.Db, &destDMSG)
	if err != nil {
		slog.Error("Database error deleting from User table", "error", err.Error())
		return err
	}
	return nil
}

func ChangeUsername(id int32, username string) error {
	stmtCU := User.UPDATE(User.Username).WHERE(User.ID.EQ(postgres.Int(int64(id)))).SET(postgres.String(username)).RETURNING(User.AllColumns)
	var dest struct {
		model.User
	}
	err := stmtCU.Query(domain.Db, &dest)
	if err != nil {
		slog.Error("Database error editing username", "error", err.Error())
		return err
	}
	return nil
}

func GetUserById(id int32) (struct{ model.User }, error) {
	stmtAuthor := User.SELECT(User.AllColumns).WHERE(User.ID.EQ(postgres.Int(int64(id))))
	var dest struct {
		model.User
	}
	err := stmtAuthor.Query(domain.Db, &dest)
	if err != nil {
		slog.Error("Database error getting user", "error", err.Error())
		return struct{ model.User }{}, err
	}
	return dest, nil
}
