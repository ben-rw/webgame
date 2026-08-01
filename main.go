package main

import (
	// "io/fs"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/ben-rw/webgame/frontend"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	templates, err := frontend.LoadTemplates()
	if err != nil {
		log.Fatal(err)
	}

	cfg := config{
		Port:         os.Getenv("PORT"),
		FilepathRoot: os.Getenv("FILEPATH_ROOT"),
		URLRoot:      os.Getenv("URL_ROOT"),
		mu:           &sync.RWMutex{},
		activeRooms:  map[string]*Room{},
		templates:    templates,
	}

	// landingFS, err := fs.Sub(LandingContent, "frontend/landing")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// landingServer := http.FileServer(http.FS(landingFS))

	mux := http.NewServeMux()

	// mux.Handle("GET /", noCacheMiddleware(landingServer))
	landingHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := cfg.templates.ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			log.Fatal(err)
		}
	})
	mux.Handle("GET /", noCacheMiddleware(landingHandler))

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
