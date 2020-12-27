package game

import (
	"encoding/json"
	"fmt"
	"server/game/geometry"
	"server/ws"
	"time"
)

// Game struct
type Game struct {
	id     string
	tps    int             // ticks per second target
	add    chan *ws.Client // channel to add clients to the game
	remove chan *ws.Client // channel to remove clients from the game

	// sending, and receiving events
	events *ws.EventQueue
	global *ws.Subscription

	// clients, teams
	walls   []*Wall
	teams   map[*Team]bool
	clients map[*ws.Client]*Team // Get the team of client
}

// NewGame func
func NewGame(id string, tps int) *Game {
	g := &Game{
		id:     id,
		tps:    tps,
		add:    make(chan *ws.Client),
		remove: make(chan *ws.Client),
		events: ws.NewEventQueue(),
		global: ws.NewSubscription("global-events"),

		walls:   []*Wall{NewWall(500, 500, 500, 100), NewWall(100, 100, 500, 100)},
		teams:   make(map[*Team]bool),
		clients: make(map[*ws.Client]*Team),
	}

	g.teams[NewTeam()] = true
	g.teams[NewTeam()] = true

	return g
}

// Handle func
func (game *Game) Handle(event *ws.ClientEvent) {
	fmt.Printf("category:'%s' name:'%s' client:'%s'\n", event.GetCategory(), event.GetName(), event.GetClient().GetID())
	switch event.Category {
	// Add game event to the queue
	case "game":
		game.events.Push(event)
	}
}

// Run func
func (game *Game) Run() {
	// start subscriptions
	go game.global.Run()
	for team := range game.teams {
		go team.events.Run()
	}

	// switch between processing game ticks, events, connecting/disconnecting clients from the game
	ticker := time.NewTicker(time.Duration(int(time.Second) / game.tps))
	
	for {
		select {
		case client := <-game.add:
			// Connect clients
			team := game.GetSmallestTeam()
			team.AddClient(client, NewChampion(client.GetID()))
			client.Subscribe(game.global)
			game.clients[client] = team

			walls := []*WallJSON{}
			for _, wall := range game.walls {
				walls = append(walls, NewWallJSON(wall))
			}
			data, _ := json.Marshal(walls)
			client.WriteMessage("walls", data)

		case client := <-game.remove:
			// Disconnect clients
			team := game.GetTeam(client)
			team.RemoveClient(client)
			client.Unsubscribe(game.global)
			delete(game.clients, client)

		case <-ticker.C:
			// Game Loop
			game.Tick()
		}
	}
}

// Tick func
func (game *Game) Tick() {
	// Empty the event queue
	for event := range game.events.Read() {
		champ := game.GetChampion(event.Client)
		switch event.Name {
		case "move-event":
			champ.MoveEvent(event)
			champ.health -= 5
		default:
			fmt.Printf("unknown game event: %s", event.Name)
		}
	}
	
	// Calculate vision of other objects. broadcast to the team
	for team := range game.teams {
		team.Update(game)
	}
	
	//data, _ := json.Marshal(game.NewTickUpdate())
	// game.global.Broadcast("update", data)
}

// LineOfSight func
func (game *Game) LineOfSight(line *geometry.Line) bool {
	for _, wall := range game.walls {
		if (line.HitsRectangle(wall.hitbox)) {
			return false
		}
	}
	return true
}


// GetChampion of client
func (game *Game) GetChampion(client *ws.Client) *Champion {
	return game.GetTeam(client).GetChampion(client)
}

// GetTeam of client
func (game *Game) GetTeam(client *ws.Client) *Team {
	return game.clients[client]
}

// GetSmallestTeam the team with lowest number of players
func (game *Game) GetSmallestTeam() *Team {
	var min *Team = nil
	for team := range game.teams {
		if min == nil {
			min = team
		} else if team.size < min.size {
			min = team
		}
	}
	return min
}

// Connect func
func (game *Game) Connect(client *ws.Client) {
	game.add <- client
}

// Disconnect func
func (game *Game) Disconnect(client *ws.Client) {
	game.remove <- client
}
