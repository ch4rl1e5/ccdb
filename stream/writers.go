package stream

import (
	"github.com/google/uuid"
	"net/http"
)

type Writers struct {
	uuid   uuid.UUID
	writer *http.ResponseWriter
	next   *Writers
}

func NewWriter() *Writers {
	return &Writers{uuid: uuid.New()}
}

func (w *Writers) AddWriter(writer *http.ResponseWriter) {
	if w.writer == nil {
		w.writer = writer
		return
	}

	w.AddWriter(w.next.writer)
}
