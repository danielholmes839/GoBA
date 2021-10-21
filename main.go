package main

import (
	"net/http"
	"os"
	"server/game"
)

// Get config variables
func config(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func serveFiles() {
	// Serve files in ./client
	fs := http.FileServer(http.Dir("./client"))
	http.Handle("/client/", http.StripPrefix("/client/", fs))

	// Serve the index.html at /game
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./client/index.html")
	})
}

func main() {
	port := config("PORT", "5000")
	room := config("ROOM", "TEST")

	mgr := game.NewGameManager()
	mgr.SetupEndpoints()
	mgr.CreateGameManually(room, nil)

	serveFiles()

	// Listen
	http.ListenAndServe(":"+port, nil)
}
