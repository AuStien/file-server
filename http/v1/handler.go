package v1

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/austien/file-server/files"
)

func NewHandler(client *files.Client) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/random", handleRandom(client))

	r.Use(mux.CORSMethodMiddleware(r))

	return r
}

func handleRandom(client *files.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		f, err := client.RandomFile()
		if err != nil {
			log.WithError(err).Error("randomFile")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(f)
		if err != nil {
			log.WithError(err).Error("marshal")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(b)
	}
}
