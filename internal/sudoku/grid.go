package sudoku

import (
	"fmt"
	"strings"
)

const (
	Size = 9
)

type Grid struct {
	values [Size][Size]int
}

func (grid *Grid) Clear(row, column int) error {
	// Checks for possible errors
	switch {
	case row < 0 || row >= Size:
		return ErrInvalidRow
	case column < 0 || column > Size:
		return ErrInvalidColumn
	}

	// Clears value
	if grid.values[row][column] != 0 {
		grid.values[row][column] = 0
	}

	return nil

}

func (grid *Grid) Set(row, column, value int) error {
	// Checks for possible errors
	switch {
	case value < 1 || value > Size:
		return ErrInvalidValue
	case row < 0 || row >= Size:
		return ErrInvalidRow
	case column < 0 || column >= Size:
		return ErrInvalidColumn
	case grid.values[row][column] != 0:
		return ErrSquareAlreadySet
	}

	// Sets value
	grid.values[row][column] = value
	return nil
}

func (grid *Grid) Get(row, column int) (int, error) {
	// Checks for possible errors
	switch {
	case row < 0 || row >= Size:
		return 0, ErrInvalidRow
	case column < 0 || column >= Size:
		return 0, ErrInvalidColumn
	}

	return grid.values[row][column], nil
}

func (grid *Grid) String() string {
	s := ""
	for _, row := range grid.values {
		rowStr := strings.Replace(fmt.Sprintf("%v", row), "0", "-", -1)
		rowStr = rowStr[1 : len(rowStr)-1]
		s += rowStr + "\n"
	}
	return s
}

func (grid *Grid) clone() Grid {
	clone := Grid{}
	for i, row := range grid.values {
		copy(clone.values[i][:], row[:])
	}

	return clone
}

func (grid *Grid) isCellEmpty(row, column int) bool {
	return grid.values[row][column] == 0
}

func (grid *Grid) isCellFilled(row, column int) bool {
	return grid.values[row][column] != 0
}

func (grid *Grid) isCompleted() bool {
	for row := 0; row < Size; row++ {
		for column := 0; column < Size; column++ {
			if grid.values[row][column] == 0 {
				return false
			}
		}
	}

	return true
}

// region Example Grids

func DifficultExampleGrid() Grid {
	return Grid{
		values: [9][9]int{
			{8, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 3, 6, 0, 0, 0, 0, 0},
			{0, 7, 0, 0, 9, 0, 2, 0, 0},
			{0, 5, 0, 0, 0, 7, 0, 0, 0},
			{0, 0, 0, 0, 4, 5, 7, 0, 0},
			{0, 0, 0, 1, 0, 0, 0, 3, 0},
			{0, 0, 1, 0, 0, 0, 0, 6, 8},
			{0, 0, 8, 5, 0, 0, 0, 1, 0},
			{0, 9, 0, 0, 0, 0, 4, 0, 0},
		},
	}

}

// endregion
