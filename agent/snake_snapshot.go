package agent

import (
	"github.com/BattlesnakeOfficial/rules"
	// "github.com/samber/mo"
	"github.com/samber/lo"
	// "log"
)

type SnakeSnapshot interface {
	ID() string
	Name() string
	Alive() bool
	Health() int
	Body() []rules.Point
	Head() rules.Point
	Length() int
	LastShout() string
	ConsideredMoves() []rules.SnakeMove
}

// SnakeSnapshot interface implementation
type snakeStatsImpl struct {
	name            string
	lastShout       string
	turnLastShouted int
}

type snakeSnapshotImpl struct {
	stats        *snakeStatsImpl
	snake        *rules.Snake
	gameSnapshot *gameSnapshotImpl
}

func (s *snakeSnapshotImpl) ID() string {
	return s.snake.ID
}

func (s *snakeSnapshotImpl) Name() string {
	return s.stats.name
}

func (s *snakeSnapshotImpl) Alive() bool {
	return s.snake.EliminatedCause == rules.NotEliminated
}

func (s *snakeSnapshotImpl) Health() int {
	if !s.Alive() {
		return 0
	} else {
		return s.snake.Health
	}
}

func (s *snakeSnapshotImpl) Body() []rules.Point {
	return s.snake.Body
}

func (s *snakeSnapshotImpl) Head() rules.Point {
	return s.snake.Body[0]
}

func (s *snakeSnapshotImpl) Length() int {
	return len(s.snake.Body)
}

func (s *snakeSnapshotImpl) LastShout() string {
	return s.stats.lastShout
}

func (s *snakeSnapshotImpl) ConsideredMoves() []rules.SnakeMove {
	possibleMoveStrs := []string{"up", "down", "left", "right"}

	isPassable := func(move string) bool {
		board := s.gameSnapshot.Board()
		target := s.getTargetPoint(move)
		if target.X < 0 || target.X >= board.Width || target.Y < 0 || target.Y >= board.Height {
			return false
		}
		return board.Cells[target.Y][target.X].IsPassable()
	}

	consideredMoves := lo.FilterMap(possibleMoveStrs, func(move string, _ int) (rules.SnakeMove, bool) {
		return rules.SnakeMove{ID: s.ID(), Move: move}, isPassable(move)
	})

	if len(consideredMoves) == 0 {
		return []rules.SnakeMove{{ID: s.ID(), Move: "up"}}
	}

	return consideredMoves
}

func (s *snakeSnapshotImpl) getTargetPoint(move string) rules.Point {
	head := s.Head()
	switch move {
	case "up":
		return rules.Point{X: head.X, Y: head.Y + 1}
	case "down":
		return rules.Point{X: head.X, Y: head.Y - 1}
	case "left":
		return rules.Point{X: head.X - 1, Y: head.Y}
	case "right":
		return rules.Point{X: head.X + 1, Y: head.Y}
	default:
		return head
	}
}
