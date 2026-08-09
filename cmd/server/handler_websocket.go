package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ben-rw/webgame/internal/protocol"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func (cfg *config) handlerWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		http.Error(w, "couldn't upgrade connection to websocket", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer conn.CloseNow()

	ctx := context.Background()

	var msgFromClient protocol.JoinMessage
	err = wsjson.Read(ctx, conn, &msgFromClient)
	if err != nil {
		http.Error(w, "couldn't read from websocket connection", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	log.Printf("received: %v", msgFromClient)

	room, _ := cfg.RoomReg.Get(msgFromClient.RoomID)
	msgToClient := protocol.JoinResponse{
		Clients: room.Clients,
	}
	err = wsjson.Write(ctx, conn, &msgToClient)
	if err != nil {
		http.Error(w, "couldn't write to websocket connection", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	log.Printf("sent: %v", msgToClient)

	conn.Close(websocket.StatusNormalClosure, "")
}
