package main

import (
	"context"
	"log"
	"syscall/js"
	"time"

	"github.com/ben-rw/webgame/internal/protocol"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func connectToWebsocket() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	wsURL, roomID, username := getClientInfo()

	log.Println(wsURL, roomID, username)

	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		return err
	}
	defer conn.CloseNow()

	msg := protocol.JoinMessage{
		RoomID: roomID,
		Name:   username,
	}

	err = wsjson.Write(ctx, conn, msg)
	if err != nil {
		return err
	}

	return nil
}

// returns websocketURL, roomID, username
func getClientInfo() (string, string, string) {
	loc := js.Global().Get("location")

	protocol := "ws://"
	if loc.Get("protocol").String() == "https:" {
		protocol = "wss://"
	}
	host := loc.Get("host").String()

	params := js.Global().Get("URLSearchParams").New(loc.Get("search"))
	roomID := params.Call("get", "room").String()
	username := params.Call("get", "name").String()

	return protocol + host + "/room/" + roomID + "/ws", roomID, username
}
