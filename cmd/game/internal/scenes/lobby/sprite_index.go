package lobby

var imgRootPath string = "assets/images/ninja_adventure/Actor/Character/"
var imgPathEnd string = "/SpriteSheet.png"

var PlayerSpriteIndex map[int]string = map[int]string{
	0: imgRootPath + "Monk" + imgPathEnd,
	1: imgRootPath + "Inspector" + imgPathEnd,
	2: imgRootPath + "MaskFrog" + imgPathEnd,
	3: imgRootPath + "Hunter" + imgPathEnd,
	4: imgRootPath + "Master" + imgPathEnd,
	5: imgRootPath + "Sultan" + imgPathEnd,
	6: imgRootPath + "Samurai" + imgPathEnd,
	7: imgRootPath + "Noble" + imgPathEnd,
}

var StartingPositions map[int]struct{ X, Y float64 } = map[int]struct{ X, Y float64 }{
	0: {X: 20, Y: 20},
	1: {X: 20, Y: 60},
	2: {X: 20, Y: 100},
	3: {X: 20, Y: 140},
	4: {X: 60, Y: 20},
	5: {X: 60, Y: 60},
	6: {X: 60, Y: 100},
	7: {X: 60, Y: 140},
}
