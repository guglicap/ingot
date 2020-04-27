package main

import (
	"github.com/ingotmc/ingot/net"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
)

func main() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	s, err := net.NewServer()
	if err != nil {
		log.Fatal(err)
	}
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	go func() {
		<-sigc
		s.Shutdown()
	}()
	err = s.Listen()
	if err != nil {
		log.Fatal(err)
	}
}
