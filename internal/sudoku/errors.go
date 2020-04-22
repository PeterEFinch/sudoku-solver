package sudoku

const (
	ErrTimeout Error = "timeout"
	//ErrGuessRequired Error = "guess_required"

	ErrInvalidValue     Error = "invalid_value"
	ErrInvalidRow       Error = "invalid_row"
	ErrInvalidColumn    Error = "invalid_column"
	ErrSquareAlreadySet Error = "square_already_set"
)

type Error string

func (e Error) Error() string {
	return string(e)
}
