package ports

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/malanak2/nextap-chat/gen/chatdb/public/model"
	. "github.com/malanak2/nextap-chat/gen/chatdb/public/table"
)

func DeleteUserMessageById(id int32) error {
	slog.Info("Deleting UserMessage", "id", id)
	stmtDelUM := UserMessage.DELETE().WHERE(UserMessage.ID.EQ(postgres.Int(int64(id)))).RETURNING(UserMessage.AllColumns)
	var destDUM []struct {
		model.UserMessage
	}
	err := stmtDelUM.Query(Db, &destDUM)
	if err != nil {
		slog.Error("Database error deleting from usermessage table", "error", err.Error())
		return err
	}
	return nil
}

func SelectUserMessagesByUserId(userID int32) ([]struct{ model.UserMessage }, error) {
	slog.Info("Selecting usermessage by user", "uid", userID)
	stmtUserMessage := UserMessage.SELECT(UserMessage.AllColumns).WHERE(UserMessage.User.EQ(postgres.Int(int64(userID))))
	var destUM []struct {
		model.UserMessage
	}
	err := stmtUserMessage.Query(Db, &destUM)
	if err != nil {
		// User has sent no messages
		if !strings.Contains(err.Error(), "no rows in result set") {
			return []struct{ model.UserMessage }{}, nil
		}
		slog.Error("Error searching the UserMessage table", "error", err.Error())
		return []struct{ model.UserMessage }{}, err
	}
	return destUM, nil
}

func InsertUserMessage(msg model.Message, uid int) error {
	slog.Info("Inserting UserMessage", "msgId", msg.ID, "uid", uid)
	// Insert into user message table
	stmt := UserMessage.INSERT(UserMessage.User, UserMessage.Message).VALUES(uid, msg.ID).RETURNING(UserMessage.AllColumns)
	var destUM struct {
		model.UserMessage
	}
	err := stmt.Query(Db, &destUM)
	if err != nil {
		slog.Error("Error inserting into the UserMessage table", "error", err, "msgId", msg.ID, "uid", uid)
		return errors.New("error inserting into the UserMessage table. Please contact an administrator")
	}
	return nil
}

func SelectMessagesByUserId(uid int, limit, page int) ([]struct{ model.UserMessage }, error) {
	slog.Info("Selecting messages by user", "uid", uid)
	stmt := UserMessage.SELECT(UserMessage.Message).WHERE(UserMessage.User.EQ(postgres.Int(int64(uid)))).LIMIT(int64(limit)).OFFSET(int64((page - 1) * limit))

	var destM []struct {
		model.UserMessage
	}

	err := stmt.Query(Db, &destM)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result") {
			return destM, errors.New("no messages found")
		}
		slog.Error("Error selecting from the UserMessage table", "error", err.Error(), "uid", uid)
		return destM, errors.New("error selecting from the UserMessage table. Please contact an administrator")
	}
	return destM, nil
}
