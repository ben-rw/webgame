package main

import (
	_ "image/png"
	"log"

	"github.com/ben-rw/webgame/cmd/game/internal/scenes/lobby"
	"github.com/ben-rw/webgame/cmd/game/internal/ws"
	"github.com/ben-rw/webgame/internal/protocol"
	"github.com/hajimehoshi/ebiten/v2"
)

type Scene interface {
	Update(messages []protocol.Message) error
	Draw(screen *ebiten.Image)
}

type Game struct {
	Conn  *ws.Connection
	Scene Scene
}

func (g *Game) Update() error {
	return g.Scene.Update(g.Conn.Check())
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Scene.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 180
}

func main() {
	conn, err := ws.ConnectToWebsocket()
	if err != nil {
		log.Fatal(err)
	}

	game := &Game{
		Conn:  conn,
		Scene: lobby.New(),
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
