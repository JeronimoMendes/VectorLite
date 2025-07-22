package algorithms

import "VectorLite/internal/vector"

type SearchAlgorithm interface {
	AddEntry(vector vector.Vector, metadata map[string]string)
	Query(queryVector *vector.Vector, k int, metric string) []Entry
	ListEntries() []Entry
}

type Entry struct {
	Vector   vector.Vector
	Metadata map[string]string
	Id       int
}