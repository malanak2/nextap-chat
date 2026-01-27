package ports

import (
	"errors"
	"log/slog"
	"strconv"
	"strings"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/malanak2/nextap-chat/gen/chatdb/public/model"
	. "github.com/malanak2/nextap-chat/gen/chatdb/public/table"
)

func DeleteMessageById(id int32) error {
	stmtDelMsg := Message.DELETE().WHERE(Message.ID.EQ(postgres.Int(int64(id))))
	slog.Info("Deleting message", "id", id)
	err := stmtDelMsg.Query(Db, nil)
	if err != nil {
		slog.Error("Database error deleting from Message table", "error", err.Error())
		return err
	}
	return nil
}

func SendMessage(content string, user int) (struct{ model.Message }, error) {
	slog.Info("Sending message", "content", content, "user", user)
	transaction, err := Db.Begin()
	if err != nil {
		slog.Error("Failed to open transaction", "error", err.Error())
		return struct{ model.Message }{}, ErrorDatabase
	}
	stmt := Message.INSERT(Message.Content).VALUES(postgres.String(content)).RETURNING(Message.ID)

	var destM struct {
		model.Message
	}
	err = stmt.Query(transaction, &destM)
	if err != nil {
		transaction.Rollback()
		return struct{ model.Message }{}, errors.New("Error inserting into the message table. Message too long?")
	}
	err = InsertUserMessage(destM.Message, user, transaction)
	if err != nil {
		transaction.Rollback()
		return struct{ model.Message }{}, err
	}
	transaction.Commit()
	return destM, nil
}

func SelectMessageById(id int) (struct {
	model.Message
	Author struct{ model.User }
}, error) {
	slog.Info("Selecting message", "id", id)
	stmt := postgres.SELECT(Message.AllColumns, User.AllColumns).WHERE(Message.ID.EQ(postgres.Int(int64(id)))).FROM(Message.INNER_JOIN(UserMessage, Message.ID.EQ(UserMessage.Message)).INNER_JOIN(User, UserMessage.User.EQ(User.ID)))

	var destM struct {
		model.Message
		Author struct {
			model.User
		}
	}

	err := stmt.Query(Db, &destM)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result") {
			return destM, errors.New("no messages with the id " + strconv.FormatInt(int64(id), 10))
		}
		slog.Error("Error selecting from the Message And/Or User table", "error", err.Error(), "msgId", id)
		return destM, errors.New("Error selecting from the Message table. Please contact an administrator. " + err.Error())
	}
	return destM, nil
}

func SelectMessagesByContent(content string, limit, page int) ([]struct{ model.Message }, error) {
	slog.Info("Selecting messages by content", "content")
	stmtSearch := Message.SELECT(Message.AllColumns).WHERE(postgres.LOWER(Message.Content).LIKE(postgres.LOWER(postgres.String("%" + content + "%")))).LIMIT(int64(limit)).OFFSET(int64((page - 1) * limit))
	var dest []struct {
		model.Message
	}
	err := stmtSearch.Query(Db, &dest)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return dest, nil
		}
		slog.Error("Error selecting from the Message table", "error", err.Error(), "txt", content)
		return []struct{ model.Message }{}, errors.New("Error selecting from the Message table. Please contact an administrator." + err.Error())
	}
	return dest, nil
}

func UpdateMessageById(id int, content string) (struct{ model.Message }, error) {
	slog.Info("Updating message", "id", id)
	stmt := Message.UPDATE(Message.Content).WHERE(Message.ID.EQ(postgres.Int(int64(id)))).SET(postgres.String(content)).RETURNING(Message.AllColumns)
	var destM struct {
		model.Message
	}
	err := stmt.Query(Db, &destM)
	if err != nil {
		slog.Error("Error updating Message table", "error", err.Error(), "msgId", id, "txt", content)
		return destM, errors.New("Error updating Message table. Please contact an administrator. " + err.Error())
	}
	return destM, nil
}

func MessageExists(id int) (bool, error) {
	stmt := Message.SELECT(Message.AllColumns).WHERE(Message.ID.EQ(postgres.Int(int64(id))))
	var destM struct {
		model.Message
	}
	err := stmt.Query(Db, &destM)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return false, nil
		}
		slog.Error("Error selecting from the Message table", "error", err.Error(), "msgId", id)
		return false, err
	}
	return true, nil
}

func DeleteMessage(id int) error {
	slog.Info("Deleting message", "id", id)
	tx, err := Db.Begin()
	if err != nil {
		slog.Error("Failed to open transaction", "error", err.Error())
		return ErrorDatabase
	}
	stmtDUM := UserMessage.DELETE().WHERE(UserMessage.Message.EQ(postgres.Int(int64(id))))
	err = stmtDUM.Query(tx, nil)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return errors.New("no message with this id")
		}
		slog.Error("Error deleting from UserMessage table", "error", err.Error(), "msgId", id)
		tx.Rollback()
		return err
	}
	stmt := Message.DELETE().WHERE(Message.ID.EQ(postgres.Int(int64(id))))
	err = stmt.Query(tx, nil)
	if err != nil {
		// This REALLY shouldn`t happen since we essentially verified the message exists since it was in usermessage table and the db constraints SHOULD make sure we are fine
		slog.Error("Error deleting from Message table", "error", err.Error(), "msgId", id)
		tx.Rollback()
		if strings.Contains(err.Error(), "no rows in result set") {
			return errors.New("no message with this id")
		}
		return err
	}
	tx.Commit()
	return nil
}

func GetAllMessages(limit, page int) ([]struct{ model.Message }, error) {
	stmt := Message.SELECT(Message.AllColumns).LIMIT(int64(limit)).OFFSET(int64((page - 1) * limit))
	var destM []struct {
		model.Message
	}
	err := stmt.Query(Db, &destM)
	if err != nil {
		slog.Error("Failed to get all messages", "error", err.Error())
		return nil, err
	}
	return destM, nil
}
