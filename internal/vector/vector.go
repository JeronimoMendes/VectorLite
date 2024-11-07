package vector

import "math"

type Vector struct {
	Values []float64
}

func NewVector(values ...float64) *Vector {
	return &Vector{Values: values}
}

func (vector *Vector) Magnitude() float64 {
	x := 0.0
	for _, value := range vector.Values {
		x += math.Pow(value, 2)
	}
	return math.Sqrt(x)
}

func (v1 *Vector) Dot_product(v2 *Vector) float64 {
	dot_product := 0.0
	for i, value1 := range v1.Values {
		value2 := v2.Values[i]
		dot_product += value1 * value2
	}
	return dot_product
}

func (v1 *Vector) Cosine_similarity(v2 *Vector) float64 {
	return v1.Dot_product(v2) / (v1.Magnitude() * v2.Magnitude())
}
