package room

import (
	"github.com/ben-rw/webgame/internal/protocol"
)

type Room struct {
	ID      string
	Players []*Player
	Scene   SceneType
}

func (r *Room) Broadcast(msg *protocol.Message) {}
