package shared

import (
	"log"

	"github.com/ben-rw/webgame/cmd/game/internal/shared/animations"
	"github.com/ben-rw/webgame/cmd/game/internal/shared/spritesheet"
	"github.com/ben-rw/webgame/internal/protocol"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	nameTagSize = 4
)

type Sprite struct {
	Img          *ebiten.Image
	X, Y, Dx, Dy float64
}

type Player struct {
	*Sprite
	Data            *protocol.PlayerData
	NameTag         *NameTag
	SpriteSheet     *spritesheet.SpriteSheet
	Animations      map[PlayerState]*animations.Animation
	ActiveAnimation *animations.Animation
	JustJoined      bool
}

type NameTag struct {
	Face          *text.GoTextFace
	X, Y          float64
	LayoutOptions text.LayoutOptions
}

type PlayerState int

const (
	Idle PlayerState = iota
	Down
	Up
	Left
	Right
	Join
)

// walks player from one place to another at speed proportionate
// to the initial distance between the player's location and the
// destination with no input from the player
func (p *Player) ScriptedWalk(destX, destY float64) {
	distX := destX - p.X
	stepX := distX / 100000
	distY := destY - p.Y
	stepY := distY / 100000
	for destX != p.X || destY != p.Y {
		if destX-p.X > stepX {
			p.X = destX
		} else {
			p.X += stepX
		}
		if destY-p.Y > stepY {
			p.Y = destY
		} else {
			p.Y += stepY
		}
		log.Println(p.X, p.Y)
	}
}

func (p *Player) GetActiveAnimation() *animations.Animation {
	if p.JustJoined {
		p.ActiveAnimation = p.Animations[Join]
		if p.ActiveAnimation.Over == true {
			p.JustJoined = false
		} else {
			return p.ActiveAnimation
		}
	}
	if p.Dx > 0 {
		return p.Animations[Right]
	}
	if p.Dx < 0 {
		return p.Animations[Left]
	}
	if p.Dy > 0 {
		return p.Animations[Down]
	}
	if p.Dx < 0 {
		return p.Animations[Up]
	}
	return p.Animations[Idle]
}

func NewPlayer(data *protocol.PlayerData, joinOrder int) *Player {
	log.Printf("playerdata sprite index: %v", data.SpriteIndex)
	imgPath := PlayerSpriteIndex[data.SpriteIndex]
	playerImg, _, err := ebitenutil.NewImageFromFileSystem(AssetsFS, imgPath)
	if err != nil {
		log.Fatal(err)
	}

	startPosition := StartingPositions[joinOrder]

	return &Player{
		Sprite: &Sprite{
			Img: playerImg,
			X:   startPosition.X,
			Y:   startPosition.Y,
			Dx:  0,
			Dy:  0,
		},
		Data: &protocol.PlayerData{
			Name:        data.Name,
			Score:       data.Score,
			Host:        data.Host,
			SpriteIndex: data.SpriteIndex,
		},
		NameTag: &NameTag{
			Face: &text.GoTextFace{
				Source: FontSrc,
				Size:   nameTagSize,
			},
			X: startPosition.X + TileSize/2,
			Y: startPosition.Y + TileSize + 2,
			LayoutOptions: text.LayoutOptions{
				PrimaryAlign: 1,
			},
		},
		SpriteSheet: spritesheet.NewSpriteSheet(4, 7, TileSize),
		Animations: map[PlayerState]*animations.Animation{
			Up:    animations.NewAnimation(5, 13, 4, 20.0),
			Down:  animations.NewAnimation(4, 12, 4, 20.0),
			Left:  animations.NewAnimation(6, 14, 4, 20.0),
			Right: animations.NewAnimation(7, 15, 4, 20.0),
			Idle:  animations.NewAnimation(0, 16, 16, 20.0),
			Join:  animations.NewAnimation(26, 27, 1, 60.0),
		},
		JustJoined: true,
	}
}
