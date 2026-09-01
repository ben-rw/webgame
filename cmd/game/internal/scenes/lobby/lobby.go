package lobby

import (
	"github.com/ben-rw/webgame/cmd/game/internal/shared"
	"github.com/ben-rw/webgame/cmd/game/internal/shared/screenproperties"
	"github.com/ben-rw/webgame/cmd/game/internal/ws"
	"github.com/ben-rw/webgame/internal/protocol"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"log"
)

type Lobby struct {
	shared.Roster
	Conn    *ws.Connection
	Sprites []*shared.Sprite
}

func NewLobby(c *ws.Connection) *Lobby {
	return &Lobby{
		Roster: shared.Roster{
			Players: map[string]*shared.Player{},
			Player:  shared.NewPlayer(&protocol.PlayerData{}, 0),
		},
		Conn:    c,
		Sprites: []*shared.Sprite{},
	}
}

func (l *Lobby) Update(messages []protocol.Message) error {
	for _, message := range messages {
		log.Printf("msg type: %v", string(message.Type))
		switch message.Type {
		case protocol.JoinResponse:
			err := l.HandleJoinResponse(message)
			if err != nil {
				log.Println(err)
				continue
			}

			log.Printf("player: %v", *l.Player.Data)

			playerUpdateData := protocol.PlayerUpdateData{
				PlayerData: l.Player.Data,
			}

			l.Conn.WriteMsg(protocol.PlayerUpdate, playerUpdateData)

		case protocol.PlayerUpdate:
			err := l.HandlePlayerUpdate(message)
			if err != nil {
				log.Println(err)
				continue
			}

		default:
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) &&
		l.Player.Data.Host == true {
		l.Conn.WriteMsg(protocol.SceneChange, protocol.SceneChangeData{
			SceneType: protocol.RandomScene,
		})
	}

	for _, player := range l.Players {
		player.ActiveAnimation = player.GetActiveAnimation()
		player.ActiveAnimation.Update()
	}

	return nil
}

func (l *Lobby) Draw(screen *ebiten.Image) {
	screen.Fill(screenproperties.BackgroundColor)

	opts := ebiten.DrawImageOptions{}
	// opts.GeoM.Translate(l.Player.X, l.Player.Y)

	// screen.DrawImage(
	// 	l.Player.Img.SubImage(
	// 		l.Player.SpriteSheet.Rect(l.Player.ActiveAnimation.Frame()),
	// 	).(*ebiten.Image),
	// 	&opts,
	// )
	//
	// opts.GeoM.Reset()

	for _, player := range l.Players {
		opts.GeoM.Translate(player.X, player.Y)

		player.ActiveAnimation = player.GetActiveAnimation()
		screen.DrawImage(
			player.Img.SubImage(
				player.SpriteSheet.Rect(player.ActiveAnimation.Frame()),
			).(*ebiten.Image),
			&opts,
		)

		opts.GeoM.Reset()
	}

	for _, player := range l.Players {
		textOpts := text.DrawOptions{
			LayoutOptions: player.NameTag.LayoutOptions,
		}
		textOpts.GeoM.Translate(player.NameTag.X, player.NameTag.Y)
		text.Draw(screen, player.Data.Name, player.NameTag.Face, &textOpts)

		textOpts.GeoM.Reset()
	}

	waitText := "Waiting for host..."
	if l.Player.Data.Host == true {
		waitText = "Press ENTER to start!"
	}
	textOpts := text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: 2,
		},
	}
	textOpts.GeoM.Translate(screenproperties.BottomRight())
	text.Draw(screen, waitText, &text.GoTextFace{Source: shared.FontSrc, Size: 8}, &textOpts)
}
