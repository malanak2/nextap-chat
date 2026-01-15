package ports

import (
	"log/slog"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/malanak2/nextap-chat/domain"
	"github.com/malanak2/nextap-chat/gen/chatdb/public/model"
	. "github.com/malanak2/nextap-chat/gen/chatdb/public/table"
)

func DeleteMessageById(id int32) error {
	stmtDelMsg := Message.DELETE().WHERE(Message.ID.EQ(postgres.Int(int64(id)))).RETURNING(Message.AllColumns)
	var destDMSG []struct {
		model.Message
	}
	err := stmtDelMsg.Query(domain.Db, &destDMSG)
	if err != nil {
		slog.Error("Database error deleting from Message table", "error", err.Error())
		return err
	}
	return nil
}
