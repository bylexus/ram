package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/bylexus/ram/db"
	"github.com/bylexus/ram/model"
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
	logger     *log.Logger
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
func Start(logger *log.Logger, conf ServerConfig) chan error {
	serverWait = make(chan error, 1)

	mux := http.NewServeMux()

	server = &Server{
		httpServer: &http.Server{
			Addr:     conf.ListenAddr,
			Handler:  mux,
			ErrorLog: logger,
		},
		logger: logger,
	}

	// Register route handlers:
	mux.Handle("/", http.FileServer(http.Dir(conf.StaticDir)))
	mux.HandleFunc("/notes", server.handleNotesRoute)

	// Start the web server in a separate goroutine,
	// to de-block the main thread that called Start.
	// After the server dies, we send the error through the returned channel
	// and closes it.
	go func() {
		logger.Printf("r.a.m. is starting on %s from %s\n", server.httpServer.Addr, conf.StaticDir)
		err := server.httpServer.ListenAndServe()
		serverWait <- err
		close(serverWait)
	}()
	return serverWait
}

// Shuts down the HTTP server
func Shutdown(logger *log.Logger) error {
	var err error = nil
	if server != nil {
		logger.Println("Initiating shutdown")
		err = server.httpServer.Shutdown(context.Background())
	}
	if err == nil {
		server = nil
		logger.Println("Shutdown successful")
	}
	return err
}

// Route Handler: /notes
func (s *Server) handleNotesRoute(w http.ResponseWriter, r *http.Request) {
	s.logger.Printf("%s %s\n", r.Method, r.RequestURI)
	switch r.Method {
	case http.MethodPut:
		s.handleNotesPush(w, r)

	}

}

// Route Handler: PUT /notes
// Reads new notes data and creates a persistend Notes entry in the DB
func (s *Server) handleNotesPush(w http.ResponseWriter, r *http.Request) {
	data, err := readJson(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %s", err)
	}
	s.logger.Printf("Got note data: %#v\n", data)
	note := model.NewNote(data.Note, data.Url, data.Tags)
	err = db.PersistNote(r.Context(), &note)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err)
	}
}

func readJson(r io.Reader) (*NewNoteData, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	noteData := NewNoteData{}
	err = json.Unmarshal(data, &noteData)
	if err != nil {
		return nil, err
	}
	return &noteData, nil
}
