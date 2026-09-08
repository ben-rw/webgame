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
	EnemyAttackCooldown      = 60
)

type Combat interface {
	Health() int
	AttackPower() int
	Damage(amount int)
	Attacking() bool
	Attack() bool
	Update()
}

type BasicCombat struct {
	health          int
	attackPower     int
	moveSpeed       float64
	projectileSpeed float64
	projectileSize  float64
	attacking       bool
}

func (b *BasicCombat) Update() {
}

func (b *BasicCombat) Attack() bool {
	b.attacking = true
	return true
}

func (b *BasicCombat) Attacking() bool {
	return b.attacking
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
		false,
	}
}

type EnemyCombat struct {
	*BasicCombat
	attackCooldown  int
	timeSinceAttack int
}

func (e *EnemyCombat) Attack() bool {
	if e.timeSinceAttack >= e.attackCooldown {
		e.attacking = true
		e.timeSinceAttack = 0
		return true
	}
	return false
}

func (e *EnemyCombat) Update() {
	e.timeSinceAttack += 1
}

func NewEnemyCombat(health, attackPower, attackCooldown int, moveSpeed, projectileSpeed, projectileSize float64) *EnemyCombat {
	return &EnemyCombat{
		NewBasicCombat(
			health,
			attackPower,
			moveSpeed,
			projectileSpeed,
			projectileSize,
		),
		attackCooldown,
		0,
	}
}
