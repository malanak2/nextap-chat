package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/gorilla/mux"
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
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	page := 1
	limit := 50
	if r.URL.Query()["limit"] != nil {
		limit, err = strconv.Atoi(r.URL.Query()["limit"][0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if r.URL.Query()["page"] != nil {
		page, err = strconv.Atoi(r.URL.Query()["page"][0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	stmt := UserMessage.SELECT(UserMessage.Message).WHERE(UserMessage.User.EQ(postgres.Int(int64(uid)))).LIMIT(int64(limit)).OFFSET(int64((page - 1) * limit))

	var destM []struct {
		model.UserMessage
	}

	err = stmt.Query(domain.Db, &destM)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result") {
			fmt.Fprintf(w, "No messages from user with the id of %d", uid)
			return
		}
		http.Error(w, "Error selecting from the UserMessage table. Please contact an administrator. "+err.Error(), http.StatusInternalServerError)
		return
	}
	marshal, err := json.Marshal(destM)
	if err != nil {
		http.Error(w, "Error converting  UserMessage table. Please contact an administrator. "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", marshal)
}

func HandleGetMessageById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	stmt := postgres.SELECT(Message.AllColumns, User.AllColumns).WHERE(Message.ID.EQ(postgres.Int(int64(id)))).FROM(Message.INNER_JOIN(UserMessage, Message.ID.EQ(UserMessage.Message)).INNER_JOIN(User, UserMessage.User.EQ(User.ID)))

	var destM struct {
		model.Message
		Author struct {
			model.User
		}
	}

	err = stmt.Query(domain.Db, &destM)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result") {
			fmt.Fprintf(w, "No messages with the id %d", id)
		}
		http.Error(w, "Error selecting from the Message table. Please contact an administrator. "+err.Error(), http.StatusInternalServerError)
	}

	marshal, err := json.Marshal(destM)
	if err != nil {
		http.Error(w, "Error converting  Message table. Please contact an administrator. "+err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "%s", marshal)
}

func HandleGetAllMessages(w http.ResponseWriter, r *http.Request) {
	page := 1
	limit := 50
	var err error
	if r.URL.Query()["limit"] != nil {
		limit, err = strconv.Atoi(r.URL.Query()["limit"][0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if r.URL.Query()["page"] != nil {
		page, err = strconv.Atoi(r.URL.Query()["page"][0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	stmt := Message.SELECT(Message.AllColumns).LIMIT(int64(limit)).OFFSET(int64((page - 1) * limit))
	var destM []struct {
		model.Message
	}
	err = stmt.Query(domain.Db, &destM)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	marshal, err := json.Marshal(destM)
	if err != nil {
		http.Error(w, "Error converting  Message table. Please contact an administrator.", http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "%s", marshal)
}

func HandleSearchMessages(w http.ResponseWriter, r *http.Request) {
	page := 1
	limit := 50
	var err error
	if r.URL.Query()["limit"] != nil {
		limit, err = strconv.Atoi(r.URL.Query()["limit"][0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if r.URL.Query()["page"] != nil {
		page, err = strconv.Atoi(r.URL.Query()["page"][0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	vars := mux.Vars(r)
	text, err := url.QueryUnescape(vars["txt"])
	if err != nil {
		http.Error(w, "Invalid text", http.StatusBadRequest)
		return
	}

	stmtSearch := Message.SELECT(Message.AllColumns).WHERE(postgres.LOWER(Message.Content).LIKE(postgres.LOWER(postgres.String("%" + text + "%")))).LIMIT(int64(limit)).OFFSET(int64((page - 1) * limit))
	var dest []struct {
		model.Message
	}
	fmt.Fprintf(os.Stdout, "Text: %s, Sql: %s", text, stmtSearch.DebugSql())
	err = stmtSearch.Query(domain.Db, &dest)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			fmt.Fprintf(w, `[]`)
			return
		}
		http.Error(w, `Database query error `+err.Error(), http.StatusInternalServerError)
	}
	marshal, err := json.Marshal(dest)
	if err != nil {
		http.Error(w, `Database marshal error `+err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "%s", marshal)
}

func HandleEditMessageById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	var body domain.EditMessage
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	stmt := Message.UPDATE(Message.Content).WHERE(Message.ID.EQ(postgres.Int(int64(id)))).SET(postgres.String(body.Content)).RETURNING(Message.AllColumns)
	var destM struct {
		model.Message
	}
	err = stmt.Query(domain.Db, &destM)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	marshal, err := json.Marshal(destM)
	if err != nil {
		http.Error(w, "Error converting  Message table. Please contact an administrator.", http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "%s", marshal)
}

func HandleDeleteMessageById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	stmtDUM := UserMessage.DELETE().WHERE(UserMessage.Message.EQ(postgres.Int(int64(id)))).RETURNING(UserMessage.AllColumns)
	var destDUM struct {
		model.UserMessage
	}
	err = stmtDUM.Query(domain.Db, &destDUM)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			http.Error(w, "No message with this id", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	stmt := Message.DELETE().WHERE(Message.ID.EQ(postgres.Int(int64(id)))).RETURNING(Message.AllColumns)
	var destM struct {
		model.Message
	}
	err = stmt.Query(domain.Db, &destM)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			http.Error(w, "No message with this id", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	marshal, err := json.Marshal(destM)
	if err != nil {
		http.Error(w, "Error converting  Message table. Please contact an administrator.", http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "%s", marshal)
}
