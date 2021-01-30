package game

import (
	"errors"
	"math/rand"
	"server/game/gameplay"
)

// Manager struct to manage games
type Manager struct {
	games  map[string]*gameplay.Game
	add    chan *gameplay.Game
	remove chan *gameplay.Game
}

// NewGameManager func
func NewGameManager() *Manager {
	return &Manager{
		games: map[string]*gameplay.Game{},
	}
}

// GetGame func
func (mgr *Manager) GetGame(code string) *gameplay.Game {
	return mgr.games[code]
}

// CreateGame - create a game with particular code
func (mgr *Manager) CreateGame(code string, permanent bool) (*gameplay.Game, error) {
	if !mgr.isCodeAvailable(code) {
		return nil, errors.New("This code is taken")
	}

	game := gameplay.NewGame(64, permanent)
	go game.Run()
	mgr.games[code] = game
	return game, nil
}

// CreateGameCode - the code used to indentify the game
func (mgr *Manager) CreateGameCode() string {
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	r := make([]rune, 6)
	for i := range r {
		r[i] = letters[rand.Intn(len(letters))]
	}
	code := string(r)

	// Create a new code in case of duplicates
	if !mgr.isCodeAvailable(code) {
		return mgr.CreateGameCode()
	}

	return code
}

func (mgr *Manager) isCodeAvailable(code string) bool {
	return mgr.games[code] == nil
}
