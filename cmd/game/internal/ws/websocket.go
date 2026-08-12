package main

import (
	"context"
	"log"
	"syscall/js"

	"github.com/ben-rw/webgame/internal/protocol"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type Connection struct {
	*websocket.Conn
	messages chan protocol.Message
}

// check for messages to pass to Update in main
func (c *Connection) Check() []protocol.Message {
	msgs := []protocol.Message{}
	for {
		select {
		case msg := <-c.messages:
			msgs = append(msgs, msg)
		default:
			return msgs
		}
	}
}

// read incoming messages, close dead connections, drop messages if buffer fills up
func (c *Connection) readLoop() {
	defer c.Close(websocket.StatusInternalError, "connection closed unexpectedly")

	msg := protocol.Message{}

	for {
		err := wsjson.Read(context.Background(), c.Conn, &msg)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
				websocket.CloseStatus(err) == websocket.StatusGoingAway {
				c.Close(websocket.StatusNormalClosure, "connection closed normally")
				return
			} else {
				log.Printf("wsjson read error: %v\n", err)
				log.Println("terminating read loop")
				return
			}
		}
		select {
		case c.messages <- msg:
			log.Printf("received msg: %v", msg)
		default:
			log.Println("buffer full: dropped a message")
		}
	}

}

// set up websocket connection, send roomID to server, start read loop
func ConnectToWebsocket() (*Connection, error) {
	ctx := context.Background()

	prot, host, roomID := getClientInfo()
	wsURL := prot + host + "/room/" + roomID + "/ws"

	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		return nil, err
	}

	joinData := protocol.JoinRequestData{
		RoomID: roomID,
	}

	msg, err := protocol.MarshalToMessage(protocol.JoinRequest, joinData)
	if err != nil {
		return nil, err
	}

	err = wsjson.Write(ctx, conn, msg)
	if err != nil {
		return nil, err
	}

	msgs := make(chan protocol.Message, 20)
	c := &Connection{
		conn,
		msgs,
	}

	go c.readLoop()

	return c, nil
}

// use syscall/js to get url info and roomID
func getClientInfo() (string, string, string) {
	loc := js.Global().Get("location")

	prot := "ws://"
	if loc.Get("protocol").String() == "https:" {
		prot = "wss://"
	}
	host := loc.Get("host").String()

	params := js.Global().Get("URLSearchParams").New(loc.Get("search"))
	roomID := params.Call("get", "room").String()

	return prot, host, roomID
}
