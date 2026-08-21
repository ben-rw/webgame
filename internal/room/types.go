package room

import (
	"github.com/coder/websocket"
)

type Player struct {
	ID          string
	Name        string
	Score       int
	Host        bool
	SpriteIndex int
	Room        *Room
	Client      *Client
}

type Client struct {
	*websocket.Conn
	Player *Player
}
