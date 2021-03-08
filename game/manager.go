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
	gameplay.Game
	OnDisconnect DisconnectHook
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
func (mgr *Manager) Stop(code string) error {
	mgr.Lock()
	defer mgr.Unlock()

	game := mgr.get(code)

	// The game does not exist
	if game == nil {
		return errors.New("This game cannot be found")
	}

	// The game is stopped
	if err := game.Stop(); err != nil {
		return err
	}

	delete(mgr.games, code)
	return nil
}

func (mgr *Manager) CreateDefaultGame() (string, *gameplay.Game) {
	mgr.Lock()
	defer mgr.Unlock()

	// Create the game
	game := gameplay.NewGame(64)
	go game.Run()

	// Set the game 
	code := mgr.createUniqueGameCode()
	mgr.set(code, &ManagedGame{*game, DefaultDisconnectHook})

	// Stop the game after 15 minutes
	go func() {
		time.Sleep(time.Minute * 15)
		mgr.Stop(code)
	}()

	return code, game
}

func (mgr *Manager) CreateCustomGame(code string, onDisconnect DisconnectHook) error {
	mgr.Lock()
	defer mgr.Unlock()

	// Don't do anything on disconnect by default
	if onDisconnect == nil {
		onDisconnect = EmptyDisconnectHook
	}

	// Check if the code is available
	if mgr.taken(code) {
		return errors.New("This code is already taken")
	}

	// Create the game
	game := gameplay.NewGame(64)
	go game.Run()
	mgr.set(code, &ManagedGame{*game, onDisconnect})

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
	for mgr.taken(code) {
		code = mgr.createGameCode()
	}
	return code
}

func (mgr *Manager) taken(code string) bool {
	return mgr.get(code) != nil
}

func (mgr *Manager) get(code string) *ManagedGame {
	return mgr.games[code]
}

func (mgr *Manager) set(code string, game *ManagedGame) {
	mgr.games[code] = game
}
