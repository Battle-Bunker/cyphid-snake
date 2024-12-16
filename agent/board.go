package agent

type CellKind int

const (
	CellEmpty CellKind = iota
	CellFood
	CellSnakeHead
	CellSnakeBody
	CellSnakeTail
)

type Cell interface {
	Kind() CellKind
	IsPassable() bool
	Coordinates() rules.Point
}

type EmptyCell struct {
	coordinates rules.Point
}

func (e EmptyCell) Kind() CellKind {
	return CellEmpty
}
func (e EmptyCell) IsPassable() bool {
	return true
}
func (e EmptyCell) Coordinates() rules.Point {
	return e.coordinates
}

type FoodCell struct {
	coordinates rules.Point
}

func (f FoodCell) Kind() CellKind {
	return CellFood
}
func (f FoodCell) IsPassable() bool {
	return true
}
func (f FoodCell) Coordinates() rules.Point {
	return f.coordinates
}

// We can differentiate snake parts by a named type:
type SnakePartType int

const (
	SnakePartHead SnakePartType = iota
	SnakePartBody
	SnakePartTail
)

type SnakePartCell struct {
	coordinates        rules.Point
	SnakeID           string
	PartType          SnakePartType
	WillVanishNextTurn bool
}

func (s SnakePartCell) Coordinates() rules.Point {
	return s.coordinates
}
func (s SnakePartCell) Kind() CellKind {
	switch s.PartType {
	case SnakePartHead:
		return CellSnakeHead
	case SnakePartBody:
		return CellSnakeBody
	case SnakePartTail:
		return CellSnakeTail
	default:
		return CellEmpty // should never happen
	}
}

func (s SnakePartCell) IsPassable() bool {
	return s.PartType == SnakePartTail && s.WillVanishNextTurn
}

type Board struct {
	Width, Height int
	Cells         [][]Cell
}

func NewBoard(g GameSnapshot) *Board {
	board := &Board{
		Width:  g.Width(),
		Height: g.Height(),
		Cells:  make([][]Cell, g.Height()),
	}

	// Initialize cells
	for y := 0; y < g.Height(); y++ {
		board.Cells[y] = make([]Cell, g.Width())
		for x := 0; x < g.Width(); x++ {
			board.Cells[y][x] = EmptyCell{coordinates: rules.Point{X: x, Y: y}}
		}
	}

	// Place food
	for _, food := range g.Food() {
		if food.Y < g.Height() && food.X < g.Width() {
			board.Cells[food.Y][food.X] = FoodCell{coordinates: food}
		}
	}

	// Place snakes
	for _, snake := range g.Snakes() {
		body := snake.Body()
		if len(body) == 0 {
			continue
		}

		// Head
		if body[0].Y < g.Height() && body[0].X < g.Width() {
			board.Cells[body[0].Y][body[0].X] = SnakePartCell{coordinates: body[0], SnakeID: snake.ID(), PartType: SnakePartHead}
		}

		// Body
		for i := 1; i < len(body)-1; i++ {
			if body[i].Y < g.Height() && body[i].X < g.Width() {
				board.Cells[body[i].Y][body[i].X] = SnakePartCell{coordinates: body[i], SnakeID: snake.ID(), PartType: SnakePartBody}
			}
		}

		// Tail
		if len(body) > 1 {
			tail := body[len(body)-1]
			if tail.Y < g.Height() && tail.X < g.Width() {
				board.Cells[tail.Y][tail.X] = SnakePartCell{coordinates: tail, SnakeID: snake.ID(), PartType: SnakePartTail, WillVanishNextTurn: snake.Health() < 100 && len(snake.Body()) >= 3}
			}
		}
	}

	return board
}