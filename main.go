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

	mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	assetsHandler := http.StripPrefix("/assets/", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/assets/", assetsHandler)

	mux.HandleFunc("POST /room", cfg.handlerCreateRoom)
	mux.HandleFunc("POST /room/join", cfg.handlerJoinRoom)

	roomHandler := http.FileServer(http.Dir(fmt.Sprintf("%v/app/room/", filepathRoot)))
	mux.Handle("GET /room/{roomID}", roomHandler)

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on: %s\n", rootURL)
	log.Fatal(server.ListenAndServe())
}
