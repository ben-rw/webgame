package lobby

import (
	"github.com/ben-rw/webgame/internal/protocol"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"image/color"
	"log"
)

type Lobby struct {
	Players []*Player
	player  *Sprite
	sprites []*Sprite
}

func New() *Lobby {
	playerImg, _, err := ebitenutil.NewImageFromFileSystem(AssetsFS, "assets/images/ninja_adventure/Actor/Character/Inspector/SpriteSheet.png")
	if err != nil {
		log.Fatal(err)
	}
	return &Lobby{player: &Sprite{
		Img: playerImg,
		X:   50.0,
		Y:   50.0,
	}}
}

func (l *Lobby) Update(messages []protocol.Message) error {
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		l.player.X += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		l.player.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		l.player.Y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		l.player.Y += 2
	}
	return nil
}

func (l *Lobby) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})

	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(l.player.X, l.player.Y)

	screen.DrawImage(
		l.player.Img.SubImage(
			image.Rect(0, 0, 16, 16),
		).(*ebiten.Image),
		&opts,
	)

	opts.GeoM.Reset()

	for _, sprite := range l.sprites {
		opts.GeoM.Translate(sprite.X, sprite.Y)

		screen.DrawImage(
			sprite.Img.SubImage(
				image.Rect(0, 0, 16, 16),
			).(*ebiten.Image),
			&opts,
		)

		opts.GeoM.Reset()
	}
}
