package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"server/game/gameplay"
	"server/ws"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type createJSON struct {
	Code    string `json:"code"`
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type joinJSON struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type infoJSON struct {
	LiveGames   int `json:"liveGames"`
	LivePlayers int `json:"livePlayers"`
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func marshall(data interface{}) []byte {
	bytes, _ := json.Marshal(data)
	return bytes
}

func validateGame(code string, game *gameplay.Game) error {
	if game == nil {
		return fmt.Errorf("Game with code \"%s\" does not exist", code)
	} else if game.GetPlayerCount() == 10 {
		return errors.New("This game is full (10/10 players)")
	}
	return nil
}

func validateName(name string, game *gameplay.Game) error {
	if len(name) == 0 {
		return errors.New("Please enter a username")
	} else if len(name) < 3 {
		return errors.New("Please enter a username with at least 3 characters")
	} else if len(name) > 20 {
		return errors.New("Please enter a username with less than 20 characters")
	} else if game != nil && game.UsernameTaken(name) {
		return fmt.Errorf("The username \"%s\" has already been taken", name)
	}

	return nil
}

// InfoAPI func
func (mgr *Manager) InfoAPI(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	players := 0
	for _, game := range mgr.games {
		players += game.GetPlayerCount()
	}

	bytes, _ := json.Marshal(&infoJSON{
		LiveGames:   len(mgr.games),
		LivePlayers: players,
	})

	w.Write(bytes)
}

// GameCreateAPI func
func (mgr *Manager) GameCreateAPI(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	name := r.URL.Query().Get("name")

	// Validate the username
	if err := validateName(name, nil); err != nil {
		w.Write(marshall(&createJSON{
			Code:    "",
			Success: false,
			Error:   err.Error(),
		}))
		return
	}

	code := mgr.CreateGameCode()
	game, _ := mgr.CreateGame(code, false)

	go func() {
		// Stop the game after 15 minutes
		time.Sleep(time.Minute * 15)
		game.Stop()
		delete(mgr.games, code)
	}()

	w.Write(marshall(&createJSON{
		Code:    code,
		Success: true,
		Error:   "",
	}))
}

// GameJoinAPI func
func (mgr *Manager) GameJoinAPI(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	// Create the websocket using the upgrader
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	code := strings.ToUpper(r.URL.Query().Get("code"))
	name := r.URL.Query().Get("name")

	game := mgr.GetGame(code)
	client := ws.NewClient(conn, name)

	// Validate the username
	if err := validateName(name, game); err != nil {
		client.WriteMessage("connection", marshall(&joinJSON{
			Success: false,
			Error:   err.Error(),
		}))
		client.Close()
		return
	}

	// Validate the game
	if err := validateGame(code, game); err != nil {
		client.WriteMessage("connection", marshall(&joinJSON{
			Success: false,
			Error:   err.Error(),
		}))
		client.Close()
		return
	}

	// Successfully connected!
	client.WriteMessage("connection", marshall(&joinJSON{
		Success: true,
		Error:   "",
	}))

	// Get the name and game code the from the request
	go client.WriteMessages()
	go client.ReceiveMessages(game)

	game.Connect(client)    // connect client to the game
	client.Wait()           // block until the websocket disconnects
	game.Disconnect(client) // disconnect client from the game

	if game.GetPlayerCount() == 0 && game.Stop() {
		delete(mgr.games, code)
	}

}
