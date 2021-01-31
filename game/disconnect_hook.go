package game

import "server/game/gameplay"

// DisconnectHandler - executed when a player disconnects
// - Currently used to stop the game when the last player disconnects
type DisconnectHook func(mgr *Manager, game *gameplay.Game, code string)

// DefaultDisconnectHandler - stop the game when the player count hits 0
func DefaultDisconnectHook(mgr *Manager, game *gameplay.Game, code string) {
	if game.GetPlayerCount() > 0 {
		return
	}

	mgr.StopGame(code)
}

// EmptyDisconnectHandler - do nothing
func EmptyDisconnectHook(mgr *Manager, game *gameplay.Game, code string) {
}
