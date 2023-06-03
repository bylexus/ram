package db

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bylexus/ram/model"
)

func PersistNote(ctx context.Context, note *model.Note) error {
	conn := GetConn()
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
			var tagsStr = string(tagsByte[:])
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

func QueryNotes(ctx context.Context) ([]*model.Note, error) {
	conn := GetConn()
	stm, err := conn.PrepareContext(ctx, "SELECT id,note,url,created,done,tags FROM note ORDER BY created DESC LIMIT 250")
	if err != nil {
		return nil, err
	}
	defer stm.Close()
	rows, err := stm.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	notes := make([]*model.Note, 0)
	for rows.Next() {
		noteObj := model.Note{}
		var id string
		var note string
		var url string
		var created time.Time
		var done bool
		var tags []byte

		err := rows.Scan(&id, &note, &url, &created, &done, &tags)
		if err != nil {
			return nil, err
		}
		noteObj.SetPhantom(false)
		noteObj.Id = id
		noteObj.Note = note
		noteObj.Url = url
		noteObj.Created = created
		noteObj.Done = done
		var tagsArr []string = make([]string, 0)
		err = json.Unmarshal(tags, &tagsArr)
		if err != nil {
			noteObj.Tags = make([]string, 0)
		} else {
			noteObj.Tags = tagsArr
		}
		notes = append(notes, &noteObj)
	}
	return notes, nil
}
