package utils

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func RegisterSigtermHandler() {
	// Add signal channel notification handler
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	sig := <-sigChan // Block until a signal is received

	switch sig {
	case os.Interrupt:
		fmt.Println("\nReceived SIGINT (Ctrl+C). Closing y2spot...")
		os.Exit(0)
	case syscall.SIGTERM:
		fmt.Println("\nReceived SIGTERM. Closing y2spot...")
		os.Exit(1)
	default:
		fmt.Println("\nReceived unexpected signal.")
		os.Exit(1)
	}
}
