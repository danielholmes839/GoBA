package gameplay

import (
	"errors"
	"fmt"
	"server/game/gameplay/geometry"
	"server/ws"
	"time"
)

type IGame interface {
	Connect(*ws.Client) error
	Disconnect(*ws.Client) error
	Run()
	Stop() error
	UsernameTaken(name string) bool
	GetPlayerCount() int
}

// ClientInfo struct
type ClientInfo struct {
	team     *Team
	champion *Champion
	score    *Score
}

// Game struct
type Game struct {
	// Game settings
	tps int // "Ticks per second"

	// Channels
	connect     chan *ws.Client
	disconnect  chan *ws.Client
	stop        chan bool
	playerCount chan int

	events *ws.ClientEventQueue // Events received are put in this queue
	global *ws.Subscription     // Events relevant to all clients are sent on this subscription

	// Game variables
	usernames   map[string]bool
	clients     map[*ws.Client]*ClientInfo
	teams       map[*Team]struct{}
	projectiles map[*Projectile]*ws.Client

	// Structures
	walls  []*geometry.Rectangle
	bushes []*geometry.Rectangle
}

// NewGame func
func NewGame(tps int) *Game {
	g := &Game{
		tps:         tps,
		connect:     make(chan *ws.Client),
		disconnect:  make(chan *ws.Client),
		stop:        make(chan bool),
		playerCount: make(chan int),
		events:      ws.NewClientEventQueue(),
		global:      ws.NewSubscription("global-events"),
		teams:       make(map[*Team]struct{}),
		clients:     make(map[*ws.Client]*ClientInfo),
		projectiles: make(map[*Projectile]*ws.Client),
		usernames:   make(map[string]bool),

		// Structures
		walls: []*geometry.Rectangle{
			geometry.NewRectangle(-2000, -2000, 7000, 2000),
			geometry.NewRectangle(-2000, -2000, 2000, 7000),
			geometry.NewRectangle(-2000, 3000, 7000, 2000),
			geometry.NewRectangle(3000, -2000, 2000, 7000),

			geometry.NewRectangle(400, 400, 950, 100),
			geometry.NewRectangle(1650, 400, 950, 100),

			geometry.NewRectangle(400, 2500, 950, 100),
			geometry.NewRectangle(1650, 2500, 950, 100),
		},

		bushes: []*geometry.Rectangle{
			geometry.NewRectangle(500, 500, 500, 300), geometry.NewRectangle(1250, 500, 500, 300), geometry.NewRectangle(2000, 500, 500, 300),
			geometry.NewRectangle(500, 2200, 500, 300), geometry.NewRectangle(1250, 2200, 500, 300), geometry.NewRectangle(2000, 2200, 500, 300),
			geometry.NewRectangle(0, 1000, 300, 1000), geometry.NewRectangle(2700, 1000, 300, 1000),
		},
	}

	g.teams[NewTeam("Red Team", "#ff0000", geometry.NewPoint(1500, 200))] = struct{}{}
	g.teams[NewTeam("Blue Team", "#0000ff", geometry.NewPoint(1500, 2800))] = struct{}{}

	return g
}

// Handle func
func (game *Game) Handle(event *ws.ClientEvent) {
	fmt.Printf("category:'%s' name:'%s' client:'%s'\n", event.GetCategory(), event.GetName(), event.GetClient().GetID())
	switch event.Category {
	case "game":
		game.events.Push(event)
	}
}

// Run func
func (game *Game) Run() {
	// switch between processing game ticks, events, connecting/disconnecting clients from the game
	tick := time.NewTicker(time.Duration(int(time.Second) / game.tps))

	defer func() {
		close(game.connect)
		close(game.disconnect)
		close(game.playerCount)
		close(game.stop)

		for client := range game.clients {
			game.disconnectClient(client)
		}

		for team := range game.teams {
			team.events.Stop()
		}

		game.global.Stop()

	}()

	for {
		select {
		case client := <-game.connect:
			// Connect clients
			game.connectClient(client)

		case client := <-game.disconnect:
			// Disconnect clients
			game.disconnectClient(client)

		case <-game.stop:
			// Stop the game
			return

		case <-tick.C:
			// Game Loop
			game.tick()

		case game.playerCount <- len(game.usernames):
		}
	}
}

// Connect clients to the game
func (game *Game) Connect(client *ws.Client) error {
	select {
	case game.connect <- client:
		return nil
	default:
		return errors.New("This game has ended")
	}
}

// Disconnect clients from the game
func (game *Game) Disconnect(client *ws.Client) error {
	select {
	case game.disconnect <- client:
		return nil
	default:
		return errors.New("This game has ended")
	}
}

// Stop the game - returns false if the game was already stopped
func (game *Game) Stop() error {
	select {
	case game.stop <- true:
		return nil
	default:
		return errors.New("This game has ended")
	}
}

// UsernameTaken func
func (game *Game) UsernameTaken(name string) bool {
	return game.usernames[name]
}

// GetPlayerCount func
func (game *Game) GetPlayerCount() int {
	return <-game.playerCount
}

// Get the team with lowest number of players
func (game *Game) getNextTeam() *Team {
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

func (game *Game) getClientChampion(client *ws.Client) *Champion {
	// Get the champion of the client
	return game.clients[client].champion
}

func (game *Game) getClientTeam(client *ws.Client) *Team {
	// Get the team of the client
	return game.clients[client].team
}

func (game *Game) getClientScore(client *ws.Client) *Score {
	// Get the info (team and champion) of the client
	return game.clients[client].score
}

func (game *Game) getClientInfo(client *ws.Client) *ClientInfo {
	// Get the info (team and champion) of the client
	return game.clients[client]
}

func (game *Game) createClientInfo(client *ws.Client, champion *Champion, team *Team) {
	// Set the info (team and champion) of the client
	score := NewScore(0, 0, 0)
	game.clients[client] = &ClientInfo{champion: champion, team: team, score: score}
}

func (game *Game) connectClient(client *ws.Client) {
	// Add client to the game
	game.usernames[client.GetUsername()] = true
	client.Subscribe(game.global)
	client.WriteMessage("setup", NewSetupUpdate(game, client))

	champion := NewChampion(client.GetID())
	team := game.getNextTeam()
	team.addClient(client, champion)

	game.createClientInfo(client, champion, team)

	// Send clients the updated teams
	game.global.Broadcast("update-teams", NewTeamsUpdate(game))
}

func (game *Game) disconnectClient(client *ws.Client) {
	// Disconnect the client
	client.Close()
	game.getClientTeam(client).removeClient(client) // Remove client from team

	delete(game.usernames, client.GetUsername()) // Remove client username from the game
	delete(game.clients, client)             // Remove client from game

	// Send clients the updated teams
	game.global.Broadcast("update-teams", NewTeamsUpdate(game))
}
