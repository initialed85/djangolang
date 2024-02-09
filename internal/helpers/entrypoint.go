package helpers

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitForExit() {
	c := make(chan os.Signal, 1024)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-c
}
