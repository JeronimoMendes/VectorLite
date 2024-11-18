package engine_test

import (
	"testing"

	"VectorLite/internal/engine"
	"VectorLite/internal/vector"

	"github.com/stretchr/testify/assert"
)

func TestNewDB(t *testing.T) {
	db := engine.NewDatabase()

	assert.NotNil(t, db)
	assert.Equal(t, 0, len(db.ListEntries()), "New database should have no entries")
}

func TestAddEntry(t *testing.T) {
	db := engine.NewDatabase()
	vec := vector.NewVector(1.5, 2.2)
	metadata := map[string]string{"text": "hello world"}

	db.AddEntry(*vec, metadata)

	assert.Equal(t, 1, len(db.ListEntries()), "Database should have one entry")
	entry := db.ListEntries()[0]
	assert.Equal(t, *vec, entry.Vector, "Vector should match the one added")
	assert.Equal(t, metadata, entry.Metadata, "Metadata should match the one added")
	assert.Equal(t, 1, entry.Id, "Id should be set correctly")
}

func TestListEntriesEmpty(t *testing.T) {
	db := engine.NewDatabase()
	entries := db.ListEntries()

	assert.Equal(t, 0, len(entries), "ListEntries should return empty list for new database")
}

func TestListEntriesWithEntries(t *testing.T) {
	db := engine.NewDatabase()
	vec1 := vector.NewVector(1.5, 2.2)
	vec2 := vector.NewVector(3.1, 4.4)
	metadata1 := map[string]string{"text": "entry1"}
	metadata2 := map[string]string{"text": "entry2"}

	db.AddEntry(*vec1, metadata1)
	db.AddEntry(*vec2, metadata2)

	entries := db.ListEntries()

	assert.Equal(t, 2, len(entries), "ListEntries should return all added entries")
	assert.Equal(t, *vec1, entries[0].Vector, "First entry vector should match")
	assert.Equal(t, metadata1, entries[0].Metadata, "First entry metadata should match")
	assert.Equal(t, *vec2, entries[1].Vector, "Second entry vector should match")
	assert.Equal(t, metadata2, entries[1].Metadata, "Second entry metadata should match")
}

func TestQuery(t *testing.T) {
	db := engine.NewDatabase()
	db.AddEntry(*vector.NewVector(1, 2, 3), map[string]string{"text": "entry1"})
	db.AddEntry(*vector.NewVector(4, 5, 6), map[string]string{"text": "entry2"})
	db.AddEntry(*vector.NewVector(7, 8, 9), map[string]string{"text": "entry3"})

	vectorA := vector.NewVector(2, 3, 4)

	// Test case 1: Basic functionality
	result := db.Query(vectorA, 2, "euclidean")
	assert.Equal(t, 2, len(result))
	// Additional assertions can be made based on the expected vector entries

	// Test case 2: Requesting more neighbors than available
	result = db.Query(vectorA, 5, "euclidean")
	assert.Equal(t, 3, len(result))

	// Test case 3: Using a different metric
	// Assuming implementation supports "manhattan" or other metrics
	result = db.Query(vectorA, 2, "manhattan")
	assert.Equal(t, 2, len(result))

	// Test case 5: Empty database
	emptyDatabase := engine.NewDatabase()
	result = emptyDatabase.Query(vectorA, 2, "euclidean")
	assert.Equal(t, 0, len(result))
}
