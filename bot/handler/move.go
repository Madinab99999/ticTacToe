package handler

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/talgat-ruby/exercises-go/exercise4/bot/pkg/httputils/request"
	"github.com/talgat-ruby/exercises-go/exercise4/bot/ticTacToe"
)

type RequestMove struct {
	Board *ticTacToe.Board `json:"board"`
	Token ticTacToe.Token  `json:"token"`
}

type ResponseMove struct {
	Index int `json:"index"`
}

func Move(w http.ResponseWriter, r *http.Request) {
	var req RequestMove
	if err := request.JSON(w, r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	position := calculateMove(req.Board, req.Token)
	resp := ResponseMove{Index: position}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to respond to move", http.StatusInternalServerError)
	}
}

func calculateMove(board *ticTacToe.Board, token ticTacToe.Token) int {
	// Если на доске только одна свободная ячейка
	if countEmptyCells(board) == 1 {
		return getAnyOccupiedCell(board)
	}

	if isBoardEmpty(board) {
		return getRandomEmptyCell(board)
	}

	for i, cell := range board {
		if cell == ticTacToe.TokenEmpty {
			board[i] = token
			if checkWinner(board, token) {
				board[i] = ticTacToe.TokenEmpty
				return i
			}
			board[i] = ticTacToe.TokenEmpty
		}
	}

	opponent := getOpponentToken(token)
	for i, cell := range board {
		if cell == ticTacToe.TokenEmpty {
			board[i] = opponent
			if checkWinner(board, opponent) {
				board[i] = ticTacToe.TokenEmpty
				return i
			}
			board[i] = ticTacToe.TokenEmpty
		}
	}

	for i, cell := range board {
		if cell == ticTacToe.TokenEmpty {
			return i
		}
	}

	return -1
}

func countEmptyCells(board *ticTacToe.Board) int {
	count := 0
	for _, cell := range board {
		if cell == ticTacToe.TokenEmpty {
			count++
		}
	}
	return count
}

func getAnyOccupiedCell(board *ticTacToe.Board) int {
	for i, cell := range board {
		if cell != ticTacToe.TokenEmpty {
			return i
		}
	}
	return -1
}

func isBoardEmpty(board *ticTacToe.Board) bool {
	for _, cell := range board {
		if cell != ticTacToe.TokenEmpty {
			return false
		}
	}
	return true
}

func getRandomEmptyCell(board *ticTacToe.Board) int {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	emptyCells := []int{}

	for i, cell := range board {
		if cell == ticTacToe.TokenEmpty {
			emptyCells = append(emptyCells, i)
		}
	}

	if len(emptyCells) == 0 {
		return -1
	}
	return emptyCells[rng.Intn(len(emptyCells))]
}

func checkWinner(board *ticTacToe.Board, token ticTacToe.Token) bool {
	winPatterns := [][3]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, // Rows
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8}, // Columns
		{0, 4, 8}, {2, 4, 6}, // Diagonals
	}

	for _, pattern := range winPatterns {
		if board[pattern[0]] == token && board[pattern[1]] == token && board[pattern[2]] == token {
			return true
		}
	}
	return false
}

func getOpponentToken(token ticTacToe.Token) ticTacToe.Token {
	if token == ticTacToe.TokenX {
		return ticTacToe.TokenO
	}
	return ticTacToe.TokenX
}
