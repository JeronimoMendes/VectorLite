package state

import (
	"VectorLite/internal/engine"
	"sync"
)

type GlobalState struct {
	DatabaseManager *engine.DatabaseManager
	Mu              sync.Mutex
}

var State = GlobalState{
	DatabaseManager: engine.NewDatabaseManager(),
}
