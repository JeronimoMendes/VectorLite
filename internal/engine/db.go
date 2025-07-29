package engine

import (
	"VectorLite/internal/algorithms"
	"VectorLite/internal/vector"
)

func NewDatabase(name string, algorithm algorithms.SearchAlgorithm) *Database {
	return &Database{
		Name:      name,
		Algorithm: algorithm,
	}
}

func (database *Database) AddEntry(vector vector.Vector, metadata map[string]string) {
	database.NumberEntries++
	entry := algorithms.Entry{
		Vector:   vector,
		Metadata: metadata,
		Id:       database.NumberEntries,
	}
	database.Algorithm.AddEntry(entry)
}

func (database *Database) ListEntries() []algorithms.Entry {
	return database.Algorithm.ListEntries()
}

func (database *Database) Query(queryVector *vector.Vector, k int, metric string) []algorithms.Entry {
	return database.Algorithm.Query(queryVector, k, metric)
}

