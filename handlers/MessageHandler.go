package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/malanak2/nextap-chat/domain"
	"github.com/malanak2/nextap-chat/ports"
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
		slog.Error("Invalid jwt detected")
		http.Error(w, ErrorBadJwt.Error(), http.StatusInternalServerError)
		return
	}
	var body domain.SendMessage
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, ErrorInvalidBody.Error(), http.StatusBadRequest)
		return
	}

	uExists, err := ports.UserExists(uid)

	if !uExists {
		http.Error(w, ErrorNoExist.Error(), http.StatusBadRequest)
		return
	}

	// Verify channel message
	if len(body.Content) == 0 {
		http.Error(w, ErrorTooShort.Error(), http.StatusBadRequest)
		return
	}
	if len(body.Content) > 1000 {
		http.Error(w, ErrorTooLong.Error(), http.StatusBadRequest)
		return
	}

	// Insert message

	// If in spec later handle it
	// Insert into message channel table
	/*stmt = MessageChannel.INSERT(MessageChannel.Channel, MessageChannel.Message).VALUES(body.Channel, destM.ID).RETURNING(MessageChannel.ID)
	var destMC struct {
		model.Message
	}
	err = stmt.Query(domain.Db, &destMC)
	if err != nil {
		http.Error(w, "Error inserting into the MessageChannel table. Please contact an administrator", http.StatusInternalServerError)
		fmt.Fprintln(os.Stderr, "Error inserting into the MessageChannel table:", err)
		return
	}*/
	msg, err := ports.SendMessage(body.Content, uid)
	if err != nil {
		http.Error(w, "Failed to send a message", http.StatusBadRequest)
		return
	}
	marshal, err := json.Marshal(msg)
	if err != nil {
		slog.Error("Failed to marshal message")
		http.Error(w, ErrorMarshal.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", marshal)
}

// HandleGetMessagesByUserId godoc
//
// @Summary			 	Messages by user
// @Description		 	Get all mesages by a user
// @Tags				message,user
// @Accept				json
// @Produce				json
// @Success				200 {object}	[]model.Message
// @Failure				400 {object}	string
// @Failure				500 {object}	string
// @Router				/user/{id}/messages [post]
func HandleGetMessagesByUserId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, ErrorInvalidId.Error(), http.StatusBadRequest)
		return
	}
	limit, page, err := GetVars(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	destM, err := ports.SelectMessagesByUserId(uid, limit, page)
	marshal, err := json.Marshal(destM)
	if err != nil {
		slog.Error("Failed to marshal Message", "error", destM)
		http.Error(w, "Failed to marshal result.", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", marshal)
}

// HandleGetMessageById godoc
//
// @Summary		 		Get a messages
// @Description 		Returns a
// @Tags				message
// @Accept				json
// @Produce				json
// @Success				200 {object}	[]model.Message
// @Failure				400 {object}	string
// @Failure				500 {object}	string
// @Security			JWTTokenBasic
// @Router				/message/{id} [get]
func HandleGetMessageById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, ErrorInvalidId.Error(), http.StatusBadRequest)
		return
	}

	destM, err := ports.SelectMessageById(id)
	if err != nil {
		http.Error(w, "Message does not exist", http.StatusBadRequest)
		return
	}
	marshal, err := json.Marshal(destM)
	if err != nil {
		slog.Error("Failed to marshal message", "error", err)
		http.Error(w, ErrorMarshal.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", marshal)
}

// HandleGetAllMessages godoc
//
// @Summary		 		Get all messages
// @Description 		Returns all messages sent on server
// @Tags				message
// @Accept				json
// @Produce				json
// @Success				200 {object}	[]model.Message
// @Failure				400 {object}	string
// @Failure				500 {object}	string
// @Security			JWTTokenBasic
// @Router				/messages [get]
func HandleGetAllMessages(w http.ResponseWriter, r *http.Request) {
	limit, page, err := GetVars(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	destM, err := ports.GetAllMessages(limit, page)
	if err != nil {
		http.Error(w, "Failed to get all messages", http.StatusInternalServerError)
		return
	}
	marshal, err := json.Marshal(destM)
	if err != nil {
		http.Error(w, "Error converting  Message table. Please contact an administrator.", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", marshal)
}

// HandleSearchMessages godoc
//
// @Summary		 		Search messages
// @Description 		Search for a message by id
// @Tags				message
// @Accept				json
// @Produce				json
// @Success				200 {object}	model.Message
// @Failure				400 {object}	string
// @Failure				500 {object}	string
// @Security			JWTTokenBasic
// @Router				/message/search/{txt} [post]
func HandleSearchMessages(w http.ResponseWriter, r *http.Request) {
	limit, page, err := GetVars(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	text, err := url.QueryUnescape(vars["txt"])
	if err != nil {
		http.Error(w, "Invalid text", http.StatusBadRequest)
		return
	}

	msgs, err := ports.SelectMessagesByContent(text, limit, page)
	marshal, err := json.Marshal(msgs)
	if err != nil {
		http.Error(w, `Failed to marshal message[]`, http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", marshal)
}

// HandleEditMessageById godoc
//
// @Summary		 		Edit a message
// @Description 		Edits a message by id
// @Tags				message
// @Accept				json
// @Produce				json
// @Success				200 {object}	model.Message
// @Failure				400 {object}	string
// @Failure				500 {object}	string
// @Security			JWTTokenBasic
// @Router				/message/{id}/update [post]
func HandleEditMessageById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, ErrorInvalidId.Error(), http.StatusBadRequest)
		return
	}
	var body domain.EditMessage
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, ErrorInvalidBody.Error(), http.StatusBadRequest)
		return
	}
	// Verify channel message
	if len(body.Content) == 0 {
		http.Error(w, ErrorTooShort.Error(), http.StatusBadRequest)
		return
	}
	if len(body.Content) > 1000 {
		http.Error(w, ErrorTooLong.Error(), http.StatusBadRequest)
		return
	}
	msgExists, err := ports.MessageExists(id)
	if !msgExists {
		if err != nil {
			http.Error(w, "Message does not exist", http.StatusBadRequest)
			return
		}
		http.Error(w, "Message does not exist", http.StatusBadRequest)
		return
	}
	destM, err := ports.UpdateMessageById(id, body.Content)
	if err != nil {
		http.Error(w, "Message not found", http.StatusBadRequest)
		return
	}
	marshal, err := json.Marshal(destM)
	if err != nil {
		http.Error(w, ErrorMarshal.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", marshal)
}

// HandleDeleteMessageById godoc
//
// @Summary		 		Deletes a message
// @Description 		Delete a message by id
// @Tags				message
// @Accept				json
// @Produce				json
// @Success				200
// @Failure				400 {object}	string
// @Failure				500 {object}	string
// @Security			JWTTokenBasic
// @Router				/message [delete]
func HandleDeleteMessageById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, ErrorInvalidId.Error(), http.StatusBadRequest)
		return
	}
	err = ports.DeleteMessage(id)
	if err != nil {
		http.Error(w, "Message not found.", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Message deleted")
}
