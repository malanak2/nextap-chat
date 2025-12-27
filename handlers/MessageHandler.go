package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/malanak2/nextap-chat/domain"
	"github.com/malanak2/nextap-chat/gen/chatdb/public/model"
	. "github.com/malanak2/nextap-chat/gen/chatdb/public/table"
)

// HandleSendMessage godoc
//
// @Summary			 	Send message
// @Description		 	Sends a message as user specified with jwt
// @Tags				message,user
// @Accept				json
// @Produce				json
// @Success				200 {object}	model.Message
// @Failure				400 {object}	string
// @Failure				500 {object}	string
// @Security			JWTTokenBasic
// @Router				/sendMessage [post]
func HandleSendMessage(w http.ResponseWriter, r *http.Request) {
	// Parse context
	uid, ok := r.Context().Value("userId").(int)
	if !ok {
		// If this happens either the login function gives out a malformed but valid jwt, or somebody has our signing key - not good either way
		http.Error(w, "Invalid jwt. Please contact a site administrator", http.StatusInternalServerError)

		return
	}
	var body domain.SendMessage
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// Verify channel and user both exist
	stmtS := Channel.SELECT(Channel.AllColumns).WHERE(Channel.ID.EQ(postgres.Int(int64(body.Channel))))
	var destC struct {
		model.Channel
	}
	err = stmtS.Query(domain.Db, &destC)
	if err != nil {
		http.Error(w, "No channel with that id", http.StatusBadRequest)
		return
	}
	stmtS = User.SELECT(User.AllColumns).WHERE(User.ID.EQ(postgres.Int(int64(uid))))
	var destU struct {
		model.Channel
	}
	err = stmtS.Query(domain.Db, &destU)
	if err != nil {
		http.Error(w, "No user with that id", http.StatusBadRequest)
		fmt.Fprintln(os.Stderr, "User with invalid ID found - probably fine though", err)
		return
	}

	// Verify channel message
	if len(body.Content) == 0 {
		http.Error(w, "Content required", http.StatusBadRequest)
		return
	}
	if len(body.Content) > 1000 {
		http.Error(w, "Content too long", http.StatusBadRequest)
		return
	}

	// Insert message

	stmt := Message.INSERT(Message.Content).VALUES(postgres.String(body.Content)).RETURNING(Message.ID)

	var destM struct {
		model.Message
	}
	err = stmt.Query(domain.Db, &destM)
	if err != nil {
		http.Error(w, "Error inserting into the message table. Message too long?", http.StatusBadRequest)
		return
	}

	// Insert into message channel table
	stmt = MessageChannel.INSERT(MessageChannel.Channel, MessageChannel.Message).VALUES(body.Channel, destM.ID).RETURNING(MessageChannel.ID)
	var destMC struct {
		model.Message
	}
	err = stmt.Query(domain.Db, &destMC)
	if err != nil {
		http.Error(w, "Error inserting into the MessageChannel table. Please contact an administrator", http.StatusInternalServerError)
		fmt.Fprintln(os.Stderr, "Error inserting into the MessageChannel table:", err)
		return
	}

	// Insert into user message table
	stmt = UserMessage.INSERT(UserMessage.User, UserMessage.Message).VALUES(uid, destM.Message.ID).RETURNING(UserMessage.AllColumns)
	var destUM struct {
		model.UserMessage
	}
	err = stmt.Query(domain.Db, &destUM)
	if err != nil {
		http.Error(w, "Error inserting into the UserMessage table. Please contact an administrator", http.StatusInternalServerError)
		fmt.Fprintln(os.Stderr, "Error inserting into the UserMessage table:", err)
		return
	}
	fmt.Fprintf(w, "%d", destM.ID)
}
