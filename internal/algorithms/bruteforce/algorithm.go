package bruteforce

import (
	"VectorLite/internal/algorithms"
	"VectorLite/internal/vector"
	"math"
	"sort"
)

type Algorithm struct {
	entries 	[]algorithms.Entry
	idCounter 	int
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

func (a *Algorithm) AddEntry(entry algorithms.Entry) {
	a.entries = append(a.entries, entry)
}

func (a *Algorithm) ListEntries() []algorithms.Entry {
	return a.entries
}

func (a *Algorithm) Query(queryVector *vector.Vector, k int, metric string) []algorithms.Entry {
	// Handle edge case where k=0
	if k <= 0 {
		return []algorithms.Entry{}
	}

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

			// Only update highestScore if we have entries
			if len(returnEntriesScores) > 0 {
				highestScore = returnEntriesScores[0].Score
			}
		}
	}

	returnEntries := []algorithms.Entry{}
	for _, i := range returnEntriesScores {
		returnEntries = append(returnEntries, i.Entry)
	}

	return returnEntries
}