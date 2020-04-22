package sudoku

type indicesConverter func(int, int) (int, int)

func squareToStandard(outer, inner int) (int, int) {
	i, j := outer%3, outer/3
	k, l := inner%3, inner/3
	row := 3*i + k
	column := 3*j + l
	return row, column
}

func standardToSquare(row, column int) (int, int) {
	i, j := row/3, column/3
	k, l := row%3, column%3
	outer := 3*j + i
	inner := 3*l + k
	return outer, inner
}

func fixedSquareConverter(fixed, variable int) (int, int) {
	return squareToStandard(fixed, variable)
}

func fixedRowConverter(fixed, variable int) (int, int) {
	return fixed, variable
}

func fixedColumnConverter(fixed, variable int) (int, int) {
	return variable, fixed
}
