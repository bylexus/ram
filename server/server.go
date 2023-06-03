package server

import (
	"context"
	"net/http"

	l "github.com/bylexus/go-stdlib/log"
)

// Defines the parameters needed for the
// HTTP server
type ServerConfig struct {
	StaticDir  string
	ListenAddr string
}

// Contains the http server infrastructure
type Server struct {
	httpServer *http.Server
	logger     *l.SeverityLogger
	Config     ServerConfig
}

// Needed for the PUT /notes route:
// contains new note data
type NewNoteData struct {
	Note string
	Url  string
	Tags string
}

var server *Server = nil
var serverWait chan error

// Starts the HTTP server.
// Returns a channel (chan error) imediately which is closed when the server
// shuts down (or otherwise encounters a problem).
func Start(logger *l.SeverityLogger, conf ServerConfig) chan error {
	serverWait = make(chan error, 1)

	mux := http.NewServeMux()

	server = &Server{
		httpServer: &http.Server{
			Addr:     conf.ListenAddr,
			Handler:  mux,
			ErrorLog: logger.GetInternalLogger(l.ERROR),
		},
		logger: logger,
		Config: conf,
	}

	// Register route handlers:
	mux.Handle("/", http.FileServer(http.Dir(conf.StaticDir)))
	mux.HandleFunc("/notes", server.handleNotesRoute)

	// Start the web server in a separate goroutine,
	// to de-block the main thread that called Start.
	// After the server dies, we send the error through the returned channel
	// and closes it.
	go func() {
		logger.Info("r.a.m. is starting on %s from %s", server.httpServer.Addr, conf.StaticDir)
		err := server.httpServer.ListenAndServe()
		serverWait <- err
		close(serverWait)
	}()
	return serverWait
}

// Shuts down the HTTP server
func Shutdown(logger *l.SeverityLogger) error {
	var err error = nil
	if server != nil {
		logger.Info("Initiating shutdown")
		err = server.httpServer.Shutdown(context.Background())
	}
	if err == nil {
		server = nil
		logger.Info("Shutdown successful")
	}
	return err
}
