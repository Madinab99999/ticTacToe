package joinGame

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/talgat-ruby/exercises-go/exercise4/bot/ticTacToe"
)

type RequestJoin struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func NewPlayer() *ticTacToe.Player {
	name := os.Getenv("NAME")
	if name == "" {
		name = "Player"
	}
	port := os.Getenv("PORT")
	url := fmt.Sprintf("http://localhost:%s", port)

	return ticTacToe.New(name, url)
}

func JoinGame(ctx context.Context) error {
	port := os.Getenv("PORT")
	botURL := fmt.Sprintf("http://localhost:%s", port)
	player := NewPlayer()
	reqBody := RequestJoin{Name: player.Name, URL: botURL}

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal join request: %w", err)
	}

	client := &http.Client{}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"http://localhost:4444/join",
		bytes.NewBuffer(reqBodyBytes))

	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to join the game: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	log.Println("Successfully joined the game")
	return nil
}
