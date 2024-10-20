package handler

import (
	"encoding/json"
	"log"
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
		log.Printf("Failed to encode response: %v", err)
	}
}

func calculateMove(board *ticTacToe.Board, token ticTacToe.Token) int {
	// Проверяем, можем ли мы выиграть
	for i, cell := range *board {
		if cell == ticTacToe.TokenEmpty {
			(*board)[i] = token
			if checkWinner(board) == token {
				(*board)[i] = ticTacToe.TokenEmpty
				return i
			}
			(*board)[i] = ticTacToe.TokenEmpty
		}
	}

	// Проверяем, может ли противник выиграть, и блокируем его
	opponent := getOpponentToken(token)
	for i, cell := range *board {
		if cell == ticTacToe.TokenEmpty {
			(*board)[i] = opponent
			if checkWinner(board) == opponent {
				(*board)[i] = ticTacToe.TokenEmpty
				return i
			}
			(*board)[i] = ticTacToe.TokenEmpty
		}
	}

	// Если центр пустой, займем его
	if (*board)[4] == ticTacToe.TokenEmpty {
		return 4
	}

	// Если углы пустые, займем угол
	for _, i := range []int{0, 2, 6, 8} {
		if (*board)[i] == ticTacToe.TokenEmpty {
			return i
		}
	}

	// Иначе займем любую пустую клетку
	for i, cell := range *board {
		if cell == ticTacToe.TokenEmpty {
			return i
		}
	}

	return -1 // Если нет доступных ходов
}

func checkWinner(board *ticTacToe.Board) ticTacToe.Token {
	winPatterns := [][3]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, // Ряды
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8}, // Столбцы
		{0, 4, 8}, {2, 4, 6}, // Диагонали
	}

	for _, pattern := range winPatterns {
		if (*board)[pattern[0]] != ticTacToe.TokenEmpty &&
			(*board)[pattern[0]] == (*board)[pattern[1]] &&
			(*board)[pattern[1]] == (*board)[pattern[2]] {
			return (*board)[pattern[0]] // Возвращаем токен победителя
		}
	}
	return ticTacToe.TokenEmpty // Если победителя нет
}

func getOpponentToken(token ticTacToe.Token) ticTacToe.Token {
	if token == ticTacToe.TokenX {
		return ticTacToe.TokenO
	}
	return ticTacToe.TokenX
}

// Алгоритм Minimax с альфа-бета отсечением
func minimax(board *ticTacToe.Board, depth int, isMaximizing bool, player ticTacToe.Token, alpha, beta int) int {
	winner := checkWinner(board)
	if winner == player {
		return 10 - depth
	} else if winner != ticTacToe.TokenEmpty {
		return depth - 10
	} else if isBoardFull(board) {
		return 0
	}

	opponent := getOpponentToken(player)

	if isMaximizing {
		maxEval := -1000000
		for i := 0; i < len(*board); i++ {
			if (*board)[i] == ticTacToe.TokenEmpty {
				(*board)[i] = player
				eval := minimax(board, depth+1, false, player, alpha, beta)
				(*board)[i] = ticTacToe.TokenEmpty
				maxEval = max(maxEval, eval)
				alpha = max(alpha, eval)
				if beta <= alpha {
					break
				}
			}
		}
		return maxEval
	} else {
		minEval := 1000000
		for i := 0; i < len(*board); i++ {
			if (*board)[i] == ticTacToe.TokenEmpty {
				(*board)[i] = opponent
				eval := minimax(board, depth+1, true, player, alpha, beta)
				(*board)[i] = ticTacToe.TokenEmpty
				minEval = min(minEval, eval)
				beta = min(beta, eval)
				if beta <= alpha {
					break
				}
			}
		}
		return minEval
	}
}

func isBoardFull(board *ticTacToe.Board) bool {
	for _, cell := range *board {
		if cell == ticTacToe.TokenEmpty {
			return false
		}
	}
	return true
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
