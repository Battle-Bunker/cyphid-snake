
package lib

import (
	"github.com/Battle-Bunker/cyphid-snake/agent"
	"github.com/BattlesnakeOfficial/rules"
)

type cellWithDist struct {
	cell agent.Cell
	dist int
}

func FindNearest(board *agent.Board, start rules.Point, predicate func(agent.Cell) bool) (agent.Cell, int) {
	visited := make(map[rules.Point]bool)
	queue := []cellWithDist{{board.Cells[start.Y][start.X], 0}}
	visited[start] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if predicate(current.cell) {
			return current.cell, current.dist
		}

		for _, neighbor := range current.cell.PassableNeighbours(board) {
			pos := neighbor.Coordinates()
			if !visited[pos] {
				visited[pos] = true
				queue = append(queue, cellWithDist{neighbor, current.dist + 1})
			}
		}
	}

	return nil, -1
}
