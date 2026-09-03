package wizards

import (
	"github.com/ben-rw/webgame/cmd/game/internal/shared"
	"github.com/ben-rw/webgame/cmd/game/internal/shared/screenproperties"
	"github.com/ben-rw/webgame/cmd/game/internal/ws"
	"github.com/ben-rw/webgame/internal/protocol"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	//	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
)

const (
	defaultMoveSpeed       = 2
	defaultProjectileSpeed = 5
	defaultProjectileSize  = 1
)

type Stats struct {
	MoveSpeed       float64
	ProjectileSpeed float64
	ProjectileSize  float64
}

type Projectile struct {
	*shared.Sprite
	Speed float64
	Size  float64
}

type Wizards struct {
	shared.Roster
	Conn        *ws.Connection
	Sprites     []*shared.Sprite
	PlayerStats map[*shared.Player]Stats
	Projectiles []*Projectile
}

func NewWizards(c *ws.Connection) *Wizards {
	log.Println("scene changed to Wizards")
	screenproperties.ScreenHeight = screenproperties.ScreenHeight * 2
	screenproperties.ScreenWidth = screenproperties.ScreenWidth * 2
	return &Wizards{
		Roster: shared.Roster{
			Players: make(map[string]*shared.Player, 8),
			Player:  shared.NewPlayer(&protocol.PlayerData{}, 0),
		},
		Conn:        c,
		Sprites:     []*shared.Sprite{},
		PlayerStats: make(map[*shared.Player]Stats, 8),
	}
}

func (m *Wizards) Update(messages []protocol.Message) error {
	for _, message := range messages {
		switch message.Type {
		case protocol.JoinResponse:
			err := m.HandleJoinResponse(message)
			if err != nil {
				log.Println(err)
				continue
			}

			m.PlayerStats[m.Player] = Stats{
				MoveSpeed:       defaultMoveSpeed,
				ProjectileSpeed: defaultProjectileSpeed,
				ProjectileSize:  defaultProjectileSize,
			}

			playerUpdateData := protocol.PlayerUpdateData{
				PlayerData: m.Player.Data,
			}

			m.Conn.WriteMsg(protocol.PlayerUpdate, playerUpdateData)

		case protocol.PlayerUpdate:
			err := m.HandlePlayerUpdate(message)
			if err != nil {
				log.Println(err)
				continue
			}

		default:
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		m.Player.X += m.PlayerStats[m.Player].MoveSpeed
		m.Player.NameTag.X += m.PlayerStats[m.Player].MoveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		m.Player.X -= m.PlayerStats[m.Player].MoveSpeed
		m.Player.NameTag.X -= m.PlayerStats[m.Player].MoveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		m.Player.Y -= m.PlayerStats[m.Player].MoveSpeed
		m.Player.NameTag.Y -= m.PlayerStats[m.Player].MoveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		m.Player.Y += m.PlayerStats[m.Player].MoveSpeed
		m.Player.NameTag.Y += m.PlayerStats[m.Player].MoveSpeed
	}

	for _, player := range m.Players {
		player.ActiveAnimation = player.GetActiveAnimation()
		player.ActiveAnimation.Update()
	}

	return nil
}

func (m *Wizards) Draw(screen *ebiten.Image) {
	screen.Fill(screenproperties.BackgroundColor)

	opts := ebiten.DrawImageOptions{}

	for _, player := range m.Players {
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
	for _, player := range m.Players {
		textOpts := text.DrawOptions{
			LayoutOptions: player.NameTag.LayoutOptions,
		}
		textOpts.GeoM.Translate(player.NameTag.X, player.NameTag.Y)
		text.Draw(screen, player.Data.Name, player.NameTag.Face, &textOpts)

		textOpts.GeoM.Reset()
	}

	waitText := "Shoot zombies for power-ups! Shoot your friends for glory!"
	textOpts := text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: 2,
		},
	}
	textOpts.GeoM.Translate(screenproperties.BottomRight())
	text.Draw(screen, waitText, &text.GoTextFace{Source: shared.FontSrc, Size: 8}, &textOpts)
}
