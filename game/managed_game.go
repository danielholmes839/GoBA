package game

import "server/game/gameplay"



// ManagedGame struct
type ManagedGame struct {
	game         *gameplay.Game
	onDisconnect DisconnectHook
}
