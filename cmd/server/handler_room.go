package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ben-rw/webgame/internal/room"
)

// generate 4 letter code for room, set creating player as host, redirect user to room at /room/{roomID}
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
		roomID = room.GenerateRoomID()
		_, ok := cfg.RoomReg.Get(roomID)
		if !ok {
			break
		}
	}

	newRoom := &room.Room{
		ID:      roomID,
		Players: []*room.Player{},
		Scene:   "lobby",
	}

	cfg.RoomReg.Set(roomID, newRoom)

	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    name,
		Path:     "/",
		HttpOnly: true,
	})

	fmt.Printf("Creating new room at %v for %v\n", "/room/"+roomID, name)

	http.Redirect(w, r, "/room/"+roomID, http.StatusSeeOther)
}

// parses url for roomID, creates Player with Player.Host set to false, adds Player to room.Players, redirects to /room/{roomID}
func (cfg *config) handlerJoinRoom(w http.ResponseWriter, r *http.Request) {
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

	player := room.Player{
		Name:  name,
		Score: 0,
		Host:  false,
	}

	err = cfg.RoomReg.AppendPlayer(roomID, &player)
	if err != nil {
		http.Error(w, "couldn't add player to room", http.StatusInternalServerError)
		log.Printf("couldn't add player to room: %v", err)
	}

	// use player.Name instead of name
	// player.Name is mutated in place during .AppendPlayer call in case of a name collision
	fmt.Printf("%v is joining room %v\n", player.Name, roomID)

	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    player.Name,
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
		return
	}
}

func (cfg *config) handlerServeTestRoom(w http.ResponseWriter, r *http.Request) {
	roomID := "TEST"
	username := "tester"

	_, ok := cfg.RoomReg.Get(roomID)
	if !ok {
		newRoom := &room.Room{
			ID:      roomID,
			Players: []*room.Player{},
			Scene:   "lobby",
		}
		cfg.RoomReg.Set(roomID, newRoom)
		err := cfg.RoomReg.AppendPlayer(roomID, &room.Player{
			Name:  username,
			Score: 0,
			Host:  true,
		})
		if err != nil {
			log.Println(err)
		}
	} else {
		err := cfg.RoomReg.AppendPlayer(roomID, &room.Player{
			Name:  username,
			Score: 0,
			Host:  false,
		})
		if err != nil {
			log.Println(err)
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    username,
		Path:     "/",
		HttpOnly: true,
	})

	err := cfg.templates.ExecuteTemplate(w, "room.html", "TEST")
	if err != nil {
		http.Error(w, "unable to serve room page: couldn't execute template", http.StatusInternalServerError)
		log.Printf("unable to serve room page: couldn't execute template: %v\n", err)
		return
	}
}
