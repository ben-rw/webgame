package minigames

import (
	// "errors"
	"github.com/ben-rw/webgame/internal/protocol"
	"github.com/ben-rw/webgame/internal/room"
	"log"
)

func ValidateMessage(msg *protocol.Message, scene room.SceneType) (*protocol.Message, error) {
	var data any
	if msg.Type == protocol.PlayerUpdate {
		log.Println("1")
		// playerUpdateData, err := msg.UnmarshalMessageData()
		// if err != nil {
		// 	return &protocol.Message{}, err
		// }
		// data, ok := playerUpdateData.(*protocol.PlayerUpdateData)
		// if !ok {
		// 	log.Println(errors.New("couldn't validate message: failed type assertion"))
		// 	return &protocol.Message{}, nil
		// }
		// if data.PlayerData.Name != "" {
		// log.Printf("validated msg data: %v, scene: %v", data, scene)
		return msg, nil
		// }
	}

	log.Printf("couldn't validate msg data: %v, scene: %v", data, scene)
	return &protocol.Message{Type: protocol.Unset}, nil
}
