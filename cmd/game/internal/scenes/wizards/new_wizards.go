package wizards

import (
	"github.com/ben-rw/webgame/cmd/game/internal/shared"
	"github.com/ben-rw/webgame/cmd/game/internal/ws"
	"github.com/ben-rw/webgame/internal/protocol"
	"image"
	"log"
)

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
	camera      *shared.Camera
	colliders   []image.Rectangle
}

func NewWizards(c *ws.Connection) *Wizards {
	log.Println("scene changed to Wizards")
	shared.ScreenHeight = shared.ScreenHeight * 2
	shared.ScreenWidth = shared.ScreenWidth * 2

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
		camera:      shared.NewCamera(0.0, 0.0),
		colliders: []image.Rectangle{
			image.Rect(100, 100, 116, 116),
		},
	}
}
