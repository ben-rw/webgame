package shared

import (
	"bytes"
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var FontSrc = getFontSrc()

func getFontSrc() *text.GoTextFaceSource {
	fontBytes, err := AssetsFS.ReadFile("assets/fonts/m3x6.ttf")
	if err != nil {
		log.Printf("error reading font file: %v", err)
	}

	fontSrc, err := text.NewGoTextFaceSource(bytes.NewReader(fontBytes))
	if err != nil {
		log.Printf("error created font face src: %v", err)
	}

	return fontSrc
}
