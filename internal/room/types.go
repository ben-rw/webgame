package room

import (
	"sync"

	"github.com/coder/websocket"
)

type Client struct {
	Name  string
	Score int
	Host  bool
	Room  *Room
	Conn  *websocket.Conn
}

type Room struct {
	ID      string
	Clients []Client
}

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

func (r *RoomRegistry) AppendClient(roomID string, client Client) {
	r.Mu.Lock()
	r.ActiveRooms[roomID].Clients = append(r.ActiveRooms[roomID].Clients, client)
	r.Mu.Unlock()
}
