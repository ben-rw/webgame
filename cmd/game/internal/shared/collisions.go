package shared

import "image"

func CheckCollisionHorizontal(sprite *Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(image.Rect(
			int(sprite.X),
			int(sprite.Y),
			int(sprite.X)+16,
			int(sprite.Y)+16,
		)) {
			if sprite.Dx > 0.0 {
				sprite.X = float64(collider.Min.X) - TileSize
			} else if sprite.Dx < 0.0 {
				sprite.X = float64(collider.Max.X)
			}
		}
	}
}

func CheckCollisionVertical(sprite *Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(image.Rect(
			int(sprite.X),
			int(sprite.Y),
			int(sprite.X)+16,
			int(sprite.Y)+16,
		)) {
			if sprite.Dy > 0.0 {
				sprite.Y = float64(collider.Min.Y) - TileSize
			} else if sprite.Dy < 0.0 {
				sprite.Y = float64(collider.Max.Y)
			}
		}
	}
}
