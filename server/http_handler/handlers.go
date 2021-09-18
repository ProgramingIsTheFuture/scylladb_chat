package http_handler

import (
	"chat/db"
	"chat/types"
	"encoding/json"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/scylladb/gocqlx/v2/qb"
)

type Error struct {
	Message string `json:"message"`
}

func sendJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	b, _ := json.Marshal(v)
	w.Write(b)
}

func GetRoom(w http.ResponseWriter, r *http.Request) {
	var rooms = []types.Room{}
	selectRooms := qb.Select("pixelchart.room").Columns("id", "name").Query(db.Session)
	err := selectRooms.Select(&rooms)
	if err != nil {
		sendJSON(w, 500, Error{Message: "Internal Server Error"})
		return
	}

	sendJSON(w, 200, rooms)

}

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	var room types.Room

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&room)
	if err != nil {
		sendJSON(w, 400, Error{Message: "Bad Request"})
		return
	}

	if room.Name == "" {
		sendJSON(w, 400, Error{Message: "Bad Request"})
		return
	}

	room.ID = gocql.TimeUUID()

	insertRoom := qb.Insert("pixelchart.room").Columns("id", "name").Query(db.Session)
	insertRoom.BindStruct(room)

	if err := insertRoom.ExecRelease(); err != nil {
		sendJSON(w, 500, Error{Message: "Internal Server Error"})
		return
	}

	w.Write([]byte(room.ID.String()))

	return
}

func GetMessagesFromRoom(w http.ResponseWriter, r *http.Request) {
	room := mux.Vars(r)["room"]

	var msg = []types.Message{}
	q := qb.Select("pixelchart.message").AllowFiltering().Where(qb.EqLit("room", room)).Columns("id", "message", "sender", "room").Query(db.Session)
	err := q.Select(&msg)
	if err != nil {
		sendJSON(w, 500, Error{Message: "Internal Server Error"})
		return
	}

	sendJSON(w, 200, msg)
}
