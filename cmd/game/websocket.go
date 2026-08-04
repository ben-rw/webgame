package main

import (
	"context"
	"syscall/js"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func connectToWebsocket() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	wsURL := getWebsocketURL()

	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		return err
	}
	defer conn.CloseNow()

	err = wsjson.Write(ctx, conn, "hi")

	return nil
}

func getWebsocketURL() string {
	loc := js.Global().Get("location")

	protocol := "ws://"
	if loc.Get("protocol").String() == "https:" {
		protocol = "wss://"
	}
	host := loc.Get("host").String()

	params := js.Global().Get("URLSearchParams").New(loc.Get("search"))
	roomID := params.Call("get", "room").String()

	return protocol + host + "/room/" + roomID + "/ws"
}
