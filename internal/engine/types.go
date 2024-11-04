package engine

import "VectorLite/internal/vector"

type Database struct {
	entries []Entry
}

type Entry struct {
	Vector   vector.Vector
	Metadata map[string]string
	Id       int
}
