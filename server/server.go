package main

import (
	"fmt"
	"log"
	"net/http"

	"server/game"
	"server/ws"

	"github.com/gorilla/websocket"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home")
}

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
		go client.ReadMessages(game)
		game.Connect(client)		// connect client to the game
		client.Wait() 				// block until the websocket disconnects
		game.Disconnect(client)		// disconnect client from the game
	}
}

func main() {
	g := game.NewGame("test-game", 30)
	go g.Run()

	http.HandleFunc("/", home)
	http.HandleFunc("/ws", wsHandler(g))

	log.Println("Listening on localhost:8080")
	http.ListenAndServe(":8080", nil)

}
