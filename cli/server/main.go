package main

import (
	"context"
	"errors"
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

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handlers.Health)
	mux.HandleFunc("/solve", handlers.Solve)
	server := &http.Server{
		Addr:    *address,
		Handler: mux,
	}

	shutdownOnSignal(server, os.Interrupt)

	log.Info().Str("address", *address).Msg("starting server")
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Error().Err(err).Str("address", *address).Msg("listen and serve failed")
		os.Exit(1)
	}
}

// shutdownOnSignal will attempt to gracefully shutdown the server when one of
// the given signals is sent.
//
// The waiting for a signal is done in a go routine to ensure that this method
// is not blocking.
func shutdownOnSignal(server *http.Server, signals ...os.Signal) {
	if len(signals) == 0 {
		log.Warn().Msg("server is not listening to any shutdown signals")
		return
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, signals...)

	go func() {
		sig := <-signalCh
		signal.Stop(signalCh)
		close(signalCh)
		log.Info().Str("signal", sig.String()).Msg("server shutting down")

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		if err := server.Shutdown(ctx); err != nil {
			log.Warn().Err(err).Msg("failed to shutdown server")
		}
		cancel()
	}()
}
