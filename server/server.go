package server

import (
	"context"
	"log"
	"net/http"
)

type ServerConfig struct {
	StaticDir  string
	ListenAddr string
}

var server *http.Server = nil
var serverWait chan error

func Start(logger *log.Logger, conf ServerConfig) chan error {
	serverWait = make(chan error, 1)

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(conf.StaticDir)))

	server = &http.Server{
		Addr:     conf.ListenAddr,
		Handler:  mux,
		ErrorLog: logger,
	}

	go func() {
		logger.Printf("r.a.m. is starting on %s from %s\n", server.Addr, conf.StaticDir)
		err := server.ListenAndServe()
		serverWait <- err
		close(serverWait)
	}()
	return serverWait
}

func Shutdown(logger *log.Logger) error {
	var err error = nil
	if server != nil {
		logger.Println("Initiating shutdown")
		err = server.Shutdown(context.Background())
	}
	if err == nil {
		server = nil
		logger.Println("Shutdown successful")
	}
	return err
}
