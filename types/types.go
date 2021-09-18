package types

import "github.com/gocql/gocql"

type Message struct {
	ID      gocql.UUID `json:"id"`
	Message string     `json:"message"`
	Sender  string     `json:"sender"`
	Room    gocql.UUID `json:"room"`
}

type Room struct {
	ID   gocql.UUID `json:"id"`
	Name string     `json:"name"`
}
