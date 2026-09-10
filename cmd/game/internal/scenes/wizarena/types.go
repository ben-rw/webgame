package wizarena

import (
	"image"
	"log"

	"github.com/ben-rw/webgame/cmd/game/internal/shared"
	"github.com/ben-rw/webgame/cmd/game/internal/shared/sound"
	"github.com/ben-rw/webgame/cmd/game/internal/ws"
	"github.com/ben-rw/webgame/internal/protocol"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

const (
	tilemapPath = "assets/maps/ninja_dungeon.json"
	songPath    = "assets/audio/music/void-construct-loop.ogg"
	introLen    = 3
)

type WizardPlayer struct {
	*shared.Player
	Combat *WizardCombat
}

type WizardCombat struct {
	*shared.BasicCombat
	Dead bool
}

type WizArena struct {
	shared.Roster
	Conn            *ws.Connection
	Sprites         []*shared.Sprite
	wizard          *WizardPlayer
	wizards         map[string]*WizardPlayer
	enemies         []*shared.Enemy
	tilemapJSON     *shared.TilemapJSON
	tileCache       map[int]*shared.Tile
	camera          *shared.Camera
	colliders       []image.Rectangle
	audioPlayer     *audio.Player
	projectiles     map[shared.ProjectileType]*shared.Projectile
	projectileCache map[shared.ProjectileType]*ebiten.Image
}

func NewWizArena(c *ws.Connection) *WizArena {
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

	audioPlayer, err := sound.NewAudioPlayer(songPath, true, introLen)
	if err != nil {
		log.Printf("couldn't create audio player: %v", err)
	}

	fireballImg, err := shared.LoadProjectile(shared.Fireball)
	if err != nil {
		log.Printf("couldn't load projectile image: %v", err)
	}

	w := &WizArena{
		Roster: shared.Roster{
			Players: make(map[string]*shared.Player, 8),
			Player:  shared.NewPlayer(&protocol.PlayerData{}, 0),
		},
		Conn:            c,
		Sprites:         []*shared.Sprite{},
		wizards:         make(map[string]*WizardPlayer, 0),
		enemies:         make([]*shared.Enemy, 0),
		projectiles:     make(map[shared.ProjectileType]*shared.Projectile, 0),
		projectileCache: make(map[shared.ProjectileType]*ebiten.Image, 0),
		tilemapJSON:     tilemap,
		tileCache:       tileCache,
		camera:          shared.NewCamera(0.0, 0.0),
		colliders: []image.Rectangle{
			image.Rect(100, 100, 116, 116),
		},
		audioPlayer: audioPlayer,
	}

	w.enemies = append(w.enemies, shared.NewEnemy(shared.Skeleton, true, 400, 300))
	w.enemies = append(w.enemies, shared.NewEnemy(shared.Skeleton, true, 200, 200))
	w.enemies = append(w.enemies, shared.NewEnemy(shared.Skeleton, true, 400, 400))

	w.projectileCache[shared.Fireball] = fireballImg

	// for _, layer := range w.tilemapJSON.Layers {
	// 	for i, id := range layer.Data {
	// 		if id == 0 {
	// 			continue
	// 		}
	// 		if id == ("any kind of hole") {
	// 			"make collider + add to colliders"
	// 		}
	// 	}
	// }

	return w
}

func NewWizard(player *shared.Player) *WizardPlayer {
	return &WizardPlayer{
		player,
		&WizardCombat{
			shared.NewBasicCombat(
				shared.DefaultPlayerHealth,
				shared.DefaultPlayerAttackPower,
				shared.DefaultPlayerMoveSpeed,
				shared.DefaultProjectileSpeed,
				shared.DefaultProjectileSize,
				shared.DefaultPlayerKnockback,
			),
			false,
		},
	}
}
