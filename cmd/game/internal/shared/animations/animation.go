package animations

type Animation struct {
	FirstFrame   int
	LastFrame    int
	Step         int
	SpeedInTps   float32
	frameCounter float32
	currentFrame int
	Over         bool
}

func NewAnimation(first, last, step int, speed float32) *Animation {
	return &Animation{
		first,
		last,
		step,
		speed,
		speed,
		first,
		false,
	}
}

func (a *Animation) Frame() int {
	return a.currentFrame
}

func (a *Animation) Update() {
	a.Over = false
	a.frameCounter -= 1.0
	if a.frameCounter < 0.0 {
		a.frameCounter = a.SpeedInTps
		a.currentFrame += a.Step
		if a.currentFrame > a.LastFrame {
			a.currentFrame = a.FirstFrame
			a.Over = true
		}
	}
}
