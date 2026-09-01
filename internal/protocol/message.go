package protocol

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

type MessageType string

const (
	Unset        MessageType = "Unset"
	JoinRequest  MessageType = "JoinRequest"
	JoinResponse MessageType = "JoinResponse"
	SceneChange  MessageType = "SceneChange"
	PlayerUpdate MessageType = "PlayerUpdate"
)

type Message struct {
	Type MessageType     `json:"type"`
	Data json.RawMessage `json:"data"`
}

type JoinRequestData struct {
	RoomID string `json:"room_id"`
}

type PlayerData struct {
	Name        string  `json:"username"`
	Score       int     `json:"score"`
	Host        bool    `json:"host"`
	SpriteIndex int     `json:"sprite_index"`
	X           float64 `json:"x"`
	Y           float64 `json:"y"`
}

type JoinResponseData struct {
	PlayerData *PlayerData   `json:"player_data"`
	PlayerList []*PlayerData `json:"player_list"`
}

type PlayerUpdateData struct {
	PlayerData *PlayerData `json:"player_data"`
}

type SceneType int

const (
	LobbyScene SceneType = iota
	RandomScene
	MemoryScene
)

type SceneChangeData struct {
	SceneType SceneType `json:"scene_type"`
}

func (m *Message) UnmarshalMessageData() (any, error) {
	var v any

	switch m.Type {
	case JoinRequest:
		v = &JoinRequestData{}
	case JoinResponse:
		v = &JoinResponseData{}
	case SceneChange:
		v = &SceneChangeData{}
	case PlayerUpdate:
		v = &PlayerUpdateData{}
	default:
		return nil, errors.New("unrecognized message format")
	}

	log.Printf("json: %s type: %v\n", m.Data, m.Type)

	err := json.Unmarshal(m.Data, v)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}

	return v, nil
}

func MarshalToMessage(mt MessageType, msgData any) (*Message, error) {
	jsonData, err := json.Marshal(msgData)
	if err != nil {
		return nil, err
	}

	msg := &Message{
		Type: mt,
		Data: jsonData,
	}

	return msg, nil
}
