package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bylexus/ram/model"
)

func PersistNote(ctx context.Context, note *model.Note) error {
	conn := Conn()
	if note.IsPhantom() {
		stm, err := conn.PrepareContext(ctx, `
		INSERT INTO note (id, note, url, tags, created, done)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
		if err != nil {
			return err
		}

		var tagsJson *string = nil
		tagsByte, err := json.Marshal(note.Tags)
		if err == nil {
			var tagsStr = fmt.Sprintf("%s", tagsByte)
			tagsJson = &tagsStr
		}

		_, err = stm.ExecContext(ctx,
			note.Id,
			note.Note,
			note.Url,
			tagsJson,
			note.Created,
			note.Done,
		)
		if err != nil {
			return err
		}
		return stm.Close()
	}

	return nil
}
