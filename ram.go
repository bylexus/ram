package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/bylexus/go-stdlib"
	"github.com/bylexus/ram/db"
	"github.com/bylexus/ram/server"
	"github.com/jessevdk/go-flags"
)

type ProgramArgs struct {
	FrontendDir string `short:"f" long:"frontend-dir" default:"./public_html"`
	ListenAddr  string `short:"l" long:"listen" default:":3333"`
}

func main() {
	// Create a logger
	logger := log.Default()

	opts := ProgramArgs{}
	_, err := flags.Parse(&opts)
	stdlib.PanicOnErr(err)

	// First things first: init the db. This makes sure the schema is
	// created and up to date.
	conn := db.Conn()
	defer conn.Close()
	db.InitDb(logger, conn)

	// Listen for an os interrupt signal
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, os.Kill)
		<-sigint

		server.Shutdown(logger)
	}()

	// Now, start the web server: start returns a channel that is closed
	// when the web server shuts down.
	// We wait so long:
	serverConf := server.ServerConfig{
		StaticDir:  opts.FrontendDir,
		ListenAddr: opts.ListenAddr,
	}

	done := server.Start(logger, serverConf)
	if err = <-done; err != http.ErrServerClosed {
		logger.Printf("Error: %s\n", err)
	} else {
		logger.Println("Server shut down")
	}
}
