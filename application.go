package main

import (
	"net/http"
	"os"
	"server/game"
)

func config(key string, defaultValue string) string {
	value := os.Getenv("PORT")
	if value == "" {
		return defaultValue
	}
	return value
}

func fileserver() {
	fs := http.FileServer(http.Dir("./client"))
	http.Handle("/client/", http.StripPrefix("/client/", fs))
	http.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./client/index.html")
	})
}

func main() {
	port := config("PORT", "5000")
	room := config("ROOM", "TEST")

	mgr := game.NewGameManager()
	mgr.SetupEndpoints()
	mgr.CreateGameManually(room, nil)
	
	fileserver()

	// Listen
	http.ListenAndServe(":"+port, nil)
}
