package game

// DisconnectHandler - executed when a player disconnects
// - Currently used to stop the game when the last player disconnects
type DisconnectHook func(mgr *Manager, game *ManagedGame, code string)

// DefaultDisconnectHandler - stop the game when the player count hits 0
func DefaultDisconnectHook(mgr *Manager, game *ManagedGame, code string) {
	if game.GetPlayerCount() > 0 {
		return
	}

	mgr.StopGame(code)
}

// EmptyDisconnectHandler - do nothing
func EmptyDisconnectHook(mgr *Manager, game *ManagedGame, code string) {
}
