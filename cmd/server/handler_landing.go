package main

import (
	"net/http"
)

type LandingPageError struct {
	JoinUsernameError   string
	JoinCodeError       string
	CreateUsernameError string
	JoinRoomFullError   string
}

func (cfg *config) handlerServeLandingPage(w http.ResponseWriter, r *http.Request) {
	err := cfg.templates.ExecuteTemplate(w, "landing.html", LandingPageError{})
	if err != nil {
		http.Error(w, "couldn't load landing page template", http.StatusInternalServerError)
		return
	}
}
