package engine

import (
	"VectorLite/internal/algorithms"
	"VectorLite/internal/vector"
)

func NewDatabase(algorithm algorithms.SearchAlgorithm) *Database {
	return &Database{
		algorithm: algorithm,
	}
}

func (database *Database) AddEntry(vector vector.Vector, metadata map[string]string) {
	database.algorithm.AddEntry(vector, metadata)
}

func (database *Database) ListEntries() []algorithms.Entry {
	return database.algorithm.ListEntries()
}

func (database *Database) Query(queryVector *vector.Vector, k int, metric string) []algorithms.Entry {
	return database.algorithm.Query(queryVector, k, metric)
}

