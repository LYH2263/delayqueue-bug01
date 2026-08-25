package console

import (
	"encoding/json"
	"net/http"

	"github.com/LYH2263/go-delayqueue"
)

type API struct {
	Broker *delayqueue.Broker
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/api/stats":
		writeJSON(w, a.Broker.Stats())
	case "/api/pending":
		writeJSON(w, a.Broker.ListPendingViews())
	case "/api/dead":
		writeJSON(w, a.Broker.SnapshotDead())
	default:
		http.NotFound(w, r)
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
