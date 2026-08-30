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
	Conn    *ws.Connection
	Players map[string]*shared.Player
	Player  *shared.Player
	Sprites []*shared.Sprite
}

func NewLobby(c *ws.Connection) *Lobby {
	return &Lobby{
		Conn:    c,
		Players: map[string]*shared.Player{},
		Player:  shared.NewPlayer(&protocol.PlayerData{}, 0),
		Sprites: []*shared.Sprite{},
	}
}

func (l *Lobby) Update(messages []protocol.Message) error {
	for _, message := range messages {
		log.Printf("msg type: %v", string(message.Type))
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
					player := shared.NewPlayer(playerData, len(l.Players))
					l.Players[playerData.Name] = player
				}
			}

			l.Player = l.Players[data.PlayerData.Name]

			log.Printf("player: %v", *l.Player.Data)

			playerUpdateData := protocol.PlayerUpdateData{
				PlayerData: l.Player.Data,
			}

			l.Conn.WriteMsg(protocol.PlayerUpdate, playerUpdateData)

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
			log.Printf("playerlist: %v", l.Players)

			if _, ok := l.Players[data.PlayerData.Name]; ok {
				l.Players[data.PlayerData.Name].Data = data.PlayerData
			} else {
				player := shared.NewPlayer(data.PlayerData, len(l.Players))
				l.Players[data.PlayerData.Name] = player
			}

			log.Printf("updated: %v", *l.Players[data.PlayerData.Name].Data)

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
