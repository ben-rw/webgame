package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ben-rw/webgame/internal/protocol"
	"github.com/ben-rw/webgame/internal/room"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

// get roomID + username, upgrade to websocket, register client to hub, check for reconnection,
// send username to client, start read loop, remove client when read loop returns
func (cfg *config) handlerWebsocket(w http.ResponseWriter, r *http.Request) {
	roomID := r.PathValue("roomID")
	_, ok := cfg.RoomReg.Get(roomID)
	if !ok {
		http.Error(w, "no active room with that room code", http.StatusBadRequest)
		return
	}

	userCookie, err := r.Cookie("username")
	if err != nil {
		http.Error(w, "missing user session data", http.StatusUnauthorized)
		log.Println(err)
		return
	}
	username := userCookie.Value

	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		http.Error(w, "couldn't upgrade connection to websocket", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	playerNames, err := getPlayerNames(cfg, roomID)
	if err != nil {
		log.Println(err)
	}

	var host bool
	if len(playerNames) == 0 {
		host = true
	}

	existingPlayer, err := checkForReconnect(cfg, roomID, username)
	if err != nil {
		log.Println(err)
	}

	if existingPlayer == nil {
		playerNames = append(playerNames, username)
	}

	data := protocol.JoinResponseData{
		Username:   username,
		PlayerList: playerNames,
	}

	msg, err := protocol.MarshalToMessage(protocol.JoinResponse, data)
	if err != nil {
		log.Println(err)
	}

	ctx := context.Background()
	err = wsjson.Write(ctx, conn, msg)
	if err != nil {
		log.Println(err)
		return
	}

	msgs := make(chan protocol.Message, 100)
	client := room.Client{
		Conn:     conn,
		Messages: msgs,
	}

	var player *room.Player
	if existingPlayer != nil {
		existingPlayer.Client = &client
		client.Player = existingPlayer
		player = existingPlayer
	} else {
		player = &room.Player{
			Name:   username,
			Score:  0,
			Host:   host,
			Client: &client,
		}
		err = cfg.RoomReg.AppendPlayer(roomID, player)
		if err != nil {
			log.Println(err)
		}
		client.Player = player
	}

	readLoop(&client)
	player.Client = nil
}

// read incoming messages, close dead connections, drop messages if buffer fills up
func readLoop(c *room.Client) {
	defer c.Close(websocket.StatusInternalError, "connection closed unexpectedly")

	msg := protocol.Message{}

	for {
		err := wsjson.Read(context.Background(), c.Conn, &msg)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
				websocket.CloseStatus(err) == websocket.StatusGoingAway {
				c.Close(websocket.StatusNormalClosure, "connection closed normally")
				return
			} else {
				log.Printf("wsjson read error: %v\n", err)
				log.Println("terminating read loop")
				return
			}
		}
		select {
		case c.Messages <- msg:
			log.Printf("received msg: %v", msg)
		default:
			log.Println("buffer full: dropped a message")
		}
	}

}

func getPlayerNames(cfg *config, roomID string) ([]string, error) {
	playerList, err := cfg.RoomReg.GetPlayerList(roomID)
	if err != nil {
		return []string{}, err
	}

	nameList := make([]string, len(playerList))
	for _, player := range playerList {
		nameList = append(nameList, player.Name)
	}

	return nameList, nil
}

func checkForReconnect(cfg *config, roomID, username string) (*room.Player, error) {
	playerList, err := cfg.RoomReg.GetPlayerList(roomID)
	if err != nil {
		return nil, err
	}

	for _, player := range playerList {
		if player.Name == username {
			return player, nil
		}
	}

	return nil, nil
}
