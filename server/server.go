package server

import (
	"fmt"
	_log "log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloudflare/tableflip"
)

type Server struct {
	Port     int
	PIDFile  string
	Name     string
	LogColor string
}

type HandlerFunc func(net.Conn)

func (server *Server) Log(msg string, a ...any) {
	_log.Printf("%s[%s]\u001B[0m: %s\n", server.LogColor, server.Name, fmt.Sprintf(msg, a...))
}

func (server *Server) Run(handlerFunc HandlerFunc) {
	upg, _ := tableflip.New(tableflip.Options{
		PIDFile: server.PIDFile,
	})
	defer upg.Stop()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP)
		for range sig {
			upg.Upgrade()
		}
	}()

	ln, err := upg.Listen("tcp", fmt.Sprintf("localhost:%d", server.Port))
	if err != nil {
		os.Exit(1)
	}
	defer ln.Close()

	server.Log("listening for connections on localhost:%d", server.Port)

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				continue
			}
			go handlerFunc(conn)
		}
	}()

	if err := upg.Ready(); err != nil {
		panic(err)
	}

	<-upg.Exit()
}
