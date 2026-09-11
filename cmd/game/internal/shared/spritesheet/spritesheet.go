package spritesheet

import (
	"image"
)

type SpriteSheet struct {
	WidthInTiles  int
	HeightInTiles int
	TileWidth     int
	TileHeight    int
}

func NewSpriteSheet(w, h, tw, th int) *SpriteSheet {
	return &SpriteSheet{
		w, h, tw, th,
	}
}

func (s *SpriteSheet) Rect(index int) image.Rectangle {
	x := (index % s.WidthInTiles) * s.TileWidth
	y := (index / s.WidthInTiles) * s.TileHeight

	return image.Rect(x, y, x+s.TileWidth, y+s.TileHeight)
}
