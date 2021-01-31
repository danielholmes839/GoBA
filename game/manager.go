package game

import (
	"errors"
	"math/rand"
	"server/game/gameplay"
	"sync"
	"time"
)

// Manager struct
type Manager struct {
	managedGames map[string]*ManagedGame
	accessing    *sync.Mutex
	editing      *sync.Mutex
}

// NewGameManager func
func NewGameManager() *Manager {
	return &Manager{
		managedGames: make(map[string]*ManagedGame),
		accessing:    &sync.Mutex{},
		editing:      &sync.Mutex{},
	}
}

// StopGame func
func (mgr *Manager) StopGame(code string) error {
	mgr.editing.Lock()
	defer mgr.editing.Unlock()

	managedGame := mgr.get(code)
	if managedGame == nil {
		return errors.New("This game cannot be found")
	}

	success := managedGame.game.Stop()
	if !success {
		return errors.New("This game is not running")
	}

	delete(mgr.managedGames, code)
	return nil
}

// CreateGameDefault - create a game with particular code
func (mgr *Manager) CreateGameDefault() (string, *gameplay.Game) {
	mgr.editing.Lock()
	defer mgr.editing.Unlock()

	code := mgr.createUniqueGameCode()
	game := gameplay.NewGame(64)
	mgr.set(code, &ManagedGame{game, EmptyDisconnectHook})

	// Stop the game after 15 minutes
	go func() {
		time.Sleep(time.Minute * 15)
		mgr.StopGame(code)
	}()

	return code, game
}

// CreateGameManually func
func (mgr *Manager) CreateGameManually(code string, onDisconnect DisconnectHook) error {
	mgr.editing.Lock()
	defer mgr.editing.Unlock()

	if onDisconnect == nil {
		onDisconnect = EmptyDisconnectHook
	}

	if !mgr.taken(code) {
		return errors.New("This code is already taken")
	}

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
	mgr.accessing.Lock()
	defer mgr.accessing.Unlock()
	return mgr.managedGames[code]
}

func (mgr *Manager) set(code string, managedGame *ManagedGame) {
	mgr.accessing.Lock()
	defer mgr.accessing.Unlock()
	mgr.managedGames[code] = managedGame
}
