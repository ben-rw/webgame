package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/joho/godotenv"
)

func main() {
	const port = "8080"
	const filepathRoot = "."
	var rootURL string = fmt.Sprintf("http://localhost:%s/", port)

	godotenv.Load()

	cfg := config{
		Port:         port,
		FilepathRoot: filepathRoot,
		RootURL:      rootURL,
		mu:           &sync.RWMutex{},
		activeRooms:  map[string]*Room{},
	}

	mux := http.NewServeMux()

	mux.Handle("/", noCacheMiddleware(http.FileServer(http.Dir(filepathRoot))))

	assetsHandler := http.StripPrefix("/assets/", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/assets/", noCacheMiddleware(assetsHandler))

	mux.HandleFunc("POST /room", cfg.handlerCreateRoom)
	mux.HandleFunc("POST /room/join", cfg.handlerJoinRoom)

	mux.HandleFunc("GET /room/{roomID}", cfg.handlerServeRoomPage)

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on: %s\n", rootURL)
	log.Fatal(server.ListenAndServe())
}
