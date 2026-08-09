package protocol

import (
	"github.com/ben-rw/webgame/internal/room"
)

type JoinMessage struct {
	Name   string
	RoomID string
}

type JoinResponse struct {
	Clients []room.Client
}
