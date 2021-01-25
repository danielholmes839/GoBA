package gameplay

import "github.com/google/uuid"

// Manager struct to manage games
type Manager struct {
	games map[uuid.UUID]*Game
	add   chan *Game
	remove   chan *Game
}

// Run func
// func (mgr *Manager) Run() {
// 	select {
// 		case game := <- mgr.add:
// 			mgr.games[uuid.New()] = game
// 			go game.Run()

// 		case game := <- mgr.remove:
// 			game.Stop()
// 			delete(mgr.games, )
// 	}
// }
