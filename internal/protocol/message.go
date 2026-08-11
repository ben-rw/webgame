package protocol

import (
	"encoding/json"
	"errors"
)

// "github.com/ben-rw/webgame/internal/room"

type MessageType string

const RoomState MessageType = "RoomState"

type Message struct {
	Type MessageType
	Data json.RawMessage
}

type RoomStateData struct {
	Username string `json:"username"`
	RoomID   string `json:"room_id"`
}

func (m *Message) UnmarshalMessageData() (any, error) {
	var v any

	switch m.Type {
	case RoomState:
		v = &RoomStateData{}
	default:
		return nil, errors.New("unrecognized message format")
	}

	err := json.Unmarshal(m.Data, v)
	if err != nil {
		return nil, err
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
