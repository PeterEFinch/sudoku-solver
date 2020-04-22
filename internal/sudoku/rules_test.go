package sudoku

import (
	"context"
	"fmt"
	"testing"
)

func TestChessKnightRule2(t *testing.T) {
	grid := Grid{
		values: [9][9]int{
			{3, 0, 0, 0, 1, 8, 2, 0, 0},
			{2, 0, 8, 0, 6, 4, 0, 0, 0},
			{0, 0, 1, 0, 0, 0, 4, 0, 8},
			{4, 0, 0, 0, 0, 0, 5, 8, 0},
			{0, 0, 0, 4, 0, 6, 0, 0, 0},
			{0, 3, 2, 0, 0, 0, 0, 0, 0},
			{6, 0, 9, 0, 0, 0, 3, 0, 0},
			{0, 0, 0, 3, 4, 0, 6, 0, 9},
			{0, 0, 3, 6, 9, 0, 0, 0, 2},
		},
	}

	solver := NewSolver(
		RowRule(),
		ColumnRule(),
		SquareRule(),
		ChessKnightRule(),
	)

	solutions := solver.Solve(context.Background(), grid)
	fmt.Println(len(solutions))
}
