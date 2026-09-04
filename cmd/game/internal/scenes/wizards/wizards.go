package wizards

import (
	"github.com/ben-rw/webgame/cmd/game/internal/shared"
	"github.com/ben-rw/webgame/cmd/game/internal/shared/screenproperties"
	"github.com/ben-rw/webgame/cmd/game/internal/ws"
	"github.com/ben-rw/webgame/internal/protocol"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image"

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
	tilemapJSON *shared.TilemapJSON
	tileImgList []*ebiten.Image
}

func NewWizards(c *ws.Connection) *Wizards {
	log.Println("scene changed to Wizards")
	screenproperties.ScreenHeight = screenproperties.ScreenHeight * 2
	screenproperties.ScreenWidth = screenproperties.ScreenWidth * 2

	tilemap, err := shared.NewTilemapJSON("assets/maps/ninja_dungeon.json")
	if err != nil {
		log.Printf("couldn't load tilemap: %v", err)
	}
	tileImgList, err := shared.NewTileImgList(tilemap)
	if err != nil {
		log.Printf("couldn't build tile image list: %v", err)
	}

	return &Wizards{
		Roster: shared.Roster{
			Players: make(map[string]*shared.Player, 8),
			Player:  shared.NewPlayer(&protocol.PlayerData{}, 0),
		},
		Conn:        c,
		Sprites:     []*shared.Sprite{},
		PlayerStats: make(map[*shared.Player]Stats, 8),
		Projectiles: make([]*Projectile, 0),
		tilemapJSON: tilemap,
		tileImgList: tileImgList,
	}
}

func (w Wizards) Update(messages []protocol.Message) error {
	for _, message := range messages {
		switch message.Type {
		case protocol.JoinResponse:
			err := w.HandleJoinResponse(message)
			if err != nil {
				log.Println(err)
				continue
			}

			w.PlayerStats[w.Player] = Stats{
				MoveSpeed:       defaultMoveSpeed,
				ProjectileSpeed: defaultProjectileSpeed,
				ProjectileSize:  defaultProjectileSize,
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
		w.Player.X += w.PlayerStats[w.Player].MoveSpeed
		w.Player.NameTag.X += w.PlayerStats[w.Player].MoveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		w.Player.X -= w.PlayerStats[w.Player].MoveSpeed
		w.Player.NameTag.X -= w.PlayerStats[w.Player].MoveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		w.Player.Y -= w.PlayerStats[w.Player].MoveSpeed
		w.Player.NameTag.Y -= w.PlayerStats[w.Player].MoveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		w.Player.Y += w.PlayerStats[w.Player].MoveSpeed
		w.Player.NameTag.Y += w.PlayerStats[w.Player].MoveSpeed
	}

	for _, player := range w.Players {
		player.ActiveAnimation = player.GetActiveAnimation()
		player.ActiveAnimation.Update()
	}

	return nil
}

func (w Wizards) Draw(screen *ebiten.Image) {
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

			tileImgIndex := shared.GetTileImgIndex(id, w.tilemapJSON)
			tileImg := w.tileImgList[tileImgIndex]

			srcX := (id - w.tilemapJSON.Tilesets[tileImgIndex].Firstgid) % w.tilemapJSON.Tilesets[tileImgIndex].Data.Columns
			srcY := (id - w.tilemapJSON.Tilesets[tileImgIndex].Firstgid) / w.tilemapJSON.Tilesets[tileImgIndex].Data.Columns

			srcX *= shared.TileSize
			srcY *= shared.TileSize

			opts.GeoM.Translate(float64(x), float64(y))

			screen.DrawImage(
				tileImg.SubImage(image.Rect(srcX, srcY, srcX+shared.TileSize, srcY+shared.TileSize)).(*ebiten.Image),
				&opts,
			)

			opts.GeoM.Reset()
		}
	}

	for _, player := range w.Players {
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
	for _, player := range w.Players {
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
