package ticTacToe

type State string

const (
	StatePending  State = "PENDING"
	StateRunning  State = "RUNNING"
	StateFinished State = "FINISHED"
)

type Game struct {
	State   State     `json:"state"`
	Players []*Player `json:"players"`
	Matches []*Match  `json:"matches"`
}

type Match struct {
	Players [2]*Player `json:"players"`
	Rounds  []*Round   `json:"rounds"`
}

type Round struct {
	Players [2]*Player `json:"players"`
	Board   *Board     `json:"board"`
	Moves   []*Move    `json:"moves"`
	Winner  *Player    `json:"winner"`
}

type Move struct {
	Player *Player `json:"players"`
	Board  *Board  `json:"board"`
}

const (
	Cols = 3
	Rows = 3
)

type Board [Cols * Rows]Token
