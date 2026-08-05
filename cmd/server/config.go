package main

import (
	"html/template"
	"sync"
)

type Client struct {
	Name  string
	Score int
	Host  bool
}

type Room struct {
	ID      string
	Clients []*Client
}

type RoomRegistry struct {
	mu          *sync.RWMutex
	activeRooms map[string]*Room
}

type config struct {
	Port         string
	FilepathRoot string
	URLRoot      string
	RoomReg      RoomRegistry
	templates    *template.Template
}

func (r *RoomRegistry) Get(roomID string) (*Room, bool) {
	r.mu.RLock()
	val, ok := r.activeRooms[roomID]
	r.mu.RUnlock()

	if !ok {
		return &Room{}, ok
	}

	return val, ok
}

func (r *RoomRegistry) Set(roomID string, room *Room) {
	r.mu.Lock()
	r.activeRooms[roomID] = room
	r.mu.Unlock()
}

func (r *RoomRegistry) AppendClient(roomID string, client *Client) {
	r.mu.Lock()
	r.activeRooms[roomID].Clients = append(r.activeRooms[roomID].Clients, client)
	r.mu.Unlock()
}
