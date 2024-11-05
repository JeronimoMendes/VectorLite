package vector_test

import (
	"VectorLite/internal/vector"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVector(t *testing.T) {
	// Test creating a vector with no values
	v := vector.NewVector()
	assert.NotNil(t, v, "Vector should not be nil")
	assert.Equal(t, 0, len(v.Values), "Vector should have no values when initialized with none")

	// Test creating a vector with multiple values
	values := []float64{1.1, 2.2, 3.3}
	v = vector.NewVector(values...)
	assert.NotNil(t, v, "Vector should not be nil")
	assert.Equal(t, values, v.Values, "Vector values should match the input values")
}

func TestNewVectorSingleValues(t *testing.T) {
	// Test creating a vector with a single value
	v := vector.NewVector(10.5)
	assert.NotNil(t, v, "Vector should not be nil")
	assert.Equal(t, 1, len(v.Values), "Vector should have one value")
	assert.Equal(t, 10.5, v.Values[0], "Vector value should match the single input value")
}
