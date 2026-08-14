package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ben-rw/webgame/internal/minigames"
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

	client := room.Client{
		Conn: conn,
	}

	playerDataList, err := getPlayerDataList(cfg, roomID)
	if err != nil {
		log.Println(err)
	}

	var host bool
	if len(playerDataList) == 0 {
		host = true
	}

	existingPlayer, exists, err := checkForReconnect(cfg, roomID, username)
	if err != nil {
		log.Println(err)
	}

	var player *room.Player
	if exists {
		existingPlayer.Client = &client
		client.Player = existingPlayer
		player = existingPlayer
	} else {
		player = &room.Player{
			Name:   username,
			Score:  0,
			Host:   host,
			Room:   cfg.RoomReg.ActiveRooms[roomID],
			Client: &client,
		}
		err = cfg.RoomReg.AppendPlayer(roomID, player)
		if err != nil {
			log.Println(err)
		}
		client.Player = player
	}

	playerData := playerToPlayerData(player)

	if !exists {
		playerDataList = append(playerDataList, playerData)
	}

	data := protocol.JoinResponseData{
		PlayerData: playerData,
		PlayerList: playerDataList,
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

	readLoop(&cfg.RoomReg, &client)
	player.Client = nil
}

// read incoming messages, close dead connections, drop messages if buffer fills up
func readLoop(r *room.RoomRegistry, c *room.Client) {
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
				log.Printf("server: wsjson read error: %v\n", err)
				log.Println("server: terminating read loop")
				return
			}
		}

		updateMsg, err := minigames.ValidateMessage(&msg, c.Player.Room.Scene)
		if err != nil {
			log.Println(err)
			continue
		} else {
			r.Broadcast(updateMsg, c.Player.Room)
		}
	}

}

func getPlayerDataList(cfg *config, roomID string) ([]*protocol.PlayerData, error) {
	playerList, err := cfg.RoomReg.GetPlayerList(roomID)
	if err != nil {
		return nil, err
	}

	dataList := make([]*protocol.PlayerData, len(playerList))
	for i, player := range playerList {
		dataList[i] = playerToPlayerData(player)
	}

	return dataList, nil
}

func playerToPlayerData(player *room.Player) *protocol.PlayerData {
	return &protocol.PlayerData{
		Name:  player.Name,
		Score: player.Score,
		Host:  player.Host,
	}
}

func checkForReconnect(cfg *config, roomID, username string) (*room.Player, bool, error) {
	playerList, err := cfg.RoomReg.GetPlayerList(roomID)
	if err != nil {
		return nil, false, err
	}

	for _, player := range playerList {
		if player.Name == username {
			return player, true, nil
		}
	}

	return nil, false, nil
}
