package api

type Grid [][]int

type SolveRequest struct {
	TimeoutMs int  `json:"timeout_ms,omitempty"`
	Grid      Grid `json:"grid,omitempty"`
}

type SolveResponse struct {
	Completed bool
	Solutions []Grid
}
