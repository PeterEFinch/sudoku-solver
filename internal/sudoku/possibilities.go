package sudoku

import (
	"fmt"
)

type possibilities [Size][Size]map[int]struct{}

func (p *possibilities) display() {
	for row := 0; row < Size; row++ {
		for column := 0; column < Size; column++ {
			if len(p[row][column]) > 1 {
				fmt.Println(row, column, p[row][column])
			}
		}
	}
}

func (p *possibilities) initialise() {
	for row := 0; row < Size; row++ {
		for column := 0; column < Size; column++ {
			p[row][column] = map[int]struct{}{1: {}, 2: {}, 3: {}, 4: {}, 5: {}, 6: {}, 7: {}, 8: {}, 9: {}}
		}
	}

}

func (p *possibilities) isInvalid() bool {
	for row := 0; row < Size; row++ {
		for column := 0; column < Size; column++ {
			if len(p[row][column]) == 0 {
				return true
			}
		}
	}

	return false
}

func (p *possibilities) isOverlapping(rowA, columnA, rowB, columnB int) bool {
	for key := range p[rowA][columnA] {
		if _, ok := p[rowB][columnB][key]; ok {
			return true
		}
	}

	for key := range p[rowB][columnB] {
		if _, ok := p[rowA][columnA][key]; ok {
			return true
		}
	}

	return false
}

func (p *possibilities) isSame(rowA, columnA, rowB, columnB int) bool {
	if len(p[rowA][columnA]) != len(p[rowB][columnB]) {
		return false
	}

	for key := range p[rowA][columnA] {
		if _, ok := p[rowB][columnB][key]; !ok {
			return false
		}
	}

	return true
}

func (p *possibilities) safeRemove(row, column, value int) {
	if isValidPair(row, column) {
		delete(p[row][column], value)
	}
}

func (p *possibilities) set(row, column, value int) {
	p[row][column] = map[int]struct{}{value: {}}
}

func (p *possibilities) remove(row, column, value int) {
	delete(p[row][column], value)
}
