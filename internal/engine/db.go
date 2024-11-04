package engine

import (
	"VectorLite/internal/vector"
)

func NewDatabase() *Database {
	return &Database{
		entries: []Entry{},
	}
}

func AddEntry(database *Database, vector vector.Vector, metadata map[string]string) {
	newEntry := Entry{
		Vector:   vector,
		Metadata: metadata,
		Id:       len(database.entries) + 1,
	}
	database.entries = append(database.entries, newEntry)
}

func ListEntries(database *Database) []Entry {
	return database.entries
}
