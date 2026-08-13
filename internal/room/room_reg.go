package room

import (
	"context"
	"errors"
	"github.com/ben-rw/webgame/internal/protocol"
	"github.com/coder/websocket/wsjson"
	"log"
	"sync"
)

type RoomRegistry struct {
	Mu          *sync.RWMutex
	ActiveRooms map[string]*Room
}

func (r *RoomRegistry) Get(roomID string) (*Room, bool) {
	r.Mu.RLock()
	val, ok := r.ActiveRooms[roomID]
	r.Mu.RUnlock()

	if !ok {
		return &Room{}, ok
	}

	return val, ok
}

func (r *RoomRegistry) Set(roomID string, room *Room) {
	r.Mu.Lock()
	r.ActiveRooms[roomID] = room
	r.Mu.Unlock()
}

// check room exists, check for username collisions, set player.Room, add playe to RoomReg
func (r *RoomRegistry) AppendPlayer(roomID string, player *Player) error {
	_, ok := r.Get(roomID)
	if !ok {
		return errors.New("room with that id doesn't exist")
	}

	r.Mu.Lock()
	player.Name = nameCollisionSolver(player.Name, r.ActiveRooms[roomID].Players)
	player.Room = r.ActiveRooms[roomID]
	r.ActiveRooms[roomID].Players = append(r.ActiveRooms[roomID].Players, player)
	r.Mu.Unlock()

	return nil
}

func (r *RoomRegistry) RemovePlayer(roomID string, currentPlayer *Player) error {
	_, ok := r.Get(roomID)
	if !ok {
		return errors.New("room with that id doesn't exist")
	}

	r.Mu.Lock()
	var i int
	var match bool
	for index, player := range r.ActiveRooms[roomID].Players {
		if player == currentPlayer {
			i = index
			match = true
			break
		}
	}
	if !match {
		return errors.New("player not found")
	}

	r.ActiveRooms[roomID].Players = append(r.ActiveRooms[roomID].Players[:i], r.ActiveRooms[roomID].Players[i+1:]...)
	r.Mu.Unlock()

	return nil
}

func (r *RoomRegistry) GetPlayerList(roomID string) ([]*Player, error) {
	_, ok := r.Get(roomID)
	if !ok {
		return []*Player{}, errors.New("room with that id doesn't exist")
	}

	r.Mu.RLock()
	playerList := make([]*Player, len(r.ActiveRooms[roomID].Players))
	_ = copy(playerList, r.ActiveRooms[roomID].Players)
	r.Mu.RUnlock()
	return playerList, nil
}

func (r *RoomRegistry) Broadcast(msg *protocol.Message, room *Room) {
	r.Mu.RLock()
	playerList := make([]*Player, len(room.Players))
	_ = copy(playerList, room.Players)
	r.Mu.RUnlock()

	for _, player := range playerList {
		if player.Client != nil {
			err := wsjson.Write(context.Background(), player.Client.Conn, msg)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
