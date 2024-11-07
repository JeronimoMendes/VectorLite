package engine

import (
	"VectorLite/internal/vector"
)

func NewDatabase() *Database {
	return &Database{
		entries: []Entry{},
	}
}

func (database *Database) AddEntry(vector vector.Vector, metadata map[string]string) {
	newEntry := Entry{
		Vector:   vector,
		Metadata: metadata,
		Id:       len(database.entries) + 1,
	}
	database.entries = append(database.entries, newEntry)
}

func (database *Database) ListEntries() []Entry {
	return database.entries
}
