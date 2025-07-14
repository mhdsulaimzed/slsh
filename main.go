package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"slsh/shell"
)

func main() {
	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Create and configure shell
	sh := shell.New()
	
	// Handle Ctrl+C gracefully
	go func() {
		<-c
		fmt.Println("\nGoodbye!")
		os.Exit(0)
	}()

	// Start the shell
	if err := sh.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}