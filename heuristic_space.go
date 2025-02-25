
package main

import (
	"github.com/Battle-Bunker/cyphid-snake/agent"
	"github.com/Battle-Bunker/cyphid-snake/boardutils"
	// "log"
)

func HeuristicSpace(snapshot agent.GameSnapshot) float64 {
	snake := snapshot.You()
	board := snapshot.Board()
	head := snake.Head()
	
	// Check if we can reach our tail
	tail := snake.Body()[len(snake.Body())-1]
	spaces, tailReachable := boardutils.FloodFill(board, head, &tail)
	
	// log.Printf("Spaces available: %d, Tail reachable: %t", spaces, tailReachable)
	
	if tailReachable {
		return 100.0
	}
	
	// If available space is greater than our length, we have room to maneuver
	if spaces >= len(snake.Body()) {
		return 100.0
	}
	
	return float64(spaces)
}
