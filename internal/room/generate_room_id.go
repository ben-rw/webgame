package room

import "math/rand"

func GenerateRoomID() string {
	const codeLength = 4
	var charSet = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	buf := make([]rune, codeLength)
	for i := range buf {
		buf[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(buf)
}
