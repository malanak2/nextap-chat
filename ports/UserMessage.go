package ports

import (
	"log/slog"
	"strings"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/malanak2/nextap-chat/domain"
	"github.com/malanak2/nextap-chat/gen/chatdb/public/model"

	. "github.com/malanak2/nextap-chat/gen/chatdb/public/table"
)

func DeleteUserMessageById(id int32) error {
	stmtDelUM := UserMessage.DELETE().WHERE(UserMessage.ID.EQ(postgres.Int(int64(id)))).RETURNING(UserMessage.AllColumns)
	var destDUM []struct {
		model.UserMessage
	}
	err := stmtDelUM.Query(domain.Db, &destDUM)
	if err != nil {
		slog.Error("Database error deleting from usermessage table", "error", err.Error())
		return err
	}
	return nil
}

func SelectUserMessagesByUserId(userID int32) ([]struct{ model.UserMessage }, error) {
	stmtUserMessage := UserMessage.SELECT(UserMessage.AllColumns).WHERE(UserMessage.User.EQ(postgres.Int(int64(userID))))
	var destUM []struct {
		model.UserMessage
	}
	err := stmtUserMessage.Query(domain.Db, &destUM)
	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result set") {
			slog.Error("Error searching the UserMessage table", "error", err.Error())
			return []struct{ model.UserMessage }{}, err
		}
	}
	return destUM, nil
}
