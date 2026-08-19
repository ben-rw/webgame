package screen

const (
	ScreenWidth  float64 = 320
	ScreenHeight float64 = 180
)

func Center() (float64, float64) {
	return ScreenWidth * 0.5, ScreenHeight * 0.5
}

func BottomCenter() (float64, float64) {
	return ScreenWidth * 0.5, ScreenHeight * 0.9
}
