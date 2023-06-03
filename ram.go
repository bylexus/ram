package main

import (
	"net/http"
	"os"
	"os/signal"

	e "github.com/bylexus/go-stdlib/err"
	l "github.com/bylexus/go-stdlib/log"
	"github.com/bylexus/ram/db"
	"github.com/bylexus/ram/server"
	"github.com/jessevdk/go-flags"
)

type ProgramArgs struct {
	FrontendDir string `short:"f" long:"frontend-dir" default:"./public_html"`
	ListenAddr  string `short:"l" long:"listen" default:":3333" description:"listener IP:Port in the form <ip:port>"`
	DbPath      string `short:"d" long:"db" default:"./ram.sqlite" description:"Path to the sqlite db file"`
}

func main() {
	// Create a logger
	logger := l.NewDefaultSeverityLogger()

	opts := ProgramArgs{}
	_, err := flags.Parse(&opts)
	e.PanicOnErr(err)

	// First things first: init the db. This makes sure the schema is
	// created and up to date.
	conn := db.Connect(opts.DbPath)
	defer conn.Close()
	db.InitDb(&logger, conn)

	// Listen for an os interrupt signal
	handleOsInterrupts(&logger)

	// start web server
	startServer(&logger, opts)
}

// Starts the web server: start returns a channel that is closed
// when the web server shuts down.
// We wait so long:
func startServer(logger *l.SeverityLogger, programArgs ProgramArgs) {
	serverConf := server.ServerConfig{
		StaticDir:  programArgs.FrontendDir,
		ListenAddr: programArgs.ListenAddr,
	}

	done := server.Start(logger, serverConf)
	if err := <-done; err != http.ErrServerClosed {
		logger.Error("%s", err)
	} else {
		logger.Info("Server shut down")
	}
}

func handleOsInterrupts(logger *l.SeverityLogger) {
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		server.Shutdown(logger)
	}()
}
