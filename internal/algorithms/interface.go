package algorithms

import "VectorLite/internal/vector"

type SearchAlgorithm interface {
	AddEntry(entry Entry)
	Query(queryVector *vector.Vector, k int, metric string) []Entry
	ListEntries() []Entry
}

type Entry struct {
	Vector   vector.Vector
	Metadata map[string]string
	Id       int
}