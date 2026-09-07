package shared

import "math"

type Camera struct {
	X, Y float64
}

func NewCamera(x, y float64) *Camera {
	return &Camera{
		X: x,
		Y: y,
	}
}

func (c *Camera) FollowTarget(targetX, targetY float64) {
	targetX += TileSize / 2
	targetY += TileSize / 2
	c.X = -targetX + ScreenWidth/2.0
	c.Y = -targetY + ScreenHeight/2.0
}

func (c *Camera) Constrain(tilemapWidthPixels, tilemapHeightPixels float64) {
	c.X = math.Min(c.X, 0.0)
	c.Y = math.Min(c.Y, 0.0)

	c.X = math.Max(c.X, ScreenWidth-tilemapWidthPixels)
	c.Y = math.Max(c.Y, ScreenHeight-tilemapHeightPixels)
}
