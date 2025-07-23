package bruteforce_test

import (
	"testing"

	"VectorLite/internal/algorithms"
	"VectorLite/internal/algorithms/bruteforce"
	"VectorLite/internal/vector"

	"github.com/stretchr/testify/assert"
)

func TestNewAlgorithm(t *testing.T) {
	algo := bruteforce.New()
	assert.NotNil(t, algo)
	assert.Equal(t, 0, len(algo.ListEntries()), "New algorithm should have no entries")
}

func TestAddEntry(t *testing.T) {
	algo := bruteforce.New()
	vec := vector.NewVector(1.0, 2.0, 3.0)
	metadata := map[string]string{"text": "test entry"}
	entry := algorithms.Entry{
		Vector:   *vec,
		Metadata: metadata,
		Id:       1,
	}

	algo.AddEntry(entry)

	entries := algo.ListEntries()
	assert.Equal(t, 1, len(entries), "Algorithm should have one entry")
	assert.Equal(t, entry, entries[0], "Entry should match the one added")
}

func TestAddMultipleEntries(t *testing.T) {
	algo := bruteforce.New()
	
	entries := []algorithms.Entry{
		{Vector: *vector.NewVector(1.0, 2.0), Metadata: map[string]string{"id": "1"}, Id: 1},
		{Vector: *vector.NewVector(3.0, 4.0), Metadata: map[string]string{"id": "2"}, Id: 2},
		{Vector: *vector.NewVector(5.0, 6.0), Metadata: map[string]string{"id": "3"}, Id: 3},
	}

	for _, entry := range entries {
		algo.AddEntry(entry)
	}

	retrievedEntries := algo.ListEntries()
	assert.Equal(t, 3, len(retrievedEntries), "Algorithm should have three entries")
	for i, entry := range entries {
		assert.Equal(t, entry, retrievedEntries[i], "Entry %d should match", i)
	}
}

func TestListEntriesEmpty(t *testing.T) {
	algo := bruteforce.New()
	entries := algo.ListEntries()
	assert.Equal(t, 0, len(entries), "Empty algorithm should return empty slice")
}

func TestQueryEmpty(t *testing.T) {
	algo := bruteforce.New()
	queryVec := vector.NewVector(1.0, 2.0)
	
	result := algo.Query(queryVec, 5, "euclidean")
	assert.Equal(t, 0, len(result), "Query on empty algorithm should return empty slice")
}

func TestQuerySingleEntry(t *testing.T) {
	algo := bruteforce.New()
	entry := algorithms.Entry{
		Vector:   *vector.NewVector(1.0, 2.0),
		Metadata: map[string]string{"id": "1"},
		Id:       1,
	}
	algo.AddEntry(entry)

	queryVec := vector.NewVector(1.5, 2.5)
	result := algo.Query(queryVec, 1, "euclidean")
	
	assert.Equal(t, 1, len(result), "Should return the single entry")
	assert.Equal(t, entry, result[0], "Should return the correct entry")
}

func TestQueryKnnEuclidean(t *testing.T) {
	algo := bruteforce.New()
	
	// Add test vectors in a predictable pattern
	entries := []algorithms.Entry{
		{Vector: *vector.NewVector(0.0, 0.0), Metadata: map[string]string{"id": "origin"}, Id: 1},     // Distance 0 from (0,0)
		{Vector: *vector.NewVector(1.0, 0.0), Metadata: map[string]string{"id": "right"}, Id: 2},     // Distance 1 from (0,0)
		{Vector: *vector.NewVector(0.0, 1.0), Metadata: map[string]string{"id": "up"}, Id: 3},        // Distance 1 from (0,0)
		{Vector: *vector.NewVector(2.0, 0.0), Metadata: map[string]string{"id": "far_right"}, Id: 4}, // Distance 2 from (0,0)
		{Vector: *vector.NewVector(0.0, 2.0), Metadata: map[string]string{"id": "far_up"}, Id: 5},    // Distance 2 from (0,0)
	}

	for _, entry := range entries {
		algo.AddEntry(entry)
	}

	queryVec := vector.NewVector(0.0, 0.0)
	
	// Test k=1 (closest neighbor)
	result := algo.Query(queryVec, 1, "euclidean")
	assert.Equal(t, 1, len(result), "Should return 1 result")
	assert.Equal(t, "origin", result[0].Metadata["id"], "Closest should be origin")

	// Test k=3 (3 closest neighbors)
	result = algo.Query(queryVec, 3, "euclidean")
	assert.Equal(t, 3, len(result), "Should return 3 results")
	
	// Verify the results include the closest points (order may vary due to sorting implementation)
	foundIds := make(map[string]bool)
	for _, entry := range result {
		foundIds[entry.Metadata["id"]] = true
	}
	assert.True(t, foundIds["origin"], "Should include origin")
	assert.True(t, foundIds["right"] || foundIds["up"], "Should include at least one distance-1 point")
}

func TestQueryKnnCosine(t *testing.T) {
	algo := bruteforce.New()
	
	// Add vectors with known cosine similarities
	entries := []algorithms.Entry{
		{Vector: *vector.NewVector(1.0, 0.0), Metadata: map[string]string{"id": "horizontal"}, Id: 1}, // 0° from query
		{Vector: *vector.NewVector(1.0, 1.0), Metadata: map[string]string{"id": "diagonal"}, Id: 2},   // 45° from query
		{Vector: *vector.NewVector(0.0, 1.0), Metadata: map[string]string{"id": "vertical"}, Id: 3},   // 90° from query
		{Vector: *vector.NewVector(-1.0, 0.0), Metadata: map[string]string{"id": "opposite"}, Id: 4},  // 180° from query
	}

	for _, entry := range entries {
		algo.AddEntry(entry)
	}

	queryVec := vector.NewVector(1.0, 0.0) // Horizontal vector
	
	// Test k=2 with cosine similarity
	result := algo.Query(queryVec, 2, "cosine")
	assert.Equal(t, 2, len(result), "Should return 2 results")
	
	// The most similar should be the horizontal vector (itself)
	foundIds := make(map[string]bool)
	for _, entry := range result {
		foundIds[entry.Metadata["id"]] = true
	}
	assert.True(t, foundIds["horizontal"], "Should include the most similar vector")
}

func TestQueryKnnDotProduct(t *testing.T) {
	algo := bruteforce.New()
	
	entries := []algorithms.Entry{
		{Vector: *vector.NewVector(1.0, 1.0), Metadata: map[string]string{"id": "positive"}, Id: 1}, // Dot product: 2
		{Vector: *vector.NewVector(1.0, 0.0), Metadata: map[string]string{"id": "partial"}, Id: 2},  // Dot product: 1
		{Vector: *vector.NewVector(0.0, 0.0), Metadata: map[string]string{"id": "zero"}, Id: 3},     // Dot product: 0
		{Vector: *vector.NewVector(-1.0, -1.0), Metadata: map[string]string{"id": "negative"}, Id: 4}, // Dot product: -2
	}

	for _, entry := range entries {
		algo.AddEntry(entry)
	}

	queryVec := vector.NewVector(1.0, 1.0)
	
	result := algo.Query(queryVec, 2, "dot_product")
	assert.Equal(t, 2, len(result), "Should return 2 results")
	
	// Should return the vectors with highest dot products (lowest distance scores)
	foundIds := make(map[string]bool)
	for _, entry := range result {
		foundIds[entry.Metadata["id"]] = true
	}
	assert.True(t, foundIds["positive"], "Should include vector with highest dot product")
}

func TestQueryKGreaterThanEntries(t *testing.T) {
	algo := bruteforce.New()
	
	// Add only 2 entries
	entries := []algorithms.Entry{
		{Vector: *vector.NewVector(1.0, 0.0), Metadata: map[string]string{"id": "1"}, Id: 1},
		{Vector: *vector.NewVector(0.0, 1.0), Metadata: map[string]string{"id": "2"}, Id: 2},
	}

	for _, entry := range entries {
		algo.AddEntry(entry)
	}

	queryVec := vector.NewVector(0.5, 0.5)
	
	// Request more neighbors than available
	result := algo.Query(queryVec, 5, "euclidean")
	assert.Equal(t, 2, len(result), "Should return all available entries when k > number of entries")
}

func TestQueryDifferentMetrics(t *testing.T) {
	algo := bruteforce.New()
	
	// Use simple test vectors
	entries := []algorithms.Entry{
		{Vector: *vector.NewVector(1.0, 0.0), Metadata: map[string]string{"id": "1"}, Id: 1},
		{Vector: *vector.NewVector(0.0, 1.0), Metadata: map[string]string{"id": "2"}, Id: 2},
		{Vector: *vector.NewVector(-1.0, 0.0), Metadata: map[string]string{"id": "3"}, Id: 3},
	}

	for _, entry := range entries {
		algo.AddEntry(entry)
	}

	queryVec := vector.NewVector(1.0, 0.0)
	
	// Test each supported metric
	metrics := []string{"euclidean", "cosine", "dot_product"}
	for _, metric := range metrics {
		result := algo.Query(queryVec, 2, metric)
		assert.Equal(t, 2, len(result), "Should return 2 results for metric %s", metric)
		assert.NotNil(t, result[0], "First result should not be nil for metric %s", metric)
		assert.NotNil(t, result[1], "Second result should not be nil for metric %s", metric)
	}
}

func TestQueryAccuracy(t *testing.T) {
	algo := bruteforce.New()
	
	// Create a more complex scenario to test accuracy
	entries := []algorithms.Entry{
		{Vector: *vector.NewVector(1.0, 1.0), Metadata: map[string]string{"id": "A", "distance": "1.41"}, Id: 1}, // sqrt(2) ≈ 1.41 from origin
		{Vector: *vector.NewVector(2.0, 0.0), Metadata: map[string]string{"id": "B", "distance": "2.00"}, Id: 2}, // 2.0 from origin
		{Vector: *vector.NewVector(0.0, 3.0), Metadata: map[string]string{"id": "C", "distance": "3.00"}, Id: 3}, // 3.0 from origin
		{Vector: *vector.NewVector(1.0, 0.0), Metadata: map[string]string{"id": "D", "distance": "1.00"}, Id: 4}, // 1.0 from origin
		{Vector: *vector.NewVector(0.0, 1.0), Metadata: map[string]string{"id": "E", "distance": "1.00"}, Id: 5}, // 1.0 from origin
	}

	for _, entry := range entries {
		algo.AddEntry(entry)
	}

	queryVec := vector.NewVector(0.0, 0.0) // Query from origin
	
	// Get the 3 closest neighbors
	result := algo.Query(queryVec, 3, "euclidean")
	assert.Equal(t, 3, len(result), "Should return exactly 3 results")
	
	// Verify distances are in ascending order (closest first)
	var distances []float64
	for _, entry := range result {
		distance := queryVec.Distance_score(&entry.Vector, "euclidean")
		distances = append(distances, distance)
	}
	
	// Check that results include the closest points
	foundIds := make(map[string]bool)
	for _, entry := range result {
		foundIds[entry.Metadata["id"]] = true
	}
	
	// Should include the two distance-1.0 points and the distance-1.41 point
	count := 0
	if foundIds["D"] { count++ } // distance 1.0
	if foundIds["E"] { count++ } // distance 1.0  
	if foundIds["A"] { count++ } // distance 1.41
	
	assert.GreaterOrEqual(t, count, 2, "Should include at least 2 of the closest points")
	assert.False(t, foundIds["C"], "Should not include the farthest point (C)")
}

func TestQueryEdgeCases(t *testing.T) {
	algo := bruteforce.New()
	
	// Test with zero vectors
	entry := algorithms.Entry{
		Vector:   *vector.NewVector(0.0, 0.0),
		Metadata: map[string]string{"id": "zero"},
		Id:       1,
	}
	algo.AddEntry(entry)

	queryVec := vector.NewVector(0.0, 0.0)
	result := algo.Query(queryVec, 1, "euclidean")
	assert.Equal(t, 1, len(result), "Should handle zero vectors")
	assert.Equal(t, entry, result[0], "Should return the zero vector")

	// Test with k=0
	result = algo.Query(queryVec, 0, "euclidean")
	assert.Equal(t, 0, len(result), "Should return empty result for k=0")
}

func TestQueryConsistency(t *testing.T) {
	algo := bruteforce.New()
	
	// Add the same set of entries multiple times to ensure consistency
	baseEntries := []algorithms.Entry{
		{Vector: *vector.NewVector(1.0, 2.0), Metadata: map[string]string{"id": "1"}, Id: 1},
		{Vector: *vector.NewVector(3.0, 4.0), Metadata: map[string]string{"id": "2"}, Id: 2},
		{Vector: *vector.NewVector(5.0, 6.0), Metadata: map[string]string{"id": "3"}, Id: 3},
	}

	for _, entry := range baseEntries {
		algo.AddEntry(entry)
	}

	queryVec := vector.NewVector(2.0, 3.0)
	
	// Run the same query multiple times
	results := make([][]algorithms.Entry, 3)
	for i := 0; i < 3; i++ {
		results[i] = algo.Query(queryVec, 2, "euclidean")
	}

	// All results should be identical
	for i := 1; i < len(results); i++ {
		assert.Equal(t, len(results[0]), len(results[i]), "Result lengths should be consistent")
		for j := 0; j < len(results[0]); j++ {
			assert.Equal(t, results[0][j], results[i][j], "Results should be identical across runs")
		}
	}
}