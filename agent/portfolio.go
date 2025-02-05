package agent

import (
	"fmt"
)

// HeuristicPortfolio represents a collection of weighted heuristics.
type HeuristicPortfolio []WeightedHeuristic

import (
	"sync/atomic"
	"time"
)

// WeightedHeuristic is an interface defining the methods to access properties of a weighted heuristic.
type WeightedHeuristic interface {
	Name() string
	F() HeuristicFunc
	Weight() float64
	NameAndWeight() string
	GetAndResetStats() (uint64, uint64) // Returns (microseconds, evaluations)
}

// HeuristicFunc is a type that represents a heuristic function.
// It takes a GameSnapshot as input and returns a float64 score.
type HeuristicFunc func(GameSnapshot) float64

func NewPortfolio(heuristics ...WeightedHeuristic) HeuristicPortfolio {
	return HeuristicPortfolio(heuristics)
}

func NewHeuristic(weight float64, name string, f HeuristicFunc) WeightedHeuristic {
	return &weightedHeuristicImpl{
		name:        name,
		f:           f,
		weight:      weight,
		microsecs:   0,
		evaluations: 0,
	}
}

// weightedHeuristicImpl represents a heuristic with an associated weight and name.
type weightedHeuristicImpl struct {
	name        string
	f           HeuristicFunc
	weight      float64
	microsecs   uint64
	evaluations uint64
}

func (w *weightedHeuristicImpl) Name() string {
	return w.name
}

func (w *weightedHeuristicImpl) F() HeuristicFunc {
	return func(snapshot GameSnapshot) float64 {
		start := time.Now()
		result := w.f(snapshot)
		elapsed := time.Since(start)
		atomic.AddUint64(&w.microsecs, uint64(elapsed.Microseconds()))
		atomic.AddUint64(&w.evaluations, 1)
		return result
	}
}

func (w *weightedHeuristicImpl) Weight() float64 {
	return w.weight
}

func (w *weightedHeuristicImpl) NameAndWeight() string {
	return fmt.Sprintf("%s, w=%.2f", w.name, w.weight)
}

func (w *weightedHeuristicImpl) GetAndResetStats() (uint64, uint64) {
	micros := atomic.SwapUint64(&w.microsecs, 0)
	evals := atomic.SwapUint64(&w.evaluations, 0)
	return micros, evals
}
