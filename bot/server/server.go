package server

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/talgat-ruby/exercises-go/exercise4/bot/handler"
)

type readyListener struct {
	net.Listener
	ready chan struct{}
	once  sync.Once
}

func (l *readyListener) Accept() (net.Conn, error) {
	l.once.Do(func() { close(l.ready) })
	return l.Listener.Accept()
}

// Start the server and handle requests using http.NewServeMux()
func StartServer() (<-chan struct{}, *http.Server) {
	ready := make(chan struct{})

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("PORT")))
	if err != nil {
		panic(err)
	}

	list := &readyListener{Listener: listener, ready: ready}
	srv := &http.Server{
		IdleTimeout: 2 * time.Minute,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", handler.Ping)
	mux.HandleFunc("POST /move", handler.Move)
	srv.Handler = mux

	go func() {
		err := srv.Serve(list)
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	return ready, srv
}
