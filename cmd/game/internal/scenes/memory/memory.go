package memory

import (
	"github.com/ben-rw/webgame/cmd/game/internal/ws"
	"github.com/ben-rw/webgame/internal/protocol"
	"github.com/hajimehoshi/ebiten/v2"
	//	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
)

type Memory struct {
	Conn *ws.Connection
	// Players map[string]*Player
	// Player  *Player
	// Sprites []*Sprite
}

func NewMemory(c *ws.Connection) *Memory {
	log.Println("scene changed to Memory")
	return &Memory{
		Conn: c,
		// Players: map[string]*Player{},
		// Player:  NewPlayer(&protocol.PlayerData{}, 0),
		// Sprites: []*Sprite{},
	}
}

func (m *Memory) Update(messages []protocol.Message) error { return nil }

func (l *Memory) Draw(screen *ebiten.Image) {}
