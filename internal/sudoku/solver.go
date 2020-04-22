package sudoku

import (
	"context"
)

// Sta
func StandardSolve(ctx context.Context, grid Grid) []Grid {
	solver := NewSolver(
		RowRule(),
		ColumnRule(),
		SquareRule(),
	)

	return solver.Solve(ctx, grid)
}

type Solver struct {
	recursion bool
	rules     []Rule
}

func NewSolver(rules ...Rule) *Solver {
	return &Solver{
		recursion: true,
		rules:     rules,
	}
}

func (s *Solver) Solve(ctx context.Context, grid Grid) []Grid {
	if s.isInvalid(&grid) {
		return nil
	}

	clone := grid.clone()

	return s.solve(ctx, &clone)
}

func (s *Solver) isInvalid(grid *Grid) bool {
	for _, rule := range s.rules {
		if rule.isInvalid(grid) {
			return true
		}
	}

	return false
}

func (s *Solver) deduction(ctx context.Context, grid *Grid) bool {
	if ctx.Err() != nil {
		return true
	}

	// Initialises possibilities and applies the restrictions
	possibilities := &possibilities{}
	possibilities.initialise()
	for row := 0; row < Size; row++ {
		for column := 0; column < Size; column++ {
			if value := grid.values[row][column]; value > 0 {
				entry := &entry{
					row:    row,
					column: column,
					value:  value,
				}

				possibilities.set(row, column, value)
				for _, rule := range s.rules {
					rule.restrict(entry, possibilities)
				}
			}
		}
	}

	// Repeatedly applies the deductions of the rules
	active := true
	for active && ctx.Err() == nil {
		active = false

		// Applies the deductions of every rule
		for _, rule := range s.rules {
			if ctx.Err() != nil {
				break
			}

			// Checks if any deduction can be made
			entry, used := rule.deduction(grid, possibilities)
			if !used {
				continue
			}

			// Record that the process is still active and
			// checks if there is something to be set
			active = true
			if entry == nil {
				continue
			}

			// Updates grid and possibilities using deduction
			grid.values[entry.row][entry.column] = entry.value
			possibilities.set(entry.row, entry.column, entry.value)
			for _, r := range s.rules {
				r.restrict(entry, possibilities)
			}

			// Checks if grid or  possibilities have become invalid
			if s.isInvalid(grid) && possibilities.isInvalid() {
				return false
			}
		}
	}

	return true
}

func (s *Solver) solve(ctx context.Context, grid *Grid) []Grid {
	// Checks if it is invalid or if is filled
	switch {
	case ctx.Err() != nil:
		return nil
	case s.isInvalid(grid):
		return nil
	case grid.isCompleted():
		return []Grid{*grid}
	}

	// Apply deduction logic
	ok := s.deduction(ctx, grid)
	if !ok {
		return nil
	}

	// Checks if it is now invalid or if is filled
	switch {
	case ctx.Err() != nil:
		return nil
	case s.isInvalid(grid):
		return nil
	case grid.isCompleted():
		return []Grid{*grid}
	}

	// A recursive method for generating solutions
	for row := 0; row < Size; row++ {
		for column := 0; column < Size; column++ {
			if ctx.Err() != nil {
				break
			}

			// Skips value if already set
			if grid.values[row][column] > 0 {
				continue
			}

			// Loops through possible values
			solutions := make([]Grid, 0)
			for value := 1; value <= Size; value++ {
				grid.values[row][column] = value
				if !s.isInvalid(grid) {
					clone := grid.clone()
					solutions = append(solutions, s.solve(ctx, &clone)...)
				}
				grid.values[row][column] = 0
			}

			return solutions
		}
	}

	return nil
}
