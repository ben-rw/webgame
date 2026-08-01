package main

import (
	"net/http"
)

func (cfg *config) handlerServeLandingPage(w http.ResponseWriter, r *http.Request) {
	err := cfg.templates.ExecuteTemplate(w, "landing.html", nil)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't load landing page template", err)
		return
	}
}
