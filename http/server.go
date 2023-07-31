package http

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/austien/file-server/files"
	v1 "github.com/austien/file-server/http/v1"
)

func NewServer(addr string, client *files.Client, rootFolder string) *http.Server {
	r := mux.NewRouter()

	r.PathPrefix("/api").Handler(http.StripPrefix("/api", v1.NewHandler(client)))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(rootFolder)))

	return &http.Server{Addr: addr, Handler: r}
}
