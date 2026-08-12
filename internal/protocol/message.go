package protocol

import (
	"encoding/json"
	"errors"
)

type MessageType string

const JoinRequest MessageType = "JoinRequest"
const JoinResponse MessageType = "JoinResponse"

type Message struct {
	Type MessageType
	Data json.RawMessage
}

type JoinRequestData struct {
	RoomID string `json:"room_id"`
}

type JoinResponseData struct {
	Username   string   `json:"username"`
	PlayerList []string `json:"player_list"`
}

func (m *Message) UnmarshalMessageData() (any, error) {
	var v any

	switch m.Type {
	case JoinRequest:
		v = &JoinRequestData{}
	case JoinResponse:
		v = &JoinResponseData{}
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
