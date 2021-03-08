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

func main() {
	port := config("PORT", "5000")
	room := config("ROOM", "TEST")

	// Create the game manager
	mgr := game.NewGameManager()
	mgr.CreateCustomGame(room, nil)

	http.HandleFunc("/join", mgr.GameJoinEndpoint)
	http.HandleFunc("/create", mgr.GameCreateEndpoint)
	http.HandleFunc("/info", mgr.InfoEndpoint)

	fs := http.FileServer(http.Dir("./client"))
	http.Handle("/client/", http.StripPrefix("/client/", fs))

	// Serve the index.html at /game
	http.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./client/index.html")
	})

	http.ListenAndServe(":"+port, nil)

}
