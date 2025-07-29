package engine

import (
	"VectorLite/internal/algorithms"
	"errors"
)

var (
	ErrDatabaseExists   = errors.New("database already exists")
	ErrDatabaseNotFound = errors.New("database not found")
)

type Database struct {
	Name          string
	Algorithm     algorithms.SearchAlgorithm
	NumberEntries int
}

type DatabaseManager struct {
	databases map[string]*Database
}

func NewDatabaseManager() *DatabaseManager {
	return &DatabaseManager{
		databases: make(map[string]*Database),
	}
}

func (dm *DatabaseManager) CreateDatabase(name string, algorithm algorithms.SearchAlgorithm) error {
	if _, exists := dm.databases[name]; exists {
		return ErrDatabaseExists
	}
	
	db := &Database{
		Name:      name,
		Algorithm: algorithm,
	}
	dm.databases[name] = db
	return nil
}

func (dm *DatabaseManager) GetDatabase(name string) (*Database, error) {
	db, exists := dm.databases[name]
	if !exists {
		return nil, ErrDatabaseNotFound
	}
	return db, nil
}

func (dm *DatabaseManager) ListDatabases() []string {
	names := make([]string, 0, len(dm.databases))
	for name := range dm.databases {
		names = append(names, name)
	}
	return names
}

func (dm *DatabaseManager) DeleteDatabase(name string) error {
	if _, exists := dm.databases[name]; !exists {
		return ErrDatabaseNotFound
	}
	delete(dm.databases, name)
	return nil
}

