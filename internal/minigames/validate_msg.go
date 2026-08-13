package minigames

import (
	"github.com/ben-rw/webgame/internal/protocol"
	"github.com/ben-rw/webgame/internal/room"
	"log"
)

func ValidateMessage(msg *protocol.Message, scene room.SceneType) (*protocol.Message, error) {
	log.Printf("validated msg: %v, scene: %v", msg, scene)
	return msg, nil
}
