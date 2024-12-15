
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
}

type EmptyCell struct{}

func (e EmptyCell) Kind() CellKind { return CellEmpty }
func (e EmptyCell) IsPassable() bool { return true }

type FoodCell struct{}

func (f FoodCell) Kind() CellKind { return CellFood }
func (f FoodCell) IsPassable() bool { return true }

// We can differentiate snake parts by a named type:
type SnakePartType int
const (
  SnakePartHead SnakePartType = iota
  SnakePartBody
  SnakePartTail
)

type SnakePartCell struct {
  Snake        *SnakeSnapshot
  PartType     SnakePartType
  // If this is a tail part that will disappear next turn, we can indicate it here:
  WillVanishNextTurn bool
}

func (spc SnakePartCell) Kind() CellKind {
  switch spc.PartType {
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

func (spc SnakePartCell) IsPassable() bool {
  return spc.PartType == SnakePartTail && spc.WillVanishNextTurn
}

type Board struct {
    Width, Height int
    Cells         [][]Cell
}
