package game

import (
	"fmt"
	"server/game/geometry"
	"server/ws"
	"time"
)

// Player struct
type Player struct {
	champion *Champion
	team     *Team
}

// NewPlayer fucn
func NewPlayer(client *ws.Client, team *Team) *Player {
	return &Player{
		champion: NewChampion(client.GetID()),
		team:     team,
	}
}

// Game struct
type Game struct {
	// Game settings
	id  string
	tps int

	// Add, remove clients
	add    chan *ws.Client
	remove chan *ws.Client

	events *ws.ClientEventQueue // Events received are put in this queue
	global *ws.Subscription     // Events relevant to all clients are sent on this subscription

	// Game variables
	walls  []*Wall
	bushes []*Bush
	teams  map[*Team]bool

	clientChampions map[*ws.Client]*Champion
	clientTeams     map[*ws.Client]*Team
}

// NewGame func
func NewGame(id string, tps int) *Game {
	g := &Game{
		id:  id,
		tps: tps,

		add:    make(chan *ws.Client),
		remove: make(chan *ws.Client),

		events: ws.NewClientEventQueue(),
		global: ws.NewSubscription("global-events"),

		walls:  []*Wall{NewWall(0, 1500, 2000, 200), NewWall(700, 600, 500, 100), NewWall(100, 100, 500, 100)},
		bushes: []*Bush{NewBush(700, 800, 500, 400), NewBush(100, 300, 500, 400)},

		teams:           make(map[*Team]bool),
		clientChampions: make(map[*ws.Client]*Champion),
		clientTeams:     make(map[*ws.Client]*Team),
	}

	g.teams[NewTeam("0", "#ff0000")] = true
	g.teams[NewTeam("1", "#0000ff")] = true

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
			game.addClient(client)

		case client := <-game.remove:
			// Disconnect clients
			game.removeClient(client)

		case <-ticker.C:
			// Game Loop
			game.tick()
		}
	}
}

// Tick func
func (game *Game) tick() {
	// Empty the event queue
	for event := range game.events.Read() {
		champ := game.getClientChampion(event.Client)
		switch event.Name {
		case "move":
			champ.setMovementDirection(event)
		default:
			fmt.Printf("unknown game event name: '%s'\n", event.Name)
		}
	}

	for _, champ := range game.clientChampions {
		champ.move(game)
	}

	// Team tick
	for team := range game.teams {
		team.tick(game)
	}
}

func (game *Game) hasLineOfSight(line *geometry.Line) bool {
	for _, wall := range game.walls {
		if line.HitsRectangle(wall.hitbox) {
			return false
		}
	}

	for _, bush := range game.bushes {
		sourceInBush := line.GetStart().HitsRectangle(bush.hitbox)
		targetInBush := line.GetEnd().HitsRectangle(bush.hitbox)

		if targetInBush && sourceInBush {
			// Both in the same bush then the bush has no effect
			continue
		}

		if targetInBush && !sourceInBush {
			// The target is in the bush but the source is not. There is no line of line of sight
			return false
		}

		if !targetInBush && !sourceInBush && line.HitsRectangle(bush.hitbox) {
			// The bush is acting as a wall
			return false
		}
	}

	return true
}

// Connect func
func (game *Game) Connect(client *ws.Client) {
	game.add <- client
}

// Disconnect func
func (game *Game) Disconnect(client *ws.Client) {
	game.remove <- client
}

func (game *Game) getClientChampion(client *ws.Client) *Champion {
	return game.clientChampions[client]
}

func (game *Game) getClientTeam(client *ws.Client) *Team {
	return game.clientTeams[client]
}

func (game *Game) setClientChampion(client *ws.Client, champion *Champion) {
	game.clientChampions[client] = champion
}

func (game *Game) setClientTeam(client *ws.Client, team *Team) {
	game.clientTeams[client] = team
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

func (game *Game) addClient(client *ws.Client) {
	client.Subscribe(game.global)
	client.WriteMessage("setup", NewSetupUpdate(game, client))

	champion := NewChampion(client.GetID())
	team := game.getNextTeam()

	team.addClient(client, champion)
	game.setClientChampion(client, champion)
	game.setClientTeam(client, team)

	// Send clients the updated teams
	game.global.Broadcast("update-teams", NewTeamsUpdate(game))
}

func (game *Game) removeClient(client *ws.Client) {
	client.Unsubscribe(game.global)
	team := game.getClientTeam(client)

	team.removeClient(client)
	delete(game.clientChampions, client)
	delete(game.clientTeams, client)

	// Send clients the updated teams
	game.global.Broadcast("update-teams", NewTeamsUpdate(game))
}
