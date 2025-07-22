package state

import (
	"VectorLite/internal/algorithms/bruteforce"
	"VectorLite/internal/engine"
	"sync"
)

type GlobalState struct {
	Database *engine.Database
	Mu       sync.Mutex
}

var State = GlobalState{
	Database: engine.NewDatabase(bruteforce.New()),
}
