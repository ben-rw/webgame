package sound

import (
	"bytes"
	"io"

	"github.com/ben-rw/webgame/cmd/game/internal/shared"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

const SampleRate = 48000

var audioContext = audio.NewContext(SampleRate)

func NewAudioPlayer(filepath string) (*audio.Player, error) {
	data, err := shared.AssetsFS.ReadFile(filepath)
	stream, err := vorbis.DecodeWithSampleRate(SampleRate, bytes.NewReader(data))
	if err != nil {
		return &audio.Player{}, err
	}
	var s io.ReadSeeker
	s = audio.NewInfiniteLoop(stream, stream.Length())
	audioPlayer, err := audioContext.NewPlayer(s)
	if err != nil {
		return &audio.Player{}, err
	}

	return audioPlayer, nil
}
