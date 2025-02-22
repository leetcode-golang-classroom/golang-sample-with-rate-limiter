package greeting

import (
	"encoding/json"
	"net/http"

	"github.com/leetcode-golang-classroom/golang-sample-with-rate-limiter/internal/types"
)

type Handler struct{}

// SetupRoute - 設定 route
func (h *Handler) SetupRoute(router *http.ServeMux) {
	router.HandleFunc("/", h.Greeting)
}

func NewRoute() *Handler {
	return &Handler{}
}

func (h *Handler) Greeting(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := types.Resposne{Message: "Hello, World!"}
	json.NewEncoder(w).Encode(response)
}
