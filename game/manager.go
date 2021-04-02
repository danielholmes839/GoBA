package game

import (
	"errors"
	"math/rand"
	"server/game/gameplay"
	"sync"
	"time"
)

// ManagedGame struct
type ManagedGame struct {
	*gameplay.Game
	onDisconnect DisconnectHook
}

// Manager struct
type Manager struct {
	sync.Mutex
	games map[string]*ManagedGame
}

// NewGameManager func
func NewGameManager() *Manager {
	return &Manager{
		sync.Mutex{},
		make(map[string]*ManagedGame),
	}
}

// StopGame func
func (mgr *Manager) StopGame(code string) error {
	mgr.Lock()
	defer mgr.Unlock()

	// The game does not exist
	game := mgr.get(code)
	if game == nil {
		return errors.New("This game cannot be found")
	}

	// The game is stopped
	success := game.Stop()
	if !success {
		return errors.New("This game is not running")
	}

	delete(mgr.games, code)
	return nil
}

// CreateGameDefault - create a game with particular code
func (mgr *Manager) CreateGameDefault() (string, *gameplay.Game) {
	mgr.Lock()
	defer mgr.Unlock()

	code := mgr.createUniqueGameCode()
	game := gameplay.NewGame(64)
	mgr.set(code, &ManagedGame{game, DefaultDisconnectHook})

	// Stop the game after 15 minutes
	go func() {
		time.Sleep(time.Minute * 15)
		mgr.StopGame(code)
	}()

	return code, game
}

// CreateGameManually func
func (mgr *Manager) CreateGameManually(code string, onDisconnect DisconnectHook) error {
	mgr.Lock()
	defer mgr.Unlock()

	// Don't do anything on disconnect by default
	if onDisconnect == nil {
		onDisconnect = EmptyDisconnectHook
	}

	// Check if the code is available
	if !mgr.taken(code) {
		return errors.New("This code is already taken")
	}

	// Create the game
	game := gameplay.NewGame(64)
	mgr.set(code, &ManagedGame{game, onDisconnect})

	return nil
}

func (mgr *Manager) createGameCode() string {
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	r := make([]rune, 6)
	for i := range r {
		r[i] = letters[rand.Intn(len(letters))]
	}
	code := string(r)
	return code
}

// CreateGameCode - SHOULD ONLY BE USED WHEN mgr.creating IS LOCKED
func (mgr *Manager) createUniqueGameCode() string {
	code := mgr.createGameCode()
	for !mgr.taken(code) {
		code = mgr.createGameCode()
	}
	return code
}

func (mgr *Manager) taken(code string) bool {
	return mgr.get(code) == nil
}

func (mgr *Manager) get(code string) *ManagedGame {
	return mgr.games[code]
}

func (mgr *Manager) set(code string, game *ManagedGame) {
	mgr.games[code] = game
}
