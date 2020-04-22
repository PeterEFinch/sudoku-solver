package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/PeterEFinch/sudoku-solver/internal/handlers"
)

func main() {
	address := flag.String("address", ":8080", "the address for the server")
	flag.Parse()

	// Adds handlers to mux
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handlers.Health)
	mux.HandleFunc("/solve", handlers.Solve)

	// Creates server
	server := &http.Server{
		Addr:    *address,
		Handler: mux,
	}

	// Allows the server to gracefully shutdown
	shutdownOnSignal(server, os.Interrupt)

	// Starts server
	log.Info().Str("address", *address).Msg("starting server")
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Str("address", *address).Msg("listen and serve failed")
	}
}

// shutdownOnSignal will shutdown the server when one of the
// given signals is sent. It is not blocking.
func shutdownOnSignal(server *http.Server, signals ...os.Signal) {
	if len(signals) == 0 {
		log.Warn().Msg("server is not listening to any shutdown signals")
		return
	}

	// Creates channel for signals and
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, signals...)

	go func() {
		// Waits for signal
		sig := <-signalCh
		signal.Stop(signalCh)
		close(signalCh)
		log.Info().Str("signal", sig.String()).Msg("server shutting down")

		// Gracefully shuts down the server
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		err := server.Shutdown(ctx)
		if err != nil {

		}

		// Cancel context
		cancel()
	}()
}
