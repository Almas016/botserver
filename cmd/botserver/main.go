package main

import (
	"botserver/internal/usecase"
	"botserver/pkg/natsclient"

	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	natsClient := natsclient.NewClient()
	botUsecase := usecase.NewBotUsecase(natsClient)

	go func() {
		if err := botUsecase.StartBotServer(); err != nil {
			log.Fatalf("Failed to start Bot Server: %v", err)
		}
	}()

	// gracefully shut down
	shutDown()
}

func shutDown() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals
}
