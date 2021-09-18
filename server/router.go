package server

import (
	"chat/server/http_handler"
	"chat/server/websocket_handler"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func Router() {
	r := mux.NewRouter()

	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	websocket_handler.NewServerMessages()

	r.HandleFunc("/get-room/", http_handler.GetRoom)
	r.HandleFunc("/create-room/", http_handler.CreateRoom)

	r.HandleFunc("/get-messages/{room}", http_handler.GetMessagesFromRoom).Methods("GET")

	r.HandleFunc("/chat/{room}/{user}", websocket_handler.ChatRoomHandler)

	http.ListenAndServe("0.0.0.0:8000", handlers.CORS(allowedOrigins, allowedMethods)(r))
}
