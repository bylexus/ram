package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path"

	"github.com/bylexus/ram/db"
	"github.com/bylexus/ram/model"
)

// Route Handler: /notes
func (s *Server) handleNotesRoute(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug("%s %s", r.Method, r.RequestURI)
	switch r.Method {
	case http.MethodPut:
		s.handleNotesPush(w, r)
	case http.MethodGet:
		s.handleNotesGet(w, r)
	}
}

// Route Handler: PUT /notes
// Reads new notes data and creates a persistend Notes entry in the DB
func (s *Server) handleNotesPush(w http.ResponseWriter, r *http.Request) {
	data, err := readNewNoteJson(r.Body)
	if err != nil {
		NewErrorJsonResponse(nil, http.StatusBadRequest, err, http.StatusBadRequest).WriteHttpResponse(w)
		s.logger.Error("%s", err)
		return
	}
	s.logger.Debug("Got note data: %#v", data)
	note := model.NewNote(data.Note, data.Url, data.Tags)
	err = db.PersistNote(r.Context(), &note)
	if err != nil {
		NewErrorJsonResponse(nil, http.StatusInternalServerError, err, http.StatusInternalServerError).WriteHttpResponse(w)
		s.logger.Error("%s", err)
		return
	}
	NewOkJsonResponse(&note).WriteHttpResponse(w)
}

// Route Handler: GET /notes
// Reads available notes, returns HTML snippets of the note.
func (s *Server) handleNotesGet(w http.ResponseWriter, r *http.Request) {
	notes, err := db.QueryNotes(r.Context())
	if r.Header.Get("Hx-Request") != "" {
		// It's a htmx request
		if err != nil {
			fmt.Fprintf(w, "<div class='error'>ERROR: %s</div>", err.Error())
			s.logger.Error("%s", err)
			return
		}
		var tplFile = "notes-list.html"
		var tplPath = path.Join(s.Config.StaticDir, "templates", "notes-list.html")
		tmpl, err := template.New(tplFile).ParseFiles(tplPath)
		if err != nil {
			fmt.Fprintf(w, "<div class='error'>ERROR: %s</div>", err.Error())
			s.logger.Error("%s", err)
			return
		}
		err = tmpl.Execute(w, notes)
		if err != nil {
			fmt.Fprintf(w, "<div class='error'>ERROR: %s</div>", err.Error())
			s.logger.Error("%s", err)
			return
		}
	} else {
		// It's another request
		if err != nil {
			NewErrorJsonResponse(nil, http.StatusInternalServerError, err, 1).WriteHttpResponse(w)
			s.logger.Error("%s", err)
			return
		}
		NewOkJsonResponse(notes).WriteHttpResponse(w)
	}
}

func readNewNoteJson(r io.Reader) (*NewNoteData, error) {
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
