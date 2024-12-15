package main

import (
  "github.com/Battle-Bunker/cyphid-snake/agent"
  "math"
)

func HeuristicFood(snapshot agent.GameSnapshot) float64 {
  foodScore := 0.0
  food := snapshot.Food()

  for _, snake := range snapshot.YourTeam() {
    if snake.Health() == 100 {
      foodScore += 1.0 // Count as distance of 1 to avoid division by zero
      continue
    }

    // Find closest food
    minDist := math.MaxFloat64
    head := snake.Head()

    for _, foodPos := range food {
      dist := math.Abs(float64(head.X-foodPos.X)) + math.Abs(float64(head.Y-foodPos.Y))
      if dist < minDist {
        minDist = dist
      }
    }

    if minDist != math.MaxFloat64 {
      foodScore += 1.0 / minDist
    }
  }

  return foodScore * 100
}
