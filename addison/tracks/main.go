package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tracks/rest"
	"tracks/tracks"
)

const (
	serverTimeoutInSeconds = 5 * time.Second
	databaseFilename       = "file:tracks.db"
	restHost               = "localhost"
	restPort               = 3000
)

func main() {
	trackstore, err := tracks.NewTrackStore(databaseFilename)
	if err != nil {
		log.Print(err)
		return
	}
	mux := rest.NewMux(trackstore)
	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", restHost, restPort),
		Handler: mux,
	}

	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Print(err)
			return
		}
	}()

	log.Printf("tracks started on %s", server.Addr)

	<-interruptChannel

	log.Printf("tracks stopped on %s", server.Addr)

	ctx, cancel := context.WithTimeout(context.Background(), serverTimeoutInSeconds)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Print(err)
		return
	}
}
