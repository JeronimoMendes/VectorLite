package state

import (
	"VectorLite/internal/engine"
	"sync"
)

type GlobalState struct {
	Database *engine.Database
	Mu       sync.Mutex
}

var State = GlobalState{
	Database: engine.NewDatabase(),
}
