package bruteforce

import (
	"VectorLite/internal/algorithms"
	"VectorLite/internal/vector"
	"math"
	"sort"
)

type Algorithm struct {
	entries []algorithms.Entry
}

type entryScore struct {
	Entry algorithms.Entry
	Score float64
}

func New() *Algorithm {
	return &Algorithm{
		entries: []algorithms.Entry{},
	}
}

func (a *Algorithm) AddEntry(vec vector.Vector, metadata map[string]string) {
	newEntry := algorithms.Entry{
		Vector:   vec,
		Metadata: metadata,
		Id:       len(a.entries) + 1,
	}
	a.entries = append(a.entries, newEntry)
}

func (a *Algorithm) ListEntries() []algorithms.Entry {
	return a.entries
}

func (a *Algorithm) Query(queryVector *vector.Vector, k int, metric string) []algorithms.Entry {
	// this is a brute force implementation of a knn algorithm
	returnEntriesScores := []entryScore{}
	highestScore := math.Inf(1) // this is actually the highest score in the return entries

	for _, entry := range a.entries {
		score := queryVector.Distance_score(&entry.Vector, metric)
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

	returnEntries := []algorithms.Entry{}
	for _, i := range returnEntriesScores {
		returnEntries = append(returnEntries, i.Entry)
	}

	return returnEntries
}