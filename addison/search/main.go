package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"search/audio"
	"search/rest"
	"syscall"
	"time"
)

const (
	serverTimeoutInSeconds = 5 * time.Second
	restHost               = "localhost"
	restPort               = 3001
	audioApiToken          = "test"
)

func main() {
	audioClient, err := audio.NewAuddioClient(audioApiToken)
	if err != nil {
		return
	}
	mux := rest.NewMux(audioClient)
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

	log.Printf("search started on %s", server.Addr)

	<-interruptChannel

	log.Printf("search stopped on %s", server.Addr)

	ctx, cancel := context.WithTimeout(context.Background(), serverTimeoutInSeconds)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return
	}
}
