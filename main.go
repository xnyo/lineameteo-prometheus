package main

import (
	"context"
	"github.com/xnyo/lineameteo-prometheus/config"
	"github.com/xnyo/lineameteo-prometheus/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const httpGraceExitTimeout = 30 * time.Second

func _main() error {
	srv := service.NewService(
		config.HTTPListen.Get(),
		strings.Split(config.WantedLocationIDs.Get(), ","),
	)
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		shutdownCtx, canc := context.WithTimeout(serverCtx, httpGraceExitTimeout)
		defer canc()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		if err := srv.HTTPServer.Shutdown(shutdownCtx); err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	log.Println("http server started")
	err := srv.HTTPServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
	return nil
}

func main() {
	if err := _main(); err != nil {
		panic(err)
	}
}
