package shared

const (
	DefaultPlayerHealth      = 3.0
	DefaultPlayerAttackPower = 1.0
	DefaultPlayerMoveSpeed   = 2.0
	DefaultProjectileSpeed   = 5.0
	DefaultProjectileSize    = 1.0
	EnemyMoveSpeed           = 0.5
	EnemyHealth              = 1.0
	EnemyAttackPower         = 1.0
)

type Combat interface {
	Health() int
	AttackPower() int
	Damage(amount int)
}

type BasicCombat struct {
	health          int
	attackPower     int
	moveSpeed       float64
	projectileSpeed float64
	projectileSize  float64
}

func (b *BasicCombat) AttackPower() int {
	return b.attackPower
}

func (b *BasicCombat) Health() int {
	return b.health
}

func (b *BasicCombat) Damage(amount int) {
	b.health -= amount
}

func (b *BasicCombat) MoveSpeed() float64 {
	return b.moveSpeed
}

func (b *BasicCombat) ProjectileSpeed() float64 {
	return b.projectileSpeed
}

func (b *BasicCombat) ProjectileSize() float64 {
	return b.projectileSize
}

func NewBasicCombat(health, attackPower int, moveSpeed, projectileSpeed, projectileSize float64) *BasicCombat {
	return &BasicCombat{
		health,
		attackPower,
		moveSpeed,
		projectileSpeed,
		projectileSize,
	}
}
