package main

import (
	"net/http"
	"os"
	"server/game"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	
	mgr := game.NewGameManager()
	mgr.CreateGame("TEST")

	// Game API
	http.HandleFunc("/join", mgr.GameJoinAPI)
	http.HandleFunc("/create", mgr.GameCreateAPI)
	http.HandleFunc("/info", mgr.InfoAPI)

	// Files
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	http.ListenAndServe(":"+port, nil)
}
