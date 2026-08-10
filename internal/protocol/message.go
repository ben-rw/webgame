package protocol

import (
// "github.com/ben-rw/webgame/internal/room"
)

type JoinMessage struct {
	RoomID string
}

type JoinResponse struct {
	Username string
}
