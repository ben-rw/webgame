package lobby

import (
	"github.com/ben-rw/webgame/cmd/game/internal/ws"
	"github.com/ben-rw/webgame/internal/protocol"
	"github.com/hajimehoshi/ebiten/v2"

	//	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"image/color"
	"log"
)

type Lobby struct {
	Conn    *ws.Connection
	Players map[string]*Player
	Player  *Player
	Sprites []*Sprite
}

func New(c *ws.Connection) *Lobby {

	return &Lobby{
		Conn:    c,
		Players: map[string]*Player{},
		Player:  NewPlayer(&protocol.PlayerData{}, 0),
		Sprites: []*Sprite{},
	}
}

func (l *Lobby) Update(messages []protocol.Message) error {
	for _, message := range messages {
		switch message.Type {
		case protocol.JoinResponse:
			joinResponseData, err := message.UnmarshalMessageData()
			if err != nil {
				log.Println(err)
				continue
			}
			data, ok := joinResponseData.(*protocol.JoinResponseData)
			if !ok {
				log.Println("failed type assertion")
				continue
			}

			for _, playerData := range data.PlayerList {
				if _, ok := l.Players[playerData.Name]; ok {
					l.Players[playerData.Name].Data = playerData
				} else {
					player := NewPlayer(playerData, len(l.Players))
					l.Players[playerData.Name] = player
				}
			}

			l.Player = l.Players[data.PlayerData.Name]

		case protocol.PlayerUpdate:
			playerUpdateData, err := message.UnmarshalMessageData()
			if err != nil {
				log.Println(err)
				continue
			}
			data, ok := playerUpdateData.(*protocol.PlayerUpdateData)
			if !ok {
				log.Println("failed type assertion")
				continue
			}
			l.Player = l.Players[data.PlayerData.Name]

		default:
		}
	}

	// if ebiten.IsKeyPressed(ebiten.KeyRight) {
	// 	l.Player.X += 2
	// }
	// if ebiten.IsKeyPressed(ebiten.KeyLeft) {
	// 	l.Player.X -= 2
	// }
	// if ebiten.IsKeyPressed(ebiten.KeyUp) {
	// 	l.Player.Y -= 2
	// }
	// if ebiten.IsKeyPressed(ebiten.KeyDown) {
	// 	l.Player.Y += 2
	// }

	// playerUpdate := &protocol.PlayerUpdateData{
	// PlayerData: l.Player.Data,
	// }
	// l.Conn.WriteMsg(protocol.PlayerUpdate, playerUpdate)

	return nil
}

func (l *Lobby) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})

	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(l.Player.X, l.Player.Y)

	screen.DrawImage(
		l.Player.Img.SubImage(
			image.Rect(0, 0, 16, 16),
		).(*ebiten.Image),
		&opts,
	)

	opts.GeoM.Reset()

	for _, sprite := range l.Players {
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
