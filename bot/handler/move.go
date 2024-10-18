package handler

import (
	"encoding/json"
	"net/http"

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

func (h *Handler) Move(w http.ResponseWriter, r *http.Request) {
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

	if board[4] == ticTacToe.TokenEmpty {
		return 4
	}

	for _, i := range []int{0, 2, 6, 8} {
		if board[i] == ticTacToe.TokenEmpty {
			return i
		}
	}

	for i, cell := range board {
		if cell == ticTacToe.TokenEmpty {
			return i
		}
	}

	return -1
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
