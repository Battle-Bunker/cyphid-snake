# AI CONTEXT

## Project Overview
This codebase implements a Battlesnake agent server for a custom variant of a Battle Snake game in which teams of snakes compete to maximize the total number of moves survived by all the snakes in their team by the end of the game (which ends, as usual when there is either 1 or 0 snake remaining on the board).  Moves survived by one's team are called 'points'. If our team has 3 snakes alive at the end of a turn then we score 3 points for that turn. The last snake standing (if any) gains points equal to its health at the end of the game.
The agent is implemented with the library called CyphidSnake, which provides the interfaces defined below. This library will use a variant of Monte Carlo Tree Search that uses a portfolio of heuristics to evaluate each board position and it searches board states with greater probability the higher it scores according to that heuristic portfolio. It will return a move with greater probability the higher it scores according to the expected score calculated by the tree search algorithm for that move. 

## Core Concepts
- **Scoring System**: Teams earn points for each turn where their snakes survive
  - Each surviving snake contributes 1 point per turn
  - The last surviving snake (if any) gains bonus points equal to its remaining health
- **Game End Conditions**: Game ends when 0-1 snakes remain on the board

## Implementation Guidelines

### Heuristic Function Requirements
1. **Purpose**: Evaluate board states to estimate expected future team survival time
2. **Characteristics**:
   - Computationally inexpensive
   - Score variance should correlate with actual quality differences between states (in terms of expected moves survived by team)
   - Focus on a single dimension of value
   - Return values scaled to approximate expected future moves survived

### Best Practices
- Keep implementations simple and elegant
- Separate distinct value dimensions into different heuristic functions
- Consider correlation between heuristic scores and actual game outcomes (turns survived)
- Optimize for efficient differentiation between board states

### Example Heuristic Function
Here's an example of a simple heuristic that returns the total team health, which is a good heuristic for our goals because it is cheap to compute and correlates with expected future moves survived by our team because it's the number of moves that would be survived if the snakes in our team neither collected any food nor died from collision. It is also, importantly, efficiently diagnostic of *differences* in expected points yet to be earned between board states because, if one of our snakes dies it deducts points that it would have earned if surviving to starvation and, if a snake eats, it adds points corresponding to additional turns that could be survived until starvation by that snake.
```go
func HeuristicHealth(snapshot agent.GameSnapshot) float64 {
    totalHealth := 0.0
    for _, allySnake := range snapshot.YourTeam() {
        totalHealth += float64(allySnake.Health())
    }
    return totalHealth
}
```

### Guiding the User
If you are asked to write a heuristic function that seems to you to recommend a scoring rule that would significantly deviate from the goal of approximating moves that will be survived by our team in the future, please politely suggest to the user how they might improve their scoring rule in order to adhere to this goal better. Strive to empathize with what valuable behavior of a Battle Snake they are trying to reward with their heuristic, given their description, and suggest an adjusted plan that would cause the heuristic to return a higher score the better the snake is doing at meeting the valuable condition or behavior, scaled such that *differences* in score between board states could be expected to correlate with differences in expected points yet to be earned from that board state.

A good heuristic function should be simple and elegant and address one dimension of value at a time. If the user suggests incorporating many dimensions of value into a single heuristic function you should suggest that they break it up into multiple heuristic functions.

### Recommending Code Location
When asked to write heuristic functions, always recommend a filename for the code to go in using the schema heuristic_<name>.go for a heuristic function called Heuristic<Name>.


## Available Interfaces

### Core Types
```go
type HeuristicFunc func(GameSnapshot) float64

type GameSnapshot interface {
    GameID() string
    Rules() rules.Ruleset
    Turn() int
    Height() int
    Width() int
    Food() []rules.Point
    Hazards() []rules.Point
    You() SnakeSnapshot
    Snakes() []SnakeSnapshot
    Teammates() []SnakeSnapshot
    YourTeam() []SnakeSnapshot
    Opponents() []SnakeSnapshot
    AllSnakes() []SnakeSnapshot
    DeadSnakes() []SnakeSnapshot
    ApplyMoves(moves []rules.SnakeMove) (GameSnapshot, error)
}

type SnakeSnapshot interface {
    ID() string
    Name() string
    Alive() bool
    Health() int
    Body() []rules.Point
    Head() rules.Point
    Length() int
    LastShout() string
    ForwardMoves() []rules.SnakeMove
}
```
h
## Imports Guide

When implementing your heuristic function, import only the packages you directly use. The `agent` package provides the core interfaces (`GameSnapshot`, `_SnakeSnapshot_`), and the `rules` package provides supporting types like `Point`. For example:

```go
package main

import (
    _ "github.com/Battle-Bunker/cyphid-snake/agent"  // Import if using GameSnapshot or SnakeSnapshot
    _ "github.com/BattlesnakeOfficial/rules"        // Import if using Point or other rules types
)
```

Only include imports that your code actually references. The compiler will help ensure you have the correct imports.__