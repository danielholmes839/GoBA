package main

import (
	"net/http"
	"os"
	"server/game"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":5000"
	}

	mgr := game.NewGameManager()
	mgr.CreateGame("TEST", true)

	// Game API
	http.HandleFunc("/join", mgr.GameJoinAPI)
	http.HandleFunc("/create", mgr.GameCreateAPI)
	http.HandleFunc("/info", mgr.InfoAPI)

	// Files
	fs := http.FileServer(http.Dir("./client"))
	http.Handle("/client/", http.StripPrefix("/client/", fs))
	http.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./client/index.html")
	})

	// Listen
	http.ListenAndServe(port, nil)
}
