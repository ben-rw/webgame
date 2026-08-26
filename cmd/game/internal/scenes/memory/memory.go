package memory

import (
	"github.com/ben-rw/webgame/cmd/game/internal/shared"
	"github.com/ben-rw/webgame/cmd/game/internal/shared/screenproperties"
	"github.com/ben-rw/webgame/cmd/game/internal/ws"
	"github.com/ben-rw/webgame/internal/protocol"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	//	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"log"
)

type Memory struct {
	Conn    *ws.Connection
	Players map[string]*shared.Player
	Player  *shared.Player
	Sprites []*shared.Sprite
}

func NewMemory(c *ws.Connection) *Memory {
	log.Println("scene changed to Memory")
	return &Memory{
		Conn:    c,
		Players: map[string]*shared.Player{},
		Player:  shared.NewPlayer(&protocol.PlayerData{}, 0),
		Sprites: []*shared.Sprite{},
	}
}

func (m *Memory) Update(messages []protocol.Message) error {
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
				if _, ok := m.Players[playerData.Name]; ok {
					m.Players[playerData.Name].Data = playerData
				} else {
					player := shared.NewPlayer(playerData, len(m.Players))
					m.Players[playerData.Name] = player
				}
			}

			m.Player = m.Players[data.PlayerData.Name]

			log.Printf("player: %v", *m.Player.Data)

			playerUpdateData := protocol.PlayerUpdateData{
				PlayerData: m.Player.Data,
			}

			m.Conn.WriteMsg(protocol.PlayerUpdate, playerUpdateData)

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
			log.Printf("playerlist: %v", m.Players)

			if _, ok := m.Players[data.PlayerData.Name]; ok {
				m.Players[data.PlayerData.Name].Data = data.PlayerData
			} else {
				player := shared.NewPlayer(data.PlayerData, len(m.Players))
				m.Players[data.PlayerData.Name] = player
			}

			log.Printf("updated: %v", *m.Players[data.PlayerData.Name].Data)

		default:
		}
	}

	return nil
}

func (m *Memory) Draw(screen *ebiten.Image) {
	screen.Fill(screenproperties.BackgroundColor)

	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(m.Player.X, m.Player.Y)

	screen.DrawImage(
		m.Player.Img.SubImage(
			image.Rect(0, 0, 16, 16),
		).(*ebiten.Image),
		&opts,
	)

	opts.GeoM.Reset()

	for _, sprite := range m.Players {
		opts.GeoM.Translate(sprite.X, sprite.Y)

		screen.DrawImage(
			sprite.Img.SubImage(
				image.Rect(0, 0, 16, 16),
			).(*ebiten.Image),
			&opts,
		)

		opts.GeoM.Reset()
	}

	for _, player := range m.Players {
		textOpts := text.DrawOptions{
			LayoutOptions: player.NameTag.LayoutOptions,
		}
		textOpts.GeoM.Translate(player.NameTag.X, player.NameTag.Y)
		text.Draw(screen, player.Data.Name, player.NameTag.Face, &textOpts)

		textOpts.GeoM.Reset()
	}

	waitText := "Hit the drums in the right order!"
	if m.Player.Data.Host == true {
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
