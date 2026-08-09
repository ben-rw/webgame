package main

import (
	"github.com/ben-rw/webgame/internal/room"
	"html/template"
)

type config struct {
	Port         string
	FilepathRoot string
	URLRoot      string
	RoomReg      room.RoomRegistry
	templates    *template.Template
}
