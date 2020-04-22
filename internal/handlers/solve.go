package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/PeterEFinch/sudoku-solver/api"
	"github.com/PeterEFinch/sudoku-solver/internal/sudoku"
)

// Solve handles requests solving a sudoku.
func Solve(rw http.ResponseWriter, req *http.Request) {
	reqBs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Err(err).Msg("failed to unmarshal request")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = req.Body.Close()
	if err != nil {
		log.Err(err).Msg("failed to close body")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	request := &api.SolveRequest{}
	err = json.Unmarshal(reqBs, request)
	if err != nil {
		log.Err(err).Msg("failed to unmarshal request")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	err = validateSolveRequest(request)
	if err != nil {
		log.Err(err).Msg("bad request")
		rw.WriteHeader(http.StatusBadRequest)
		_, _ = rw.Write([]byte(err.Error()))
		return
	}

	timeout := time.Duration(request.TimeoutMs) * time.Millisecond
	ctx, cancel := context.WithTimeout(req.Context(), timeout)
	defer cancel()

	response, err := solve(ctx, request)
	if err != nil {
		log.Err(err).Msg("solver failed to run")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	resBs, err := json.Marshal(response)
	if err != nil {
		log.Err(err).Msg("failed to marshall response")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = rw.Write(resBs)
	if err != nil {
		log.Err(err).Msg("failed to write response")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// solve converts the request to a grid which can be solved
// by the sudoku solver.
func solve(ctx context.Context, request *api.SolveRequest) (*api.SolveResponse, error) {
	grid := sudoku.Grid{}
	for r, row := range request.Grid {
		for c, entry := range row {
			if entry > 0 {
				err := grid.Set(r, c, entry)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	solutions := sudoku.StandardSolve(ctx, grid)

	response := &api.SolveResponse{
		Completed: ctx.Err() == nil,
		Solutions: make([]api.Grid, len(solutions)),
	}
	for i, s := range solutions {
		response.Solutions[i] = make([][]int, sudoku.Size)
		for row := 0; row < sudoku.Size; row++ {
			response.Solutions[i][row] = make([]int, sudoku.Size)
			for column := 0; column < sudoku.Size; column++ {
				response.Solutions[i][row][column], _ = s.Get(row, column)
			}
		}
	}

	return response, nil
}

// validateSolveRequest validates a solve request.
func validateSolveRequest(request *api.SolveRequest) error {
	if len(request.Grid) != sudoku.Size {
		return fmt.Errorf("invalid number of rows %d (expected %d)", len(request.Grid), sudoku.Size)
	}

	for r, row := range request.Grid {
		if len(row) != sudoku.Size {
			return fmt.Errorf("row %d has invalid number of columns %d (expected %d)", r, len(row), sudoku.Size)
		}

		for c, entry := range row {
			if entry < 0 || entry > sudoku.Size {
				return fmt.Errorf("invalid entry: %d at position (%d, %d) ", entry, r, c)
			}
		}
	}

	return nil
}
