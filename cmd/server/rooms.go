package main

import (
	_ "embed"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
)

type Room struct {
	ID      string
	Players []*Player
}

// generates 4 letter code for room, sets creating player as host, redirects user to room at /room/{roomID}
func (cfg *config) handlerCreateRoom(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "unable to parse form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("username")
	if name == "" {
		err := cfg.templates.ExecuteTemplate(w, "landing.html", LandingPageError{CreateUsernameError: "player must provide a username"})
		if err != nil {
			http.Error(w, "unable to serve room page with error message", http.StatusInternalServerError)
		}
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

	fmt.Printf("Creating new room at %v for %v\n", "/room/"+roomID, name)

	http.Redirect(w, r, "/room/"+roomID, http.StatusSeeOther)
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

	//TODO - need to check that a player with the same name isn't already in the room
	//Scores need to be stored in a way that a player who disconnects and reconnects doesn't lose their points

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "unable to parse form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("username")
	if name == "" {
		err := cfg.templates.ExecuteTemplate(w, "landing.html", LandingPageError{JoinUsernameError: "player must provide a username"})
		if err != nil {
			http.Error(w, "unable to serve room page with error message", http.StatusInternalServerError)
		}
		return
	}

	roomID := r.FormValue("roomID")
	roomID = strings.ToUpper(roomID)
	cfg.mu.RLock()
	_, ok := cfg.activeRooms[roomID]
	cfg.mu.RUnlock()
	if !ok {
		err := cfg.templates.ExecuteTemplate(w, "landing.html", LandingPageError{JoinCodeError: "no active rooms with that room code"})
		if err != nil {
			http.Error(w, "unable to serve room page with error message", http.StatusInternalServerError)
		}
		return
	}

	player := Player{
		Name:  name,
		Score: 0,
		Host:  false,
	}

	cfg.mu.Lock()
	cfg.activeRooms[roomID].Players = append(cfg.activeRooms[roomID].Players, &player)
	cfg.mu.Unlock()

	fmt.Printf("%v is joining room %v\n", name, roomID)

	http.Redirect(w, r, "/room/"+roomID, http.StatusSeeOther)
}

// parses url for roomID, serves room page html at /room/{roomID}
func (cfg *config) handlerServeRoomPage(w http.ResponseWriter, r *http.Request) {
	roomID := r.PathValue("roomID")
	fmt.Printf("roomID: %v\n", roomID)

	roomID = strings.ToUpper(roomID)
	cfg.mu.RLock()
	_, ok := cfg.activeRooms[roomID]
	cfg.mu.RUnlock()
	if !ok {
		http.Error(w, "unable to serve room page: unregistered room code", http.StatusInternalServerError)
		return
	}

	err := cfg.templates.ExecuteTemplate(w, "room.html", roomID)
	if err != nil {
		http.Error(w, "unable to serve room page: couldn't execute template", http.StatusInternalServerError)
	}
}

func (cfg *config) handlerServeTestRoom(w http.ResponseWriter, r *http.Request) {
	err := cfg.templates.ExecuteTemplate(w, "room.html", nil)
	if err != nil {
		http.Error(w, "unable to serve room page: couldn't execute template", http.StatusInternalServerError)
	}
}
