package shared

import (
	"github.com/ben-rw/webgame/internal/protocol"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Player struct {
	*Sprite
	Data    *protocol.PlayerData
	NameTag *NameTag
}

type Sprite struct {
	Img  *ebiten.Image
	X, Y float64
}

type NameTag struct {
	Face          *text.GoTextFace
	X, Y          float64
	LayoutOptions text.LayoutOptions
}
