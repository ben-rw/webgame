package main

import (
	_ "embed"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
)

//go:embed app/room/index.html
var roomHTMLTemplate string

type Room struct {
	ID      string
	Players []*Player
}

// generates 4 letter code for room, sets creating player as host, redirects user to room at /room/{roomID}
func (cfg *config) handlerCreateRoom(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "unable to parse form", err)
		return
	}

	name := r.FormValue("username")
	if name == "" {
		respondWithError(w, http.StatusBadRequest, "player must provide a username", nil)
		return
	}

	//gen roomID + ensure there isn't an active room with the same roomID
	var roomID string
	for {
		roomID = generateRoomID()

		cfg.mu.RLock()
		_, ok := cfg.activeRooms[roomID]
		cfg.mu.RUnlock()

		if !ok {
			break
		}
	}

	host := Player{
		Name:  name,
		Score: 0,
		Host:  true,
	}

	room := Room{
		ID:      roomID,
		Players: []*Player{&host},
	}

	cfg.mu.Lock()
	cfg.activeRooms[roomID] = &room
	cfg.mu.Unlock()

	newUrl := cfg.RootURL + "room/" + roomID
	fmt.Printf("Creating new room at %v for %v\n", newUrl, name)

	http.Redirect(w, r, newUrl, http.StatusSeeOther)
}

func generateRoomID() string {
	const codeLength = 4
	var charSet = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	buf := make([]rune, codeLength)
	for i := range buf {
		buf[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(buf)
}

// parses url for roomID, creates player with player.Host set to false, adds player to room.Players, redirects to /room/{roomID}
func (cfg *config) handlerJoinRoom(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "unable to parse form", err)
		return
	}

	name := r.FormValue("username")
	if name == "" {
		respondWithError(w, http.StatusBadRequest, "player must provide a username", nil)
		return
	}

	roomID := r.FormValue("roomID")
	roomID = strings.ToUpper(roomID)
	cfg.mu.RLock()
	_, ok := cfg.activeRooms[roomID]
	cfg.mu.RUnlock()
	if !ok {
		respondWithError(w, http.StatusBadRequest, "no active rooms with that room code", nil)
		return
	}

	roomUrl := cfg.RootURL + "room/" + roomID

	player := Player{
		Name:  name,
		Score: 0,
		Host:  false,
	}

	cfg.mu.Lock()
	cfg.activeRooms[roomID].Players = append(cfg.activeRooms[roomID].Players, &player)
	cfg.mu.Unlock()

	fmt.Printf("%v is joining room %v\n", name, roomID)

	http.Redirect(w, r, roomUrl, http.StatusSeeOther)
}

// parses url for roomID, serves room page html at /room/{roomID}
func (cgf *config) handlerServeRoomPage(w http.ResponseWriter, r *http.Request) {
	roomID := r.PathValue("roomID")

	roomHTML := fmt.Sprintf(roomHTMLTemplate, roomID)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(roomHTML))
}
