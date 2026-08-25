package minigames

import (
	"errors"
	"github.com/ben-rw/webgame/internal/minigames/lobby"
	"github.com/ben-rw/webgame/internal/protocol"
	"log"
	"math/rand"
)

var numberOfMinigames = 2

// lobby, random omitted as they are not games
var minigameMap = map[int]protocol.SceneType{
	2: protocol.MemoryScene,
	//placeholder
	3: protocol.LobbyScene,
}

func ValidateMessage(msg *protocol.Message, scene protocol.SceneType) (*protocol.Message, error) {
	if msg.Type == protocol.SceneChange {
		sceneChangeData, err := msg.UnmarshalMessageData()
		if err != nil {
			return &protocol.Message{Type: protocol.Unset}, err
		}
		scdata, ok := sceneChangeData.(*protocol.SceneChangeData)
		if !ok {
			err = errors.New("couldn't validate message: failed type assertion")
			return &protocol.Message{Type: protocol.Unset}, err
		}

		switch scdata.SceneType {
		case protocol.RandomScene:
			data := protocol.SceneChangeData{
				SceneType: minigameMap[rand.Intn(numberOfMinigames)+numberOfMinigames],
			}
			log.Printf("randomly selected scene: %v\n", data.SceneType)
			msg, err := protocol.MarshalToMessage(protocol.SceneChange, data)
			if err != nil {
				return &protocol.Message{Type: protocol.Unset}, err
			}
			return msg, nil
		case protocol.LobbyScene:
			data := protocol.SceneChangeData{
				SceneType: protocol.LobbyScene,
			}
			msg, err := protocol.MarshalToMessage(protocol.SceneChange, data)
			if err != nil {
				return &protocol.Message{Type: protocol.Unset}, err
			}
			return msg, nil
		case protocol.MemoryScene:
			data := protocol.SceneChangeData{
				SceneType: protocol.MemoryScene,
			}
			msg, err := protocol.MarshalToMessage(protocol.SceneChange, data)
			if err != nil {
				return &protocol.Message{Type: protocol.Unset}, err
			}
			return msg, nil
		}
	} else {
		switch scene {
		case protocol.LobbyScene:
			msg, err := lobby.ValidateMsg(msg)
			if err != nil {
				log.Printf("lobby couldn't validate msg: %v", err)
			}
			return msg, nil
		case protocol.MemoryScene:
			msg, err := lobby.ValidateMsg(msg)
			if err != nil {
				log.Printf("lobby couldn't validate msg: %v", err)
			}
			return msg, nil
		}
	}

	return &protocol.Message{Type: protocol.Unset}, nil
}
