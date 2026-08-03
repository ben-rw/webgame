package main

import (
	"html/template"
	"sync"
)

type config struct {
	Port         string
	FilepathRoot string
	URLRoot      string
	mu           *sync.RWMutex
	activeRooms  map[string]*Room
	templates    *template.Template
}
