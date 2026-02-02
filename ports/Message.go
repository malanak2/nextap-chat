package ports

import (
	"database/sql"
	"log/slog"
	"strings"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/malanak2/nextap-chat/domain"
	"github.com/malanak2/nextap-chat/gen/chatdb/public/model"
	. "github.com/malanak2/nextap-chat/gen/chatdb/public/table"
)

func DeleteMessageById(id int32, tx *sql.Tx) error {
	stmtDelMsg := Message.DELETE().WHERE(Message.ID.EQ(postgres.Int(int64(id))))
	slog.Info("Deleting message", "id", id)
	var dest struct {
	}
	err := stmtDelMsg.Query(tx, &dest)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil
		}
		slog.Error("Database error deleting from Message table", "error", err.Error())
		return err
	}
	return nil
}

func SendMessage(content string, user int) (domain.Message, error) {
	slog.Info("Sending message", "content", content, "user", user)
	transaction, err := Db.Begin()
	if err != nil {
		slog.Error("Failed to open transaction", "error", err.Error())
		return domain.Message{}, ErrorDatabase
	}
	stmt := Message.INSERT(Message.Content, Message.UserID).VALUES(postgres.String(content), postgres.Int(int64(user))).RETURNING(Message.ID)

	var destM struct {
		model.Message
	}
	err = stmt.Query(transaction, &destM)
	if err != nil {
		transaction.Rollback()
		return domain.Message{}, ErrorTooLong
	}
	err = transaction.Commit()
	if err != nil {
		slog.Error("Failed to commit", "error", err)
		return domain.Message{}, err
	}
	return domain.Message{ID: int(destM.ID), Content: destM.Content}, nil
}

func SelectMessageById(id int) (domain.MessageWithAuthor, error) {
	slog.Info("Selecting message", "id", id)
	stmt := postgres.SELECT(Message.AllColumns, User.ID, User.Username).WHERE(Message.ID.EQ(postgres.Int(int64(id)))).FROM(Message.INNER_JOIN(User, Message.UserID.EQ(User.ID)))

	var destM struct {
		model.Message
		Author struct {
			model.User
		}
	}

	err := stmt.Query(Db, &destM)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result") {
			return domain.MessageWithAuthor{}, ErrorNoResult
		}
		slog.Error("Error selecting from the Message And/Or User table", "error", err.Error(), "msgId", id)
		return domain.MessageWithAuthor{}, ErrorDatabase
	}
	return domain.MessageWithAuthor{Message: domain.Message{ID: int(destM.ID), Content: destM.Content}, Author: domain.User{ID: int(destM.Author.ID), Username: destM.Author.Username}}, nil
}

func SelectMessagesByContent(content string, limit, page int) ([]domain.Message, error) {
	slog.Info("Selecting messages by content", "content")
	stmtSearch := Message.SELECT(Message.AllColumns).WHERE(postgres.LOWER(Message.Content).LIKE(postgres.LOWER(postgres.String("%" + content + "%")))).LIMIT(int64(limit)).OFFSET(int64((page - 1) * limit))
	var dest []struct {
		model.Message
	}
	err := stmtSearch.Query(Db, &dest)
	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result set") {
			slog.Error("Error selecting from the Message table", "error", err.Error(), "txt", content)
			return []domain.Message{}, ErrorDatabase
		}

	}
	var ret []domain.Message
	for _, v := range dest {
		destConverted := domain.Message{
			ID:      int(v.ID),
			Content: v.Content,
		}
		ret = append(ret, destConverted)
	}
	return ret, nil
}

func UpdateMessageById(id int, content string) (domain.Message, error) {
	slog.Info("Updating message", "id", id)
	stmt := Message.UPDATE(Message.Content).WHERE(Message.ID.EQ(postgres.Int(int64(id)))).SET(postgres.String(content)).RETURNING(Message.AllColumns)
	var destM struct {
		model.Message
	}
	err := stmt.Query(Db, &destM)
	if err != nil {
		slog.Error("Error updating Message table", "error", err.Error(), "msgId", id, "txt", content)
		return domain.Message{}, ErrorDatabase
	}
	return domain.Message{
		ID:      int(destM.ID),
		Content: destM.Content,
	}, nil
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
	tx, err := Db.Begin()
	if err != nil {
		slog.Error("Failed to open transaction", "error", err.Error())
		return ErrorDatabase
	}

	stmt := Message.DELETE().WHERE(Message.ID.EQ(postgres.Int(int64(id))))
	var dests int
	err = stmt.Query(tx, &dests)
	if err != nil {
		slog.Error("Error deleting from Message table", "error", err.Error(), "msgId", id)
		tx.Rollback()
		if strings.Contains(err.Error(), "no rows in result set") {
			return ErrorNoResult
		}
		return err
	}
	tx.Commit()
	return nil
}

func GetAllMessages(limit, page int) ([]domain.Message, error) {
	stmt := Message.SELECT(Message.AllColumns).LIMIT(int64(limit)).OFFSET(int64((page - 1) * limit))
	var destM []struct {
		model.Message
	}
	err := stmt.Query(Db, &destM)
	if err != nil {
		slog.Error("Failed to get all messages", "error", err.Error())
		return nil, err
	}
	var ret []domain.Message
	for _, v := range destM {
		destConverted := domain.Message{
			ID:      int(v.ID),
			Content: v.Content,
		}
		ret = append(ret, destConverted)
	}
	return ret, nil
}

func SelectMessagesByUserId(id, limit, page int) ([]domain.Message, error) {
	stmt := Message.SELECT(Message.ID).WHERE(Message.UserID.EQ(postgres.Int(int64(id)))).LIMIT(int64(limit)).OFFSET(int64((page - 1) * limit))
	var dest []struct {
		model.Message
	}
	err := stmt.Query(Db, &dest)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return []domain.Message{}, ErrorNoResult
		}
	}
	var ret []domain.Message
	for _, v := range dest {
		destConverted := domain.Message{
			ID:      int(v.ID),
			Content: v.Content,
		}
		ret = append(ret, destConverted)
	}
	return ret, nil
}
