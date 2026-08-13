package minigames

import (
	"github.com/ben-rw/webgame/internal/protocol"
	"github.com/ben-rw/webgame/internal/room"
	"log"
)

func ValidateMessage(msg *protocol.Message, scene room.SceneType) (*protocol.Message, error) {
	data, err := msg.UnmarshalMessageData()
	if err != nil {
		return &protocol.Message{}, err
	}
	log.Printf("validated msg data: %v, scene: %v", data, scene)
	return msg, nil
}
