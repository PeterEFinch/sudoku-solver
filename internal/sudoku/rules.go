package sudoku

// Rule represents a sudoku rule, whether it be a standard
// rule, such as each row must have the numbers 1 to 9, or
// a custom rule, such as the diagonal rule where the
// diagonals must also have
type Rule interface {
	// isInvalid must return true if and only if the grid
	// violates this rule.
	isInvalid(grid *Grid) bool

	// deduction must return the first valid entry it finds
	// and true if either an entry was found or the possibilities
	// where changed.
	deduction(grid *Grid, possibilities *possibilities) (*entry, bool)

	// restrict must restrict the possibilities (in line
	// with this rule) based on the entry given.
	restrict(entry *entry, possibilities *possibilities)
}

type entry struct {
	row    int
	column int
	value  int
}

// region TrivialRule

func TrivialRule() Rule {
	return trivialRule{}
}

type trivialRule struct{}

func (trivialRule) isInvalid(_ *Grid) bool {
	return false
}

func (trivialRule) deduction(grid *Grid, possibilities *possibilities) (*entry, bool) {
	for row := 0; row < Size; row++ {
		for column := 0; column < Size; column++ {
			if grid.values[row][column] == 0 && len(possibilities[row][column]) == 1 {
				for key := range possibilities[row][column] {
					return &entry{
						row:    row,
						column: column,
						value:  key,
					}, true
				}
			}
		}
	}

	return nil, false
}

func (trivialRule) restrict(_ *entry, _ *possibilities) {}

// endregion

// region Row Rule

// RowRule
func RowRule() Rule {
	return rowRule{}
}

type rowRule struct{}

func (rowRule) isInvalid(grid *Grid) bool {
	var entries [Size]bool
	for row := 0; row < Size; row++ {
		// Resets entries
		for i := range entries {
			entries[i] = false
		}

		// Loops through columns
		for column := 0; column < Size; column++ {
			value := grid.values[row][column]

			switch {
			case value == 0:
			case entries[value-1]:
				return true
			default:
				entries[value-1] = true
			}
		}
	}

	return false
}

func (rowRule) deduction(grid *Grid, possibilities *possibilities) (*entry, bool) {
	entry := singlePositionLogic(grid, possibilities, fixedRowConverter)
	if entry != nil {
		return entry, true
	}

	applied := false
	for n := 2; n < 9; n++ {
		if limitPossibilitiesLogic(grid, possibilities, fixedRowConverter, n) {
			applied = true
		}
	}

	return nil, applied
}

func (rowRule) restrict(entry *entry, possibilities *possibilities) {
	for column := 0; column < Size; column++ {
		if column != entry.column {
			possibilities.remove(entry.row, column, entry.value)
		}
	}
}

// endregion

// region Column Rule

func ColumnRule() Rule {
	return columnRule{}
}

type columnRule struct{}

func (columnRule) isInvalid(grid *Grid) bool {
	var entries [Size]bool
	for column := 0; column < Size; column++ {
		// Resets entries
		for i := range entries {
			entries[i] = false
		}

		// Loops through rows
		for row := 0; row < Size; row++ {
			value := grid.values[row][column]

			switch {
			case value == 0:
			case entries[value-1]:
				return true
			default:
				entries[value-1] = true
			}
		}
	}

	return false
}

func (columnRule) deduction(grid *Grid, possibilities *possibilities) (*entry, bool) {
	entry := singlePositionLogic(grid, possibilities, fixedColumnConverter)
	if entry != nil {
		return entry, true
	}

	applied := false
	for n := 2; n < 9; n++ {
		if limitPossibilitiesLogic(grid, possibilities, fixedColumnConverter, n) {
			applied = true
		}
	}

	return nil, applied
}

func (columnRule) restrict(entry *entry, possibilities *possibilities) {
	for row := 0; row < Size; row++ {
		if row != entry.row {
			possibilities.remove(row, entry.column, entry.value)
		}
	}
}

// endregion

// region Square Rule

func SquareRule() Rule {
	return squareRule{}
}

type squareRule struct{}

func (squareRule) isInvalid(grid *Grid) bool {
	var entries [Size]bool
	for outer := 0; outer < Size; outer++ {
		// Resets entries
		for i := range entries {
			entries[i] = false
		}

		for inner := 0; inner < Size; inner++ {
			row, column := squareToStandard(outer, inner)
			value := grid.values[row][column]

			switch {
			case value == 0:
			case entries[value-1]:
				return true
			default:
				entries[value-1] = true
			}
		}
	}

	return false
}

func (squareRule) deduction(grid *Grid, possibilities *possibilities) (*entry, bool) {
	entry := singlePositionLogic(grid, possibilities, squareToStandard)
	if entry != nil {
		return entry, true
	}

	applied := false
	for n := 2; n < 9; n++ {
		if limitPossibilitiesLogic(grid, possibilities, squareToStandard, n) {
			applied = true
		}
	}

	return nil, applied
}

func (squareRule) restrict(entry *entry, possibilities *possibilities) {
	outer, inner := standardToSquare(entry.row, entry.column)
	for variable := 0; variable < Size; variable++ {
		if variable != inner {
			row, column := squareToStandard(outer, variable)
			possibilities.remove(row, column, entry.value)
		}
	}
}

// endregion

// region Chess King Rule

func ChessKingRule() Rule {
	return chessKingRule{}
}

type chessKingRule struct{}

func (chessKingRule) isInvalid(grid *Grid) bool {
	for row := 0; row < Size; row++ {
		for column := 0; column < Size; column++ {
			if grid.isCellEmpty(row, column) {
				continue
			}

			value := grid.values[row][column]
			for r := row - 1; r <= row+1; r++ {
				for c := column - 1; c <= column+1; c++ {
					if r != row && c != column &&
						isValidPair(r, c) &&
						grid.values[r][c] == value {
						return true
					}
				}
			}
		}
	}

	return false
}

func (chessKingRule) deduction(_ *Grid, _ *possibilities) (*entry, bool) {
	return nil, false
}

func (chessKingRule) restrict(entry *entry, possibilities *possibilities) {
	possibilities.safeRemove(entry.row+1, entry.column-1, entry.value)
	possibilities.safeRemove(entry.row+1, entry.column, entry.value)
	possibilities.safeRemove(entry.row+1, entry.column+1, entry.value)
	possibilities.safeRemove(entry.row, entry.column+1, entry.value)
	possibilities.safeRemove(entry.row-1, entry.column+1, entry.value)
	possibilities.safeRemove(entry.row-1, entry.column, entry.value)
	possibilities.safeRemove(entry.row-1, entry.column-1, entry.value)
	possibilities.safeRemove(entry.row, entry.column-1, entry.value)
}

// endregion

// region Chess King Rule

func ChessKnightRule() Rule {
	return chessKnightRule{}
}

type chessKnightRule struct{}

func (chessKnightRule) isInvalid(grid *Grid) bool {
	for row := 0; row < Size; row++ {
		for column := 0; column < Size; column++ {
			if grid.isCellEmpty(row, column) {
				continue
			}

			value := grid.values[row][column]
			for near := -1; near <= 1; near += 2 {
				for far := -2; far <= 2; far += 4 {
					r, c := row+near, column+far
					if isValidPair(r, c) && grid.values[r][c] == value {
						return true
					}

					r, c = row+far, column+near
					if isValidPair(r, c) && grid.values[r][c] == value {
						return true
					}
				}
			}
		}
	}

	return false
}

func (chessKnightRule) deduction(_ *Grid, _ *possibilities) (*entry, bool) {
	return nil, false
}

func (chessKnightRule) restrict(entry *entry, possibilities *possibilities) {
	possibilities.safeRemove(entry.row+2, entry.column-1, entry.value)
	possibilities.safeRemove(entry.row+2, entry.column+1, entry.value)
	possibilities.safeRemove(entry.row-2, entry.column-1, entry.value)
	possibilities.safeRemove(entry.row-2, entry.column+1, entry.value)
	possibilities.safeRemove(entry.row+1, entry.column+2, entry.value)
	possibilities.safeRemove(entry.row-1, entry.column+2, entry.value)
	possibilities.safeRemove(entry.row+1, entry.column-2, entry.value)
	possibilities.safeRemove(entry.row-1, entry.column-2, entry.value)
}

// endregion

// region Helpers

func isValidPair(row, column int) bool {
	if row >= 0 && row < Size &&
		column >= 0 && column < Size {
		return true
	}

	return false
}

func singlePositionLogic(grid *Grid, possibilities *possibilities, convert indicesConverter) *entry {
	// Loops through fixed index and value
	for fixed := 0; fixed < Size; fixed++ {
		for value := 1; value <= Size; value++ {
			count := 0
			candidateRow := 0
			candidateColumn := 0

			// Loops through variable index
			for variable := 0; variable < Size; variable++ {
				row, column := convert(fixed, variable)

				// Checks if value is already used
				if grid.values[row][column] == value {
					count = 0
					break
				}

				// Checks if value is possible
				if _, ok := possibilities[row][column][value]; ok {
					count++
					if count > 1 {
						break
					}

					candidateRow = row
					candidateColumn = column
				}
			}

			// Checks if there was only once place to place the value, the value is place in the candidate row/column
			if count == 1 {
				return &entry{
					row:    candidateRow,
					column: candidateColumn,
					value:  value,
				}
			}
		}
	}

	return nil
}

func limitPossibilitiesLogic(grid *Grid, possibilities *possibilities, convert indicesConverter, n int) bool {
	applied := false

	for fixed := 0; fixed < Size; fixed++ {
		for variable := 0; variable < Size; variable++ {
			// Skips entries with incorrect number of possibilities
			row, column := convert(fixed, variable)
			if len(possibilities[row][column]) != n {
				continue
			}

			// Initialises variables
			matches := 0
			irrelevant := 0
			var ignore [Size]bool

			// Determines the number of entries in row with same possibilities
			for other := 0; other < Size; other++ {
				r, c := convert(fixed, other)

				// Counts the number of the entries with the same possibilities
				if possibilities.isSame(row, column, r, c) {
					matches++
					ignore[other] = true
					continue
				}

				// Counts number of irrelevant entries
				if grid.values[r][c] != 0 || !possibilities.isOverlapping(row, column, r, c) {
					irrelevant++
					ignore[other] = true
					continue
				}
			}

			// Stops if number of matches doesn't is not n or if nothing can happen
			if matches != n || matches+irrelevant == Size {
				continue
			}

			// Changes the other possibilities
			applied = true
			m := possibilities[row][column]
			for other, skip := range ignore {
				if skip {
					continue
				}

				r, c := convert(fixed, other)
				for k := range m {
					possibilities.remove(r, c, k)
				}
			}
		}
	}

	return applied
}

// endregion
