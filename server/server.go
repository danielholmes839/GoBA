package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"server/ws"

	"github.com/gorilla/websocket"
)

// Update struct
type Update struct {
	Message string  `json:"message"`
	Points  []Point `json:"points"`
}

// Point struct
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Game struct
type Game struct {
	id      string
	clients int
	events  *ws.Subscription
	//champions map[string]*Champion
}

// Connect func
func (game *Game) Connect() {
	game.clients++
}

// Disconnect func
func (game *Game) Disconnect() {
	game.clients--
}

// Run func
func (game *Game) Run() {
	// start subscriptions
	go game.events.Run()

	counter := 0
	ticker := time.NewTicker(time.Second/5)


	for range ticker.C {

		points := make([]Point, game.clients)
		for i := 0; i < game.clients; i++ {
			points[i] = Point{X: rand.Intn(100), Y: rand.Intn(100)}
		}

		fmt.Println(game.clients)

		message := fmt.Sprintf("tick %d", counter)
		data, err := json.Marshal(&Update{Message: message, Points: points})
		if err != nil {
			game.events.Broadcast([]byte(err.Error()))
		} else {
			game.events.Broadcast(data)
		}

		
		counter++
	}
}

// NewGame func
func NewGame(id string) *Game {
	return &Game{
		id:     id,
		clients: 0,
		events: ws.NewSubscription("events"),
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home")
}

func wsHandler(game *Game) func(w http.ResponseWriter, r *http.Request) {
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
		game.Connect()
		go client.Write()
		go client.Read()
		client.Subscribe(game.events)
		client.Wait() // block until the websocket disconnects
		game.Disconnect()
	}
}

func main() {
	game := NewGame("test-game")
	go game.Run()

	http.HandleFunc("/", home)
	http.HandleFunc("/ws", wsHandler(game))

	log.Println("Listening on localhost:8080")
	http.ListenAndServe(":8080", nil)

}
