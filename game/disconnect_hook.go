package game

import "server/game/gameplay"

// DisconnectHook
// Currently used to stop the game when the last player disconnects
type DisconnectHook func(mgr *Manager, game gameplay.IGame, code string)

// DefaultDisconnectHandler
// stops the game when it hits 0 players
func DefaultDisconnectHook(mgr *Manager, game gameplay.IGame, code string) {
	if game.GetPlayerCount() > 0 {
		return
	}

	mgr.Stop(code)
}

// EmptyDisconnectHandler
func EmptyDisconnectHook(mgr *Manager, game gameplay.IGame, code string) {
}
