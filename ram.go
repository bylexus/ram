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

type ServeArgs struct {
	FrontendDir string `short:"f" long:"frontend-dir" default:"./public_html"`
	ListenAddr  string `short:"l" long:"listen" default:":3333" description:"listener IP:Port in the form <ip:port>"`
	DbPath      string `short:"d" long:"db" default:"./ram.sqlite" description:"Path to the sqlite db file"`
}

type ProgramArgs struct {
	Help bool `short:"h" long:"help"`
}

func main() {
	// Create a logger
	logger := l.NewDefaultSeverityLogger()

	programArgs := ProgramArgs{}
	serveArgs := ServeArgs{}

	parser := flags.NewParser(&programArgs, flags.Default)
	parser.AddCommand("serve", "Runs the web server", "Starts the notes app by firing up the web server", &serveArgs)

	_, err := parser.Parse()
	if err != nil {
		if flags.WroteHelp(err) {
			os.Exit(1)
		} else {
			e.PanicOnErr(err)
		}
	}

	// First things first: init the db. This makes sure the schema is
	// created and up to date.
	conn := db.Connect(serveArgs.DbPath)
	defer conn.Close()
	db.InitDb(&logger, conn)

	// Listen for an os interrupt signal
	handleOsInterrupts(&logger)

	// start web server
	startServer(&logger, serveArgs)
}

// Starts the web server: start returns a channel that is closed
// when the web server shuts down.
// We wait so long:
func startServer(logger *l.SeverityLogger, serveArgs ServeArgs) {
	serverConf := server.ServerConfig{
		StaticDir:  serveArgs.FrontendDir,
		ListenAddr: serveArgs.ListenAddr,
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
