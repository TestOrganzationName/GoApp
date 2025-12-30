package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Stone int

const (
	Empty Stone = iota
	Black
	White
)

func (s Stone) String() string {
	switch s {
	case Black:
		return "black"
	case White:
		return "white"
	default:
		return "empty"
	}
}

type Board struct {
	Size   int       `json:"size"`
	Grid   [][]Stone `json:"grid"`
	Turn   Stone     `json:"turn"`
	Passes int       `json:"passes"`
	mu     sync.RWMutex
}

func NewBoard(size int) *Board {
	grid := make([][]Stone, size)
	for i := range grid {
		grid[i] = make([]Stone, size)
	}
	return &Board{
		Size: size,
		Grid: grid,
		Turn: Black,
	}
}

func (b *Board) IsValidMove(row, col int) bool {
	if row < 0 || row >= b.Size || col < 0 || col >= b.Size {
		return false
	}
	return b.Grid[row][col] == Empty
}

func (b *Board) PlaceStone(row, col int) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.IsValidMove(row, col) {
		return false
	}

	b.Grid[row][col] = b.Turn

	// Remove captured opponent stones
	opponent := White
	if b.Turn == White {
		opponent = Black
	}

	// Check all adjacent positions for captures
	directions := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for _, dir := range directions {
		newRow, newCol := row+dir[0], col+dir[1]
		if b.isInBounds(newRow, newCol) && b.Grid[newRow][newCol] == opponent {
			if !b.hasLiberties(newRow, newCol, make(map[[2]int]bool)) {
				b.removeGroup(newRow, newCol)
			}
		}
	}

	// Check if the placed stone group has liberties (suicide rule)
	if !b.hasLiberties(row, col, make(map[[2]int]bool)) {
		b.Grid[row][col] = Empty // Remove the stone
		return false
	}

	b.Passes = 0
	b.nextTurn()
	return true
}

func (b *Board) isInBounds(row, col int) bool {
	return row >= 0 && row < b.Size && col >= 0 && col < b.Size
}

func (b *Board) hasLiberties(row, col int, visited map[[2]int]bool) bool {
	if visited[[2]int{row, col}] {
		return false
	}
	visited[[2]int{row, col}] = true

	stone := b.Grid[row][col]
	directions := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	for _, dir := range directions {
		newRow, newCol := row+dir[0], col+dir[1]
		if !b.isInBounds(newRow, newCol) {
			continue
		}

		if b.Grid[newRow][newCol] == Empty {
			return true // Found a liberty
		}

		if b.Grid[newRow][newCol] == stone {
			if b.hasLiberties(newRow, newCol, visited) {
				return true
			}
		}
	}

	return false
}

func (b *Board) removeGroup(row, col int) {
	stone := b.Grid[row][col]
	if stone == Empty {
		return
	}

	b.Grid[row][col] = Empty
	directions := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	for _, dir := range directions {
		newRow, newCol := row+dir[0], col+dir[1]
		if b.isInBounds(newRow, newCol) && b.Grid[newRow][newCol] == stone {
			b.removeGroup(newRow, newCol)
		}
	}
}

func (b *Board) Pass() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Passes++
	b.nextTurn()
}

func (b *Board) nextTurn() {
	if b.Turn == Black {
		b.Turn = White
	} else {
		b.Turn = Black
	}
}

func (b *Board) IsGameOver() bool {
	return b.Passes >= 2
}

func (b *Board) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for i := range b.Grid {
		for j := range b.Grid[i] {
			b.Grid[i][j] = Empty
		}
	}
	b.Turn = Black
	b.Passes = 0
}

var gameBoard = NewBoard(9)

type MoveRequest struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

type MoveResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Board    *Board `json:"board"`
	GameOver bool   `json:"gameOver"`
}

func handleGetBoard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	gameBoard.mu.RLock()
	defer gameBoard.mu.RUnlock()
	json.NewEncoder(w).Encode(gameBoard)
}

func handleMove(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var move MoveRequest
	if err := json.NewDecoder(r.Body).Decode(&move); err != nil {
		json.NewEncoder(w).Encode(MoveResponse{
			Success: false,
			Message: "Invalid request",
		})
		return
	}

	success := gameBoard.PlaceStone(move.Row, move.Col)
	message := "Move successful"
	if !success {
		message = "Invalid move! The position is occupied or violates game rules."
	}

	gameBoard.mu.RLock()
	gameOver := gameBoard.IsGameOver()
	gameBoard.mu.RUnlock()

	json.NewEncoder(w).Encode(MoveResponse{
		Success:  success,
		Message:  message,
		Board:    gameBoard,
		GameOver: gameOver,
	})
}

func handlePass(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	gameBoard.Pass()
	gameBoard.mu.RLock()
	gameOver := gameBoard.IsGameOver()
	gameBoard.mu.RUnlock()

	json.NewEncoder(w).Encode(MoveResponse{
		Success:  true,
		Message:  "Player passed",
		Board:    gameBoard,
		GameOver: gameOver,
	})
}

func handleReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	gameBoard.Reset()

	json.NewEncoder(w).Encode(MoveResponse{
		Success:  true,
		Message:  "Game reset",
		Board:    gameBoard,
		GameOver: false,
	})
}

func main() {
	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// API endpoints
	http.HandleFunc("/api/board", handleGetBoard)
	http.HandleFunc("/api/move", handleMove)
	http.HandleFunc("/api/pass", handlePass)
	http.HandleFunc("/api/reset", handleReset)

	port := ":8080"
	fmt.Printf("Starting Go Web App on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
