package gameplay

import (
	"fmt"
	"server/game/gameplay/geometry"
	"server/ws"
	"time"
)

// ClientInfo struct
type ClientInfo struct {
	team     *Team
	champion *Champion
}

// Game struct
type Game struct {
	// Game settings
	tps       int
	permanent bool
	stopped   bool

	// Add, remove clients
	connect     chan *ws.Client
	disconnect  chan *ws.Client
	stop        chan bool
	playerCount chan int

	events *ws.ClientEventQueue // Events received are put in this queue
	global *ws.Subscription     // Events relevant to all clients are sent on this subscription

	// Game variables
	walls     []*Wall
	bushes    []*Bush
	teams     map[*Team]struct{}
	clients   map[*ws.Client]*ClientInfo
	usernames map[string]bool
}

// NewGame func
func NewGame(tps int, permanent bool) *Game {
	g := &Game{
		tps:         tps,
		permanent:   permanent,
		stopped:     false,
		connect:     make(chan *ws.Client),
		disconnect:  make(chan *ws.Client),
		stop:        make(chan bool),
		playerCount: make(chan int),
		events:      ws.NewClientEventQueue(),
		global:      ws.NewSubscription("global-events"),
		teams:       make(map[*Team]struct{}),
		clients:     make(map[*ws.Client]*ClientInfo),
		usernames:   make(map[string]bool),

		// Structures
		walls: []*Wall{
			NewWall(-2000, -2000, 7000, 2000),
			NewWall(-2000, -2000, 2000, 7000),
			NewWall(-2000, 3000, 7000, 2000),
			NewWall(3000, -2000, 2000, 7000),

			NewWall(400, 400, 950, 100),
			NewWall(1650, 400, 950, 100),

			NewWall(400, 2500, 950, 100),
			NewWall(1650, 2500, 950, 100),
		},

		bushes: []*Bush{
			NewBush(500, 500, 500, 300), NewBush(1250, 500, 500, 300), NewBush(2000, 500, 500, 300),
			NewBush(500, 2200, 500, 300), NewBush(1250, 2200, 500, 300), NewBush(2000, 2200, 500, 300),
			NewBush(0, 1000, 300, 1000), NewBush(2700, 1000, 300, 1000),
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
	defer func() {
		for client := range game.clients {
			game.disconnectClient(client)
		}

		for team := range game.teams {
			team.events.Close()
		}

		game.global.Close()
	}()

	// switch between processing game ticks, events, connecting/disconnecting clients from the game
	tick := time.NewTicker(time.Duration(int(time.Second) / game.tps))

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
			if !game.permanent {
				return
			}

		case <-tick.C:
			// Game Loop
			if len(game.usernames) > 0 {
				game.tick()
			}

		case game.playerCount <- len(game.usernames):
		}
	}
}

// Connect clients to the game
func (game *Game) Connect(client *ws.Client) {
	game.connect <- client
}

// Disconnect clients from the game
func (game *Game) Disconnect(client *ws.Client) {
	game.disconnect <- client
}

// Stop the game
func (game *Game) Stop() bool {
	if game.stopped {
		return true
	}

	if game.permanent {
		return false
	}

	game.stopped = true
	game.stop <- true
	return true
}

func (game *Game) connectClient(client *ws.Client) {
	// Add client to the game
	game.usernames[client.GetName()] = true
	client.Subscribe(game.global)
	client.WriteMessage("setup", NewSetupUpdate(game, client))

	champion := NewChampion(client.GetID())
	team := game.getNextTeam()
	team.addClient(client, champion)

	game.setClientInfo(client, champion, team)

	// Send clients the updated teams
	game.global.Broadcast("update-teams", NewTeamsUpdate(game))
}

func (game *Game) disconnectClient(client *ws.Client) {
	// Disconnect the client
	game.getClientTeam(client).removeClient(client) // Remove client from team
	delete(game.usernames, client.GetName())        // Remove client username from the game
	delete(game.clients, client)                    // Remove client from game
	client.Close()

	// Send clients the updated teams
	game.global.Broadcast("update-teams", NewTeamsUpdate(game))
}

// ############### HELPERS ##############

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

func (game *Game) getClientInfo(client *ws.Client) *ClientInfo {
	// Get the info (team and champion) of the client
	return game.clients[client]
}

func (game *Game) setClientChampion(client *ws.Client, champion *Champion) {
	// Set the champion of the client
	game.clients[client].champion = champion
}

func (game *Game) setClientTeam(client *ws.Client, team *Team) {
	// Set the team of the client
	game.clients[client].team = team
}

func (game *Game) setClientInfo(client *ws.Client, champion *Champion, team *Team) {
	// Set the info (team and champion) of the client
	game.clients[client] = &ClientInfo{champion: champion, team: team}
}
