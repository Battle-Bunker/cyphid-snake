
package boardutils

import (
	"github.com/Battle-Bunker/cyphid-snake/agent"
	"github.com/BattlesnakeOfficial/rules"
)

type CellWithDist struct {
	Cell agent.Cell
	Dist int
}

func FindNearest(board *agent.Board, start rules.Point, predicate func(agent.Cell) bool) (agent.Cell, int) {
	visited := make(map[rules.Point]bool)
	queue := []CellWithDist{{board.Cells[start.Y][start.X], 0}}
	visited[start] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if predicate(current.Cell) {
			return current.Cell, current.Dist
		}

		for _, neighbor := range current.Cell.PassableNeighbours(board) {
			pos := neighbor.Coordinates()
			if !visited[pos] {
				visited[pos] = true
				queue = append(queue, CellWithDist{neighbor, current.Dist + 1})
			}
		}
	}

	return nil, -1
}

// FloodFill returns the count of reachable cells and whether a target position is reachable
func FloodFill(board *agent.Board, start rules.Point, target *rules.Point) (int, bool) {
	return 0, false
}
