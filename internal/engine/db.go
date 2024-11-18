package engine

import (
	"VectorLite/internal/vector"
	"math"
	"sort"
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

func (database *Database) Query(query_vector *vector.Vector, k int, metric string) []Entry {
	// this is a brute force implementation of a knn algorithm
	returnEntriesScores := []entryScore{}
	highestScore := math.Inf(1) // this is actually the highest score in the return entries

	for _, entry := range database.entries {
		score := query_vector.Distance_score(&entry.Vector, metric)
		if score < highestScore || len(returnEntriesScores) < k {
			returnEntriesScores = append(returnEntriesScores, entryScore{Entry: entry, Score: score})

			// sorts the returnEntriesScores by score in ascending order
			sort.Slice(returnEntriesScores, func(i, j int) bool {
				return returnEntriesScores[i].Score > returnEntriesScores[j].Score
			})

			// here we need to remove
			if k < len(returnEntriesScores) {
				cut_from := len(returnEntriesScores) - k
				returnEntriesScores = returnEntriesScores[cut_from:]
			}

			highestScore = returnEntriesScores[0].Score
		}
	}

	returnEntries := []Entry{}
	for _, i := range returnEntriesScores {
		returnEntries = append(returnEntries, i.Entry)
	}

	return returnEntries
}
