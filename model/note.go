package model

import (
	"strings"
	"time"

	"github.com/bylexus/go-stdlib/slices"
	stdstr "github.com/bylexus/go-stdlib/strings"
	"github.com/google/uuid"
)

type Note struct {
	Id      string    `json:"id"`
	Note    string    `json:"note"`
	Url     string    `json:"url"`
	Tags    []string  `json:"tags"`
	Created time.Time `json:"created"`
	Done    bool      `json:"done"`

	// new entries are phantom entries (not yet persisted)
	phantom bool
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
		newNote.Tags = slices.Filter(&tagSlice, func(el *string) bool { return *el != "" })
	}

	return newNote
}

func (n *Note) IsPhantom() bool {
	return n.phantom
}

func (n *Note) SetPhantom(isPhantom bool) {
	n.phantom = isPhantom
}
