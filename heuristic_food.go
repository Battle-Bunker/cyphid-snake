
package main

import (
	"github.com/Battle-Bunker/cyphid-snake/agent"
	"github.com/Battle-Bunker/cyphid-snake/boardutils"
)

func HeuristicFood(snapshot agent.GameSnapshot) float64 {
	snake := snapshot.You()
	if snake.Health() == 100 {
		return 100.0 // Same as original - full health means no food needed
	}

	board := snapshot.Board()
	head := snake.Head()

	isFoodCell := func(cell agent.Cell) bool {
		return cell.Kind() == agent.CellFood
	}

	_, dist := boardutils.FindNearest(board, head, isFoodCell)
	if dist == -1 {
		return 0.0 // No reachable food
	}

	return 100.0 / float64(dist)
}
