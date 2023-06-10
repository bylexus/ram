package server

import (
	"context"
	"net/http"
	"path"

	l "github.com/bylexus/go-stdlib/log"
)

// Defines the parameters needed for the
// HTTP server
type ServerConfig struct {
	StaticDir  string
	ListenAddr string
}

type ServerContextKey string

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
	// configure global middlewares:
	var sessionMiddleware = NewSessionMiddleware(logger)
	var authMiddleware = NewAuthMiddleware(logger)

	// all requests use a cookie session:
	globalHandler := sessionMiddleware.WrapHandler(
		// limit body data sent to 4k
		http.MaxBytesHandler(mux, 4*1024),
	)

	server = &Server{
		httpServer: &http.Server{
			Addr:     conf.ListenAddr,
			Handler:  globalHandler,
			ErrorLog: logger.GetInternalLogger(l.ERROR),
		},
		logger: logger,
		Config: conf,
	}

	// Register route handlers, wrap in middlewares where apropriate:
	mux.Handle("/notes", authMiddleware.WrapHandler(NewNotesRouter(logger, server)))
	mux.Handle("/guest/", http.FileServer(http.Dir(conf.StaticDir)))
	mux.Handle("/resources/", http.FileServer(http.Dir(conf.StaticDir)))

	// make sure that all non-auth routes are configured first (see above)
	mux.Handle("/", authMiddleware.WrapHandler(http.FileServer(http.Dir(path.Join(conf.StaticDir, "auth")))))

	// Start the web server in a separate goroutine,
	// to de-block the main thread that called Start.
	// After the server dies, we send the error through the returned channel
	// and closes it.
	go func() {
		logger.Info("r.a.m. is starting on %s from %s", server.httpServer.Addr, conf.StaticDir)
		logger.Debug("Server config is: %#v", server.Config)
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
