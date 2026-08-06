package main

import (
	"context"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"log"
	"net/http"
	"time"
)

func (cfg *config) handlerWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		http.Error(w, "couldn't upgrade connection to websocket", http.StatusInternalServerError)
		return
	}
	defer conn.CloseNow()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var v any
	err = wsjson.Read(ctx, conn, &v)
	if err != nil {
		http.Error(w, "couldn't read websocket connection", http.StatusInternalServerError)
		return
	}

	log.Printf("received: %v", v)

	conn.Close(websocket.StatusNormalClosure, "")
}
