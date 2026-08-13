package room

import (
	"github.com/coder/websocket"
)

type Player struct {
	Name   string
	Score  int
	Host   bool
	Room   *Room
	Client *Client
}

type Client struct {
	*websocket.Conn
	Player *Player
}
