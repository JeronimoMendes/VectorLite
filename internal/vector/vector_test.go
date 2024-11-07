package vector_test

import (
	"VectorLite/internal/vector"
	"math"
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

func TestMagnitude(t *testing.T) {
	tests := []struct {
		vector vector.Vector
		want   float64
	}{
		{vector.Vector{Values: []float64{3, 4}}, 5},
		{vector.Vector{Values: []float64{1, 2, 2}}, 3},
		{vector.Vector{Values: []float64{0, 0, 0}}, 0},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := tt.vector.Magnitude(); math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Magnitude() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDotProduct(t *testing.T) {
	tests := []struct {
		v1, v2 vector.Vector
		want   float64
	}{
		{vector.Vector{Values: []float64{1, 3, -5}}, vector.Vector{Values: []float64{4, -2, -1}}, 3},
		{vector.Vector{Values: []float64{1, 2, 3}}, vector.Vector{Values: []float64{4, 5, 6}}, 32},
		{vector.Vector{Values: []float64{1, 0, 0}}, vector.Vector{Values: []float64{0, 1, 0}}, 0},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := tt.v1.Dot_product(&tt.v2); math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Dot_product() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		v1, v2 vector.Vector
		want   float64
	}{
		{vector.Vector{Values: []float64{1, 0}}, vector.Vector{Values: []float64{0, 1}}, 0},
		{vector.Vector{Values: []float64{1, 0}}, vector.Vector{Values: []float64{1, 0}}, 1},
		{vector.Vector{Values: []float64{1, 2, 3}}, vector.Vector{Values: []float64{4, 5, 6}}, 0.974631846},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := tt.v1.Cosine_similarity(&tt.v2); math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Cosine_similarity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEuclideanDistance(t *testing.T) {
	tests := []struct {
		v1, v2 vector.Vector
		want   float64
	}{
		{vector.Vector{Values: []float64{0, 0}}, vector.Vector{Values: []float64{0, 0}}, 0},
		{vector.Vector{Values: []float64{1, 2}}, vector.Vector{Values: []float64{4, 6}}, 5},
		{vector.Vector{Values: []float64{1, 0, 0}}, vector.Vector{Values: []float64{0, 1, 0}}, math.Sqrt(2)},
		{vector.Vector{Values: []float64{1, 2, 3}}, vector.Vector{Values: []float64{1, 2, 3}}, 0},
		{vector.Vector{Values: []float64{-1, -2, -3}}, vector.Vector{Values: []float64{-4, -5, -6}}, math.Sqrt(27)},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := tt.v1.Euclidean_distance(&tt.v2); math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Euclidean_distance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		input    vector.Vector
		expected vector.Vector
	}{
		{vector.Vector{Values: []float64{3, 4}}, vector.Vector{Values: []float64{0.6, 0.8}}},
		{vector.Vector{Values: []float64{5, 0}}, vector.Vector{Values: []float64{1, 0}}},
		{vector.Vector{Values: []float64{0, 0}}, vector.Vector{Values: []float64{0, 0}}},
		{vector.Vector{Values: []float64{-2, -2}}, vector.Vector{Values: []float64{-math.Sqrt(2) / 2, -math.Sqrt(2) / 2}}},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			normalized_vector := tt.input.Normalize()
			for i, v := range normalized_vector.Values {
				if math.Abs(v-tt.expected.Values[i]) > 1e-9 {
					t.Errorf("Normalize() = %v, want %v", tt.input.Values, tt.expected.Values)
				}
			}

			// Check if normalized vector's magnitude is 1 (with tolerance for floating point arithmetic)
			if mag := normalized_vector.Magnitude(); math.Abs(mag-1) > 1e-9 && mag != 0 { // ignore the zero vector
				t.Errorf("Normalized vector's magnitude is %v, want 1", mag)
			}
		})
	}
}
