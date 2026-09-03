package memory

import (
	"github.com/ben-rw/webgame/cmd/game/internal/shared"
	"github.com/ben-rw/webgame/cmd/game/internal/shared/screenproperties"
)

// starts bottom center, to bottom left, then top left
var PlayerPositions = map[int]struct{ X, Y float64 }{
	0: {X: screenproperties.ScreenWidth * .5, Y: screenproperties.ScreenHeight * .9},
	1: {X: screenproperties.ScreenWidth * .25, Y: screenproperties.ScreenHeight * .9},
	2: {X: screenproperties.ScreenWidth * .1, Y: screenproperties.ScreenHeight * .9},
	3: {X: screenproperties.ScreenWidth * .1, Y: screenproperties.ScreenHeight * .74},
	4: {X: screenproperties.ScreenWidth * .1, Y: screenproperties.ScreenHeight * .58},
	5: {X: screenproperties.ScreenWidth * .1, Y: screenproperties.ScreenHeight * .42},
	6: {X: screenproperties.ScreenWidth * .1, Y: screenproperties.ScreenHeight * .26},
	7: {X: screenproperties.ScreenWidth * .1, Y: screenproperties.ScreenHeight * .1},
}

// put player at the back of the line after thier turn, remove furthest positions
// based on number of eliminated players, making the line shorter
func getInLine(player *shared.Player, totalPlayers, eliminatedPlayers int) {
	lineIndex := totalPlayers - eliminatedPlayers
	player.ScriptedWalk(PlayerPositions[lineIndex].X, PlayerPositions[lineIndex].Y)
}
