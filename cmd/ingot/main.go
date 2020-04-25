package main

import (
	"github.com/ingotmc/ingot/net"
	"log"
	"os"
	"os/signal"
)

func main() {
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
