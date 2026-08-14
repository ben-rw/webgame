package lobby

import (
	"github.com/ben-rw/webgame/internal/protocol"
	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	*Sprite
	Data *protocol.PlayerData
}

type Sprite struct {
	Img  *ebiten.Image
	X, Y float64
}

func NewPlayer(data *protocol.PlayerData, img *ebiten.Image) *Player {
	return &Player{
		Sprite: &Sprite{
			Img: img,
			X:   50,
			Y:   50,
		},
		Data: &protocol.PlayerData{
			Name:  data.Name,
			Score: data.Score,
			Host:  data.Host,
		},
	}
}
