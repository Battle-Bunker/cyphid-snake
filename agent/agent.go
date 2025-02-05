package agent

import (
	"github.com/Battle-Bunker/cyphid-snake/lib"
	"github.com/BattlesnakeOfficial/rules"
	"github.com/BattlesnakeOfficial/rules/client"

	// "github.com/samber/mo"
	"fmt"
	"log"
	// "math"
	"slices"
	"strings"

	"github.com/samber/lo"
	"github.com/samber/lo/parallel"
)

// Update the SnakeAgent structure to include SnakeMetadataResponse
type SnakeAgent struct {
	Portfolio             HeuristicPortfolio
	Metadata              client.SnakeMetadataResponse
	Temperature           float64
	LogPerformanceStats bool
}

// SnakeAgentOption defines a function type for configuring a SnakeAgent
type SnakeAgentOption func(*SnakeAgent)

// WithTemperature sets the temperature for the snake agent
func WithTemperature(temp float64) SnakeAgentOption {
	return func(sa *SnakeAgent) {
		sa.Temperature = temp
	}
}

// WithPerformanceLogging enables or disables performance logging
func WithPerformanceLogging(enabled bool) SnakeAgentOption {
	return func(sa *SnakeAgent) {
		sa.LogPerformanceStats = enabled
	}
}

func NewSnakeAgent(portfolio HeuristicPortfolio, metadata client.SnakeMetadataResponse, opts ...SnakeAgentOption) *SnakeAgent {
	sa := &SnakeAgent{
		Portfolio:             portfolio,
		Metadata:              metadata,
		Temperature:           5.0,  // default temperature
		LogPerformanceStats: true, // default to true
	}

	// Apply all options
	for _, opt := range opts {
		opt(sa)
	}

	return sa
}

// Keep NewSnakeAgentWithTemp for backward compatibility
func NewSnakeAgentWithTemp(portfolio HeuristicPortfolio, temperature float64, metadata client.SnakeMetadataResponse) *SnakeAgent {
	return NewSnakeAgent(portfolio, metadata, WithTemperature(temperature))
}

func (sa *SnakeAgent) ChooseMove(snapshot GameSnapshot) client.MoveResponse {
	you := snapshot.You()
	consideredMoves := you.ConsideredMoves()

	consideredMoveStrs := lo.Map(consideredMoves, func(move rules.SnakeMove, _ int) string { return move.Move })
	slices.Sort(consideredMoveStrs)
	log.Printf("\n\n ### Start Turn %d: Considered Moves = %v", snapshot.Turn(), consideredMoveStrs)

	if sa.LogPerformanceStats {
		defer func() {
			log.Printf("### Performance Stats:")
			for _, h := range sa.Portfolio {
				micros, evals := h.GetAndResetStats()
				if evals > 0 {
					avgMicros := float64(micros) / float64(evals)
					log.Printf("###   %25s: %6d evals, %8.2f µs/eval, %8d µs total",
						h.Name(), evals, avgMicros, micros)
				}
			}
		}()
	}

	// If only one move is available, return it immediately
	if len(consideredMoveStrs) == 1 {
		return client.MoveResponse{
			Move:  consideredMoveStrs[0],
			Shout: "I'm moving " + consideredMoveStrs[0],
		}
	}

	// map: move -> set(state snapshots)
	nextStatesMap := make(map[string][]GameSnapshot)
	for _, move := range consideredMoveStrs {
		nextStatesMap[move] = sa.generateNextStates(snapshot, move)
	}

	// slice of maps, for each heuristic, giving mapping: move -> aggScore
	heuristicScores := parallel.Map(sa.Portfolio, func(heuristic WeightedHeuristic, _ int) map[string]float64 {
		return sa.weightedScoresForHeuristic(heuristic, nextStatesMap, consideredMoveStrs)
	})

	totalHeuristicWeight := lo.SumBy(sa.Portfolio, func(heuristic WeightedHeuristic) float64 {
		return heuristic.Weight()
	})

	// slice of scores aligned with consideredMoveStrs
	normalizedScores := lo.Map(consideredMoveStrs, func(move string, _ int) float64 {
		return lo.SumBy(heuristicScores, func(scores map[string]float64) float64 {
			return scores[move] / totalHeuristicWeight
		})
	})

	probs := lib.SoftmaxWithTemp(normalizedScores, sa.Temperature)

	log.Printf("### %36s: %s", "Aggregate move weights", strings.Join(lo.Map(consideredMoveStrs, func(move string, i int) string {
		return fmt.Sprintf("%s=%6.1f", move, normalizedScores[i])
	}), ", "))
	log.Printf("### %36s: %s", "Aggregate move probabilities", strings.Join(lo.Map(consideredMoveStrs, func(move string, i int) string {
		return fmt.Sprintf("%s=%5.1f%%", move, probs[i]*100)
	}), ", "))

	chosenMove := consideredMoveStrs[lib.SampleFromWeights(probs)]

	return client.MoveResponse{
		Move:  chosenMove,
		Shout: "I'm moving " + chosenMove,
	}
}

func (sa *SnakeAgent) weightedScoresForHeuristic(heuristic WeightedHeuristic, nextStatesMap map[string][]GameSnapshot, consideredMoveStrs []string) map[string]float64 {
	type moveScore struct {
		move  string
		score float64
	}

	scores := parallel.Map(consideredMoveStrs, func(move string, _ int) moveScore {
		states := nextStatesMap[move]
		// Parallelize state evaluation
		stateScores := parallel.Map(states, func(state GameSnapshot, _ int) float64 {
			return heuristic.F()(state)
		})
		mean := lo.Mean(stateScores)
		return moveScore{
			move:  move,
			score: mean,
		}
	})

	moveScores := make(map[string]float64)
	for _, score := range scores {
		moveScores[score.move] = score.score
	}

	log.Printf("MoveScores for %25s: %s", heuristic.NameAndWeight(), strings.Join(lo.Map(consideredMoveStrs, func(move string, _ int) string {
		return fmt.Sprintf("%s=%6.1f", move, moveScores[move])
	}), ", "))

	weightedScores := lo.MapValues(moveScores, func(score float64, _ string) float64 {
		return score * heuristic.Weight()
	})

	return weightedScores
}

func (sa *SnakeAgent) generateNextStates(snapshot GameSnapshot, move string) []GameSnapshot {
	var nextStates []GameSnapshot
	yourID := snapshot.You().ID()

	// Generate all possible move combinations for other snakes
	presetMoves := map[string]rules.SnakeMove{yourID: {ID: yourID, Move: move}}
	moveCombinations := generateConsideredMoveCombinations(snapshot.Snakes(), presetMoves)

	// log.Printf("Trying move %s, combinations: %v", move, getMoveComboList(moveCombinations))

	for _, combination := range moveCombinations {
		// Convert the combination map to a slice
		var moveSlice []rules.SnakeMove
		for _, m := range combination {
			moveSlice = append(moveSlice, m)
		}

		if snapshot == nil {
			log.Fatalf("Snapshot is nil before applying moves")
		}
		nextState, err := snapshot.ApplyMoves(moveSlice)

		if err != nil {
			log.Fatalf("Error applying moves: %v", err)
		} else { // Debug the state after ApplyMoves call
			// log.Printf("Next state after applying move: %+v", nextState)
		}
		if nextState != nil {
			nextStates = append(nextStates, nextState)
		}
	}
	// log.Printf("Generated next states: %+v", nextStates)

	return nextStates
}

func generateConsideredMoveCombinations(snakes []SnakeSnapshot, presetMoves map[string]rules.SnakeMove) []map[string]rules.SnakeMove {
	presetSnakeIDs := lo.Keys(presetMoves)

	nonPresetSnakes := lo.Filter(snakes, func(snake SnakeSnapshot, _ int) bool {
		return !lo.Contains(presetSnakeIDs, snake.ID())
	})

	// If there are no non-preset snakes, return just our preset combination
	if len(nonPresetSnakes) == 0 {
		return []map[string]rules.SnakeMove{presetMoves}
	}

	nonPresetMoves := lo.Map(nonPresetSnakes, func(snake SnakeSnapshot, _ int) []rules.SnakeMove {
		return snake.ConsideredMoves()
	})

	moveCombinations := lib.CartesianProduct(nonPresetMoves...)

	// mix in preset moves to each combo and convert to map from snakeID->move
	mappedCombinations := make([]map[string]rules.SnakeMove, len(moveCombinations))
	for moveSet := range moveCombinations {
		combination := lo.Assign(presetMoves)
		for j, move := range moveSet {
			combination[nonPresetSnakes[j].ID()] = move
		}
		mappedCombinations = append(mappedCombinations, combination)
	}

	return mappedCombinations
}

// for convenient debug printing of move combo collection
func getMoveComboList(moveCombinations []map[string]rules.SnakeMove) [][]string {
	var result [][]string
	for _, combo := range moveCombinations {
		var moves []string
		for _, snakeMove := range combo {
			moves = append(moves, snakeMove.Move)
		}
		result = append(result, moves)
	}
	return result
}
