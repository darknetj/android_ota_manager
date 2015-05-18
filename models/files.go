package models
import (
  "net/http"
)

type File struct {
    Name string
    Size int64
}

func FilesHandler(w http.ResponseWriter, r *http.Request) {}
