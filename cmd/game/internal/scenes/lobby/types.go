package lobby

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	Name  string
	Score int
	Host  bool
}

func (p *Player) New(username string) *Player {
	return &Player{
		Name: username,
	}
}

type Sprite struct {
	Img  *ebiten.Image
	X, Y float64
}
