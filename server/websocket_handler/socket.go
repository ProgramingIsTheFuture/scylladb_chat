package websocket_handler

import (
	"chat/db"
	"chat/types"
	"fmt"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/scylladb/gocqlx/v2/qb"
)

// Actions Types
const (
	SEND_MESSAGE    = "SEND"
	RECEIVE_MESSAGE = "RECEIVE"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type RequestMessage struct {
	Action  string        `json:"action"`
	Message types.Message `json:"message"`
}

type ServerMessages struct {
	Conn *websocket.Conn
	User string
}

type Server struct {
	rooms map[string][]ServerMessages
}

var server Server

func NewServerMessages() {
	server = Server{rooms: map[string][]ServerMessages{}}
	return
}

func (s *Server) AddUser(conn *websocket.Conn, room, user string) {
	s.rooms[room] = append(s.rooms[room], ServerMessages{Conn: conn, User: user})
}

func ChatRoomHandler(w http.ResponseWriter, r *http.Request) {
	room := mux.Vars(r)["room"]
	user := mux.Vars(r)["user"]
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	server.AddUser(conn, room, user)

	var reader RequestMessage
	var chat = make(chan types.Message)
	var valid = make(chan bool)
	var resp = make(chan types.Message)
	for {
		err := conn.ReadJSON(&reader)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if reader.Action != SEND_MESSAGE {
			continue
		}

		reader.Message.ID = gocql.TimeUUID()
		reader.Message.Room, err = gocql.ParseUUID(room)
		if err != nil {
			continue
		}

		go HandleMessages(reader.Message, chat, valid)

		go SaveMessage(chat, valid, resp)

		go Response(resp, room, user)
	}
}

func HandleMessages(msg types.Message, chat chan<- types.Message, valid chan<- bool) {
	if msg.Message == "" || msg.Sender == "" {
		valid <- false
	}

	valid <- true
	chat <- msg
}

func SaveMessage(chat <-chan types.Message, valid <-chan bool, resp chan<- types.Message) {
	if <-valid {
		msg := <-chat
		msgInsert := qb.Insert("pixelchart.message").Columns("id", "message", "sender", "room").Query(db.Session)
		msgInsert.BindStruct(msg)
		if err := msgInsert.ExecRelease(); err != nil {
			return
		}

		resp <- msg
	}
}

type IDResp struct {
	ID string `json:"id"`
}

func Response(resp <-chan types.Message, room, user string) {
	msg := <-resp

	for _, i := range server.rooms[room] {
		i.Conn.WriteJSON(msg)
		/*
				if i.User != user {
				} else {
			response := IDResp{ID: msg.ID.String()}
			i.Conn.WriteJSON(response)
				}*/
	}
}
