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
