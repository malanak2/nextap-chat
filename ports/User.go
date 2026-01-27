package ports

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/malanak2/nextap-chat/gen/chatdb/public/model"
	. "github.com/malanak2/nextap-chat/gen/chatdb/public/table"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(name string, password string) (model.User, error) {
	slog.Info("Creating user", "Username", name)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Failed to hash password", "error", err.Error(), "password", password)
		return model.User{}, ErrorHashing
	}
	stmt := User.INSERT(User.Username, User.Password).VALUES(
		name,
		hash,
	).RETURNING(User.AllColumns)
	var dest struct {
		model.User
	}
	err = stmt.Query(Db, &dest)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"User_username_key\"") {
			return model.User{}, ErrorDuplicateUsername
		}
		slog.Error("Error inserting into User", "error", err.Error())
		return model.User{}, ErrorDatabase
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
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return dest, ErrorNoResult
		}
		return dest, ErrorDatabase
	}
	return dest, nil
}

func DeleteUser(id int32) error {
	slog.Info("Deleting user", "id", id)
	// Get all messages for user
	transaction, err := Db.Begin()
	if err != nil {
		return ErrorDatabase
	}
	stmt := Message.SELECT(Message.ID).WHERE(Message.UserID.EQ(postgres.Int(int64(id))))
	var dest []struct {
		id int
	}
	err = stmt.Query(transaction, &dest)
	if err != nil {
		transaction.Rollback()
		if strings.Contains(err.Error(), "no rows in result set") {
			return ErrorNoResult
		}
		return ErrorDatabase
	}
	for i := 0; i < len(dest); i++ {
		err = DeleteMessageById(int32(dest[i].id), transaction)
		if err != nil {
			transaction.Rollback()
			slog.Error("Database error deleting from Message table", "error", err.Error())
			return ErrorDatabase
		}
	}
	// Finally, delete user
	stmtDelUser := User.DELETE().WHERE(User.ID.EQ(postgres.Int(int64(id)))).RETURNING(User.AllColumns)
	var destDMSG []struct {
		model.User
	}
	err = stmtDelUser.Query(transaction, &destDMSG)
	if err != nil {
		transaction.Rollback()
		slog.Error("Database error deleting from User table", "error", err.Error())
		return ErrorDatabase
	}
	transaction.Commit()
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
		return ErrorDatabase
	}
	return nil
}

func GetUserById(id int32) (struct{ model.User }, error) {
	slog.Info("Getting user", "id", id)
	stmtAuthor := User.SELECT(User.ID, User.Username).WHERE(User.ID.EQ(postgres.Int(int64(id))))
	var dest struct {
		model.User
	}
	err := stmtAuthor.Query(Db, &dest)
	if err != nil {
		slog.Error("Database error getting user", "error", err.Error())
		return struct{ model.User }{}, ErrorDatabase
	}
	return dest, nil
}

func UserLogin(username, password string) (model.User, error) {
	slog.Info("User login", "username", username)
	stmt := User.SELECT(User.AllColumns).FROM(User).WHERE(postgres.AND(User.Username.EQ(postgres.Text(username))))
	var dest struct {
		model.User
	}
	err := stmt.Query(Db, &dest)
	if err != nil {
		return model.User{}, ErrorBadCredentials
	}
	err = bcrypt.CompareHashAndPassword([]byte(dest.Password), []byte(password))
	if err != nil {
		return model.User{}, ErrorBadCredentials
	}
	return dest.User, nil
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
		model.User
	}
	err := stmtS.Query(Db, &destU)
	if err != nil {
		slog.Warn("User with invalid ID found - probably fine though", "error", err.Error(), "uid", uid)
		return false, errors.New("no user with that id")
	}
	return true, nil
}
