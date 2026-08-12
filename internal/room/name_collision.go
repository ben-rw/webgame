package room

import "math/rand"

func nameCollisionSolver(name string, playerList []*Player) string {
	insults := []string{"smelly", "butt-ugly", "bald", "cheapskate", "hairy", "tiny", "pudgy", "boring", "wack", "weird", "cat owner", "train lover", "big dummy"}
	//I have a cat, she's great. Trains are okay.

	var collision bool = false
	for _, player := range playerList {
		if player.Name == name {
			name = name + " the " + insults[rand.Intn(len(insults))]
			collision = true
			break
		}
	}

	if collision {
		name = nameCollisionSolver(name, playerList)
	}

	return name
}
