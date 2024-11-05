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
	assert.Equal(t, 0, len(engine.ListEntries(db)), "New database should have no entries")
}

func TestAddEntry(t *testing.T) {
	db := engine.NewDatabase()
	vec := vector.NewVector(1.5, 2.2)
	metadata := map[string]string{"text": "hello world"}

	engine.AddEntry(db, *vec, metadata)

	assert.Equal(t, 1, len(engine.ListEntries(db)), "Database should have one entry")
	entry := engine.ListEntries(db)[0]
	assert.Equal(t, *vec, entry.Vector, "Vector should match the one added")
	assert.Equal(t, metadata, entry.Metadata, "Metadata should match the one added")
	assert.Equal(t, 1, entry.Id, "Id should be set correctly")
}

func TestListEntriesEmpty(t *testing.T) {
	db := engine.NewDatabase()
	entries := engine.ListEntries(db)

	assert.Equal(t, 0, len(entries), "ListEntries should return empty list for new database")
}

func TestListEntriesWithEntries(t *testing.T) {
	db := engine.NewDatabase()
	vec1 := vector.NewVector(1.5, 2.2)
	vec2 := vector.NewVector(3.1, 4.4)
	metadata1 := map[string]string{"text": "entry1"}
	metadata2 := map[string]string{"text": "entry2"}

	engine.AddEntry(db, *vec1, metadata1)
	engine.AddEntry(db, *vec2, metadata2)

	entries := engine.ListEntries(db)

	assert.Equal(t, 2, len(entries), "ListEntries should return all added entries")
	assert.Equal(t, *vec1, entries[0].Vector, "First entry vector should match")
	assert.Equal(t, metadata1, entries[0].Metadata, "First entry metadata should match")
	assert.Equal(t, *vec2, entries[1].Vector, "Second entry vector should match")
	assert.Equal(t, metadata2, entries[1].Metadata, "Second entry metadata should match")
}
