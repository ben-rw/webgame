package main

import (
	_ "embed"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
)

// generates 4 letter code for room, sets creating client as host, redirects user to room at /room/{roomID}
func (cfg *config) handlerCreateRoom(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "unable to parse form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("username")
	if name == "" {
		err := cfg.templates.ExecuteTemplate(w, "landing.html", LandingPageError{CreateUsernameError: "client must provide a username"})
		if err != nil {
			http.Error(w, "unable to serve room page with error message", http.StatusInternalServerError)
		}
		return
	}

	//gen roomID + ensure there isn't an active room with the same roomID
	var roomID string
	for {
		roomID = generateRoomID()
		_, ok := cfg.RoomReg.Get(roomID)
		if !ok {
			break
		}
	}

	host := Client{
		Name:  name,
		Score: 0,
		Host:  true,
	}

	room := Room{
		ID:      roomID,
		Clients: []*Client{&host},
	}

	cfg.RoomReg.Set(roomID, &room)

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

// parses url for roomID, creates client with client.Host set to false, adds client to room.Clients, redirects to /room/{roomID}
func (cfg *config) handlerJoinRoom(w http.ResponseWriter, r *http.Request) {

	//TODO - need to check that a client with the same name isn't already in the room
	//Scores need to be stored in a way that a client who disconnects and reconnects doesn't lose their points

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
	_, ok := cfg.RoomReg.Get(roomID)
	if !ok {
		err := cfg.templates.ExecuteTemplate(w, "landing.html", LandingPageError{JoinCodeError: "no active rooms with that room code"})
		if err != nil {
			http.Error(w, "unable to serve room page with error message", http.StatusInternalServerError)
		}
		return
	}

	client := Client{
		Name:  name,
		Score: 0,
		Host:  false,
	}

	cfg.RoomReg.AppendClient(roomID, &client)

	fmt.Printf("%v is joining room %v\n", name, roomID)

	http.Redirect(w, r, "/room/"+roomID, http.StatusSeeOther)
}

// parses url for roomID, serves room page html at /room/{roomID}
func (cfg *config) handlerServeRoomPage(w http.ResponseWriter, r *http.Request) {
	roomID := r.PathValue("roomID")
	fmt.Printf("roomID: %v\n", roomID)

	roomID = strings.ToUpper(roomID)
	_, ok := cfg.RoomReg.Get(roomID)
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
