package wizards

import (
	"github.com/ben-rw/webgame/cmd/game/internal/shared"
	"github.com/ben-rw/webgame/cmd/game/internal/shared/screenproperties"
	"github.com/ben-rw/webgame/cmd/game/internal/ws"
	"github.com/ben-rw/webgame/internal/protocol"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"math"

	"log"
)

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
	targetX += shared.TileSize / 2
	targetY += shared.TileSize / 2
	c.X = -targetX + screenproperties.ScreenWidth/2.0
	c.Y = -targetY + screenproperties.ScreenHeight/2.0
}

const (
	defaultMoveSpeed       = 2
	defaultProjectileSpeed = 5
	defaultProjectileSize  = 1
	tilemapPath            = "assets/maps/ninja_dungeon.json"
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
	PlayerStats map[string]*Stats
	Projectiles []*Projectile
	tilemapJSON *shared.TilemapJSON
	tileCache   map[int]*shared.Tile
	camera      *Camera
}

func NewWizards(c *ws.Connection) *Wizards {
	log.Println("scene changed to Wizards")
	screenproperties.ScreenHeight = screenproperties.ScreenHeight * 2
	screenproperties.ScreenWidth = screenproperties.ScreenWidth * 2

	tilemap, err := shared.NewTilemapJSON(tilemapPath)
	if err != nil {
		log.Printf("couldn't load tilemap: %v", err)
	}
	tileCache, err := shared.NewTileCache(tilemap)
	if err != nil {
		log.Printf("couldn't build tile cache: %v", err)
	}

	return &Wizards{
		Roster: shared.Roster{
			Players: make(map[string]*shared.Player, 8),
			Player:  shared.NewPlayer(&protocol.PlayerData{}, 0),
		},
		Conn:        c,
		Sprites:     []*shared.Sprite{},
		PlayerStats: make(map[string]*Stats, 8),
		Projectiles: make([]*Projectile, 0),
		tilemapJSON: tilemap,
		tileCache:   tileCache,
		camera:      NewCamera(0.0, 0.0),
	}
}

func (w *Wizards) Update(messages []protocol.Message) error {
	for _, message := range messages {
		switch message.Type {
		case protocol.JoinResponse:
			err := w.HandleJoinResponse(message)
			if err != nil {
				log.Println(err)
				continue
			}

			for _, player := range w.Players {
				if _, ok := w.PlayerStats[player.Data.Name]; !ok {
					w.PlayerStats[player.Data.Name] = &Stats{
						MoveSpeed:       defaultMoveSpeed,
						ProjectileSpeed: defaultProjectileSpeed,
						ProjectileSize:  defaultProjectileSize,
					}
				}
			}

			playerUpdateData := protocol.PlayerUpdateData{
				PlayerData: w.Player.Data,
			}

			w.Conn.WriteMsg(protocol.PlayerUpdate, playerUpdateData)

		case protocol.PlayerUpdate:
			err := w.HandlePlayerUpdate(message)
			if err != nil {
				log.Println(err)
				continue
			}

		default:
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		w.Player.X += w.PlayerStats[w.Player.Data.Name].MoveSpeed
		w.Player.NameTag.X += w.PlayerStats[w.Player.Data.Name].MoveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		w.Player.X -= w.PlayerStats[w.Player.Data.Name].MoveSpeed
		w.Player.NameTag.X -= w.PlayerStats[w.Player.Data.Name].MoveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		w.Player.Y -= w.PlayerStats[w.Player.Data.Name].MoveSpeed
		w.Player.NameTag.Y -= w.PlayerStats[w.Player.Data.Name].MoveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		w.Player.Y += w.PlayerStats[w.Player.Data.Name].MoveSpeed
		w.Player.NameTag.Y += w.PlayerStats[w.Player.Data.Name].MoveSpeed
	}

	for _, player := range w.Players {
		player.ActiveAnimation = player.GetActiveAnimation()
		player.ActiveAnimation.Update()
	}

	w.camera.FollowTarget(w.Player.X, w.Player.Y)

	return nil
}

func (w *Wizards) Draw(screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}

	for _, layer := range w.tilemapJSON.Layers {
		for i, id := range layer.Data {
			if id == 0 {
				continue
			}
			x := i % layer.Width
			y := i / layer.Width

			x *= shared.TileSize
			y *= shared.TileSize

			tile := shared.Tile{}

			if id&int(shared.FlagFlippedHorizontally) != 0 {
				tile.Flips.HorizontalFlip = true
			}
			if id&int(shared.FlagFlippedVertically) != 0 {
				tile.Flips.VerticalFlip = true
			}
			if id&int(shared.FlagFlippedDiagonally) != 0 {
				tile.Flips.DiagonalFlip = true
			}

			id &= ^(int(shared.FlagFlippedHorizontally) |
				int(shared.FlagFlippedVertically) |
				int(shared.FlagFlippedDiagonally) |
				int(shared.FlagRotatedHexagonal120))

			switch {
			case tile.Flips.HorizontalFlip && tile.Flips.VerticalFlip:
				opts.GeoM.Scale(-1, -1)
				x += 16
				y += 16
			case tile.Flips.DiagonalFlip && tile.Flips.HorizontalFlip:
				opts.GeoM.Translate(-shared.TileSize/2, -shared.TileSize/2)
				opts.GeoM.Rotate(math.Pi / 2)
				opts.GeoM.Translate(shared.TileSize/2, shared.TileSize/2)
			case tile.Flips.HorizontalFlip:
				opts.GeoM.Scale(-1, 1)
				x += 16
			case tile.Flips.DiagonalFlip && tile.Flips.VerticalFlip:
				opts.GeoM.Translate(-shared.TileSize/2, -shared.TileSize/2)
				opts.GeoM.Rotate(3 * math.Pi / 2)
				opts.GeoM.Translate(shared.TileSize/2, shared.TileSize/2)
			case tile.Flips.VerticalFlip:
				opts.GeoM.Scale(1, -1)
				y += 16
			default:
			}

			opts.GeoM.Translate(float64(x), float64(y))

			opts.GeoM.Translate(w.camera.X, w.camera.Y)

			screen.DrawImage(
				w.tileCache[id].Img,
				&opts,
			)

			opts.GeoM.Reset()
		}
	}

	for _, player := range w.Players {
		opts.GeoM.Translate(player.X, player.Y)

		opts.GeoM.Translate(w.camera.X, w.camera.Y)

		player.ActiveAnimation = player.GetActiveAnimation()
		screen.DrawImage(
			player.Img.SubImage(
				player.SpriteSheet.Rect(player.ActiveAnimation.Frame()),
			).(*ebiten.Image),
			&opts,
		)

		opts.GeoM.Reset()
	}
	for _, player := range w.Players {
		textOpts := text.DrawOptions{
			LayoutOptions: player.NameTag.LayoutOptions,
		}
		textOpts.GeoM.Translate(player.NameTag.X, player.NameTag.Y)

		textOpts.GeoM.Translate(w.camera.X, w.camera.Y)

		text.Draw(screen, player.Data.Name, player.NameTag.Face, &textOpts)

		textOpts.GeoM.Reset()
	}

	waitText := "Shoot zombies for power-ups! Shoot your friends for glory!"
	textOpts := text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignEnd,
		},
	}
	textOpts.GeoM.Translate(screenproperties.BottomRight())

	text.Draw(screen, waitText, &text.GoTextFace{Source: shared.FontSrc, Size: 8}, &textOpts)
}
