package room

import "math/rand"
import "log"

func nameCollisionSolver(name string, playerList []*Player) string {
	insults := []string{"smelly", "uggo", "stinky", "bald", "cheap", "hairy", "tiny", "pudgy", "boring", "wack", "weird", "cat owner", "train lover", "big dummy"}
	tooLong := []string{"too much", "long name", "shhh", "chill out", "relax", "try again", "y so long"}
	//I have a cat, she's great. Trains are okay.

	log.Printf("name len: %v", len(name))
	if len(name) > 10 {
		name = tooLong[rand.Intn(len(tooLong))]
	}

	var collision bool = false
	for _, player := range playerList {
		if player.Name == name {
			name = insults[rand.Intn(len(insults))]
			collision = true
			break
		}
	}

	if collision {
		name = nameCollisionSolver(name, playerList)
	}

	return name
}
