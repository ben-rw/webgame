package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ben-rw/webgame/internal/room"
	"github.com/coder/websocket"
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
		roomID = room.GenerateRoomID()
		_, ok := cfg.RoomReg.Get(roomID)
		if !ok {
			break
		}
	}

	newRoom := room.Room{
		ID:      roomID,
		Clients: []room.Client{},
	}

	host := room.Client{
		Name:  name,
		Score: 0,
		Host:  true,
		Room:  &newRoom,
		Conn:  &websocket.Conn{},
	}

	cfg.RoomReg.Set(roomID, &newRoom)
	cfg.RoomReg.AppendClient(roomID, host)

	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    name,
		Path:     "/",
		HttpOnly: true,
	})

	fmt.Printf("Creating new room at %v for %v\n", "/room/"+roomID, name)

	http.Redirect(w, r, "/room/"+roomID, http.StatusSeeOther)
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
			log.Printf("unable to serve room page with error message: %v", err)
		}
		return
	}

	roomID := r.FormValue("roomID")
	roomID = strings.ToUpper(roomID)
	_, ok := cfg.RoomReg.Get(roomID)
	if !ok {
		err := cfg.templates.ExecuteTemplate(w, "landing.html", LandingPageError{JoinCodeError: "no active room with that room code"})
		if err != nil {
			http.Error(w, "unable to serve room page with error message", http.StatusInternalServerError)
			log.Printf("unable to serve room page with error message: %v\n", err)
		}
		return
	}

	client := room.Client{
		Name:  name,
		Score: 0,
		Host:  false,
	}

	cfg.RoomReg.AppendClient(roomID, client)

	fmt.Printf("%v is joining room %v\n", name, roomID)

	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    name,
		Path:     "/",
		HttpOnly: true,
	})

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
		log.Printf("unable to serve room page: couldn't execute template: %v\n", err)
	}
}

func (cfg *config) handlerServeTestRoom(w http.ResponseWriter, r *http.Request) {
	err := cfg.templates.ExecuteTemplate(w, "room.html", nil)
	if err != nil {
		http.Error(w, "unable to serve room page: couldn't execute template", http.StatusInternalServerError)
		log.Printf("unable to serve room page: couldn't execute template: %v\n", err)
	}
}
