package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ben-rw/webgame/internal/protocol"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

// get roomID + username, upgrade to websocket, register client to hub,
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

	data := protocol.JoinResponseData{
		Username:   userCookie.Value,
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
