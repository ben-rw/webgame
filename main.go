package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

//go:embed app/landing
var landingContent embed.FS

func main() {
	godotenv.Load()

	cfg := config{
		Port:         os.Getenv("PORT"),
		FilepathRoot: os.Getenv("FILEPATH_ROOT"),
		URLRoot:      os.Getenv("URL_ROOT"),
		mu:           &sync.RWMutex{},
		activeRooms:  map[string]*Room{},
	}

	landingFS, err := fs.Sub(landingContent, "app/landing")
	if err != nil {
		log.Fatal(err)
	}
	landingServer := http.FileServer(http.FS(landingFS))

	mux := http.NewServeMux()

	mux.Handle("GET /", noCacheMiddleware(landingServer))

	mux.HandleFunc("POST /room", cfg.handlerCreateRoom)
	mux.HandleFunc("POST /room/join", cfg.handlerJoinRoom)

	mux.HandleFunc("GET /room/{roomID}", cfg.handlerServeRoomPage)

	server := http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	log.Printf("Serving on: %s%s/\n", cfg.URLRoot, cfg.Port)
	log.Fatal(server.ListenAndServe())
}
