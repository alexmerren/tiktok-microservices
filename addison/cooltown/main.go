package main

import (
	"context"
	"cooltown/rest"
	"cooltown/search"
	"cooltown/tracks"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	serverTimeoutInSeconds = 5 * time.Second
	restHost               = "localhost"
	restPort               = 3002
)

func main() {
	trackClient, err := tracks.NewTrackClient()
	if err != nil {
		return
	}

	searchClient, err := search.NewSearchClient()
	if err != nil {
		return
	}

	mux := rest.NewMux(trackClient, searchClient)
	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", restHost, restPort),
		Handler: mux,
	}

	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return
		}
	}()

	log.Printf("cooltown started on %s", server.Addr)

	<-interruptChannel

	log.Printf("cooltown stopped on %s", server.Addr)

	ctx, cancel := context.WithTimeout(context.Background(), serverTimeoutInSeconds)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return
	}
}
