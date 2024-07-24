package handler

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

type File struct {
	l *log.Logger
	s *files.Storage
}

func NewFiles(l *log.Logger, s *files.Storage) *File {
	return &File{l, s}
}

func (f *File) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fileName := vars["filename"]
	if id == "" || fileName == "" {
		f.l.Println("No id/filename found")
		http.Error(rw, "No id/filename found", http.StatusBadRequest)
		return
	}
	f.SaveFile(id, fileName, rw, r)
}

func (f *File) SaveFile(id string, filename string, rw http.ResponseWriter, r *http.Request) {
	f.l.Println("Save file for the id", id, "filename", filename)
	fp := filepath.Join(id, filename)
	err := files.Store(fp)
	if err != nil {
		f.l.Println("Unable to save the file")
		http.Error(rw, "Unable to save the file", http.StatusInternalServerError)
	}

}
