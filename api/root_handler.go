package api

import "net/http"

type RootHandler struct{}

func NewRootHandler() *RootHandler {
	return &RootHandler{}
}

func (h *RootHandler) HandleGetRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/contacts", http.StatusMovedPermanently)
}
