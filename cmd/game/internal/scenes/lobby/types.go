package lobby

import (
	"github.com/ben-rw/webgame/internal/protocol"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
)

type Player struct {
	*Sprite
	Data *protocol.PlayerData
}

type Sprite struct {
	Img  *ebiten.Image
	X, Y float64
}

func NewPlayer(data *protocol.PlayerData) *Player {
	imgPath := PlayerSpriteIndex[data.SpriteIndex]
	playerImg, _, err := ebitenutil.NewImageFromFileSystem(AssetsFS, imgPath)
	if err != nil {
		log.Fatal(err)
	}

	return &Player{
		Sprite: &Sprite{
			Img: playerImg,
			X:   50,
			Y:   50,
		},
		Data: &protocol.PlayerData{
			Name:        data.Name,
			Score:       data.Score,
			Host:        data.Host,
			SpriteIndex: data.SpriteIndex,
		},
	}
}
