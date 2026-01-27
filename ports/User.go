package ports

import (
	"errors"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/golang-jwt/jwt/v5"
	"github.com/malanak2/nextap-chat/gen/chatdb/public/model"
	. "github.com/malanak2/nextap-chat/gen/chatdb/public/table"
)

func CreateUser(name string, password string) (model.User, error) {
	// TODO: validate password to match expectation
	slog.Info("Creating user", "Username", name)
	stmt := User.INSERT(User.Username).VALUES(
		name,
	).RETURNING(User.AllColumns)
	var dest struct {
		model.User
	}
	err := stmt.Query(Db, &dest)
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
	slog.Info("Searching for user", "text", text)
	stmtSearch := User.SELECT(User.AllColumns).WHERE(postgres.LOWER(User.Username).LIKE(postgres.LOWER(postgres.String("%" + text + "%")))).LIMIT(int64(limit)).OFFSET(int64((pageNo - 1) * limit))
	var dest []struct {
		model.User
	}
	err := stmtSearch.Query(Db, &dest)
	return dest, err
}

func DeleteUser(id int32) error {
	slog.Info("Deleting user", "id", id)
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
	err = stmtDelUser.Query(Db, &destDMSG)
	if err != nil {
		slog.Error("Database error deleting from User table", "error", err.Error())
		return err
	}
	return nil
}

func ChangeUsername(id int32, username string) error {
	slog.Info("Changing username", "id", id, "username", username)
	stmtCU := User.UPDATE(User.Username).WHERE(User.ID.EQ(postgres.Int(int64(id)))).SET(postgres.String(username)).RETURNING(User.AllColumns)
	var dest struct {
		model.User
	}
	err := stmtCU.Query(Db, &dest)
	if err != nil {
		slog.Error("Database error editing username", "error", err.Error())
		return err
	}
	return nil
}

func GetUserById(id int32) (struct{ model.User }, error) {
	slog.Info("Getting user", "id", id)
	stmtAuthor := User.SELECT(User.AllColumns).WHERE(User.ID.EQ(postgres.Int(int64(id))))
	var dest struct {
		model.User
	}
	err := stmtAuthor.Query(Db, &dest)
	if err != nil {
		slog.Error("Database error getting user", "error", err.Error())
		return struct{ model.User }{}, err
	}
	return dest, nil
}

func UserLogin(username, password string) (string, error) {
	slog.Info("User login", "username", username)
	// TODO: verify password
	stmt := User.SELECT(User.AllColumns).FROM(User).WHERE(postgres.AND(User.Username.EQ(postgres.Text(username))))
	var dest struct {
		model.User
	}
	err := stmt.Query(Db, &dest)
	if err != nil {
		return "", errors.New("Invalid username or password")
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// Expires in about a month
		"exp":      jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
		"username": username,
		"userId":   dest.ID,
	})
	s, _ := t.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return s, nil
}

func GetAllUsers(limit int, pageNo int) ([]struct{ model.User }, error) {
	stmt := User.SELECT(User.ID, User.Username).FROM(User).LIMIT(int64(limit)).OFFSET(int64((pageNo - 1) * limit))
	var dest []struct {
		model.User
	}
	err := stmt.Query(Db, &dest)
	return dest, err
}

func UserExists(uid int) (bool, error) {
	stmtS := User.SELECT(User.AllColumns).WHERE(User.ID.EQ(postgres.Int(int64(uid))))
	var destU struct {
		model.Channel
	}
	err := stmtS.Query(Db, &destU)
	if err != nil {
		slog.Warn("User with invalid ID found - probably fine though", "error", err.Error(), "uid", uid)
		return false, errors.New("no user with that id")
	}
	return true, nil
}
