package main

import (
	"github.com/ben-rw/webgame/cmd/game/internal/scenes/lobby"
	"github.com/ben-rw/webgame/cmd/game/internal/scenes/memory"
	"github.com/ben-rw/webgame/cmd/game/internal/ws"
	"log"
)

func StartNewScene(scene string, c *ws.Connection) Scene {
	switch scene {
	case "lobby":
		return lobby.New(c)
	case "memory":
		return memory.New(c)
	default:
		log.Println("invalid scene name")
		return nil
	}
}
