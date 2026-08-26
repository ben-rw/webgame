package lobby

import (
	"log"

	"github.com/ben-rw/webgame/cmd/game/internal/shared"
	"github.com/ben-rw/webgame/internal/protocol"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	nameTagSize = 4
)

func NewPlayer(data *protocol.PlayerData, joinOrder int) *shared.Player {
	log.Printf("playerdata sprite index: %v", data.SpriteIndex)
	imgPath := PlayerSpriteIndex[data.SpriteIndex]
	playerImg, _, err := ebitenutil.NewImageFromFileSystem(shared.AssetsFS, imgPath)
	if err != nil {
		log.Fatal(err)
	}

	startPosition := StartingPositions[joinOrder]

	return &shared.Player{
		Sprite: &shared.Sprite{
			Img: playerImg,
			X:   startPosition.X,
			Y:   startPosition.Y,
		},
		Data: &protocol.PlayerData{
			Name:        data.Name,
			Score:       data.Score,
			Host:        data.Host,
			SpriteIndex: data.SpriteIndex,
		},
		NameTag: &shared.NameTag{
			Face: &text.GoTextFace{
				Source: shared.FontSrc,
				Size:   nameTagSize,
			},
			X: startPosition.X + shared.TileSize/2,
			Y: startPosition.Y + shared.TileSize + 2,
			LayoutOptions: text.LayoutOptions{
				PrimaryAlign: 1,
			},
		},
	}
}
