package room

import (
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
