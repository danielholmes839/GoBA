package main

import (
	"log"
	"net/http"
	"os"

	game "server/gameplay"
	"server/ws"

	"github.com/gorilla/websocket"
)

func wsHandler(game *game.Game) func(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatalln(err)
		}

		client := ws.NewClient(conn)
		go client.WriteMessages()
		go client.ReceiveMessages(game)
		game.Connect(client)    // connect client to the game
		client.Wait()           // block until the websocket disconnects
		game.Disconnect(client) // disconnect client from the game
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	g := game.NewGame(64)
	go g.Run()

	// Game websocket connection
	http.HandleFunc("/ws", wsHandler(g))
	// Files
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	http.ListenAndServe(":"+port, nil)
}
