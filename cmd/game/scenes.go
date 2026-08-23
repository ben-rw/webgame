package main

import (
	"github.com/ben-rw/webgame/cmd/game/internal/scenes/lobby"
	"github.com/ben-rw/webgame/cmd/game/internal/scenes/memory"
	"github.com/ben-rw/webgame/cmd/game/internal/ws"
	"github.com/ben-rw/webgame/internal/protocol"
	"log"
)

func StartNewScene(sceneType protocol.SceneType, c *ws.Connection) Scene {
	switch sceneType {
	case protocol.LobbyScene:
		return lobby.NewLobby(c)
	case protocol.MemoryScene:
		return memory.NewMemory(c)
	default:
		log.Println("invalid scene name")
		return nil
	}
}
