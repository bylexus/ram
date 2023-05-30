package model

import (
	"strings"
	"time"

	stdstr "github.com/bylexus/go-stdlib/strings"
	"github.com/google/uuid"
)

type Note struct {
	Id      string
	Note    string
	Url     string
	Tags    []string
	Created time.Time
	Done    bool

	// new entries are phantom entries (not yet persisted)
	phantom bool
}

func (n *Note) IsPhantom() bool {
	return n.phantom
}

func (n *Note) SetPhantom(isPhantom bool) {
	n.phantom = isPhantom
}

func NewNote(note string, url string, tags string) Note {
	newNote := Note{
		Id:      uuid.NewString(),
		Note:    note,
		Url:     url,
		Created: time.Now(),
		Done:    false,
		Tags:    make([]string, 0),
		phantom: true,
	}
	tagSlice, err := stdstr.SplitRe(strings.TrimSpace(tags), `[,;\s]+`)
	if err == nil {
		newNote.Tags = tagSlice
	}

	return newNote
}
