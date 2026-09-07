package screenproperties

import "image/color"

var BackgroundColor = color.RGBA{0, 105, 150, 1}

const (
	defaultScreenWidth  = 320
	defaultScreenHeight = 180
)

var (
	ScreenWidth  float64 = defaultScreenWidth
	ScreenHeight float64 = defaultScreenHeight
)

func Center() (float64, float64) {
	return ScreenWidth * 0.5, ScreenHeight * 0.5
}

func CenterRight() (float64, float64) {
	return ScreenWidth * 0.9, ScreenHeight * 0.5
}

func BottomCenter() (float64, float64) {
	return ScreenWidth * 0.5, ScreenHeight * 0.9
}

func BottomRight() (float64, float64) {
	return ScreenWidth * 0.99, ScreenHeight * 0.9
}

func TopRight() (float64, float64) {
	return ScreenWidth * 0.99, ScreenHeight * 0.1
}
