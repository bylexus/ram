package model

import "time"

type Note struct {
	Id      int64
	Note    string
	Url     string
	Tags    []string
	Created time.Time
	Done    bool
}
