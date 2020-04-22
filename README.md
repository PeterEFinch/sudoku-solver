# Sudoku Solver

Sudoku Solver is a server provides REST calls to solve sudoku
puzzles. 

## Getting started

To download the repository use the command:
```
go get github.com/PeterEFinch/sudoku-solver
```
The server can be started using running: 
```
go run .cli/server/main.go
``` 
However, we recommend using docker and docker-compose.
To start the server using docker compose simply run
```
docker-compose up --build
```
The server is now running at *localhost:8080*.

## Example REST Call
To solve a sudoku puzzle calls should be made to 
*localhost:8080/solve*. For example try the call:
```
{
  "timeout_ms": 5000,
  "grid": [
      [8, 0, 0, 0, 0, 0, 0, 0, 0],
      [0, 0, 3, 6, 0, 0, 0, 0, 0],
      [0, 7, 0, 0, 9, 0, 2, 0, 0],
      [0, 5, 0, 0, 0, 7, 0, 0, 0],
      [0, 0, 0, 0, 4, 5, 7, 0, 0],
      [0, 0, 0, 1, 0, 0, 0, 3, 0],
      [0, 0, 1, 0, 0, 0, 0, 6, 8],
      [0, 0, 8, 5, 0, 0, 0, 1, 0],
      [0, 9, 0, 0, 0, 0, 4, 0, 0]
  ]
}
```
