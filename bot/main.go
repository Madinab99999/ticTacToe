package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/talgat-ruby/exercises-go/exercise4/bot/joinGame"
	"github.com/talgat-ruby/exercises-go/exercise4/bot/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ready, srv := server.StartServer()
	<-ready

	if err := joinGame.JoinGame(ctx); err != nil {
		log.Fatalf("Failed to join the game: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
}
