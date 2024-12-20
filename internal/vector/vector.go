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

func (v1 *Vector) Normalize() *Vector {
	magnitude := v1.Magnitude()
	new_values := []float64{}
	for _, value := range v1.Values {
		new_value := value / magnitude
		new_values = append(new_values, new_value)
	}
	new_vector := Vector{Values: new_values}
	return &new_vector
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

func (v1 *Vector) Euclidean_distance(v2 *Vector) float64 {
	x := 0.0
	for i, value1 := range v1.Values {
		value2 := v2.Values[i]
		x += math.Pow(value1-value2, 2)
	}
	return math.Sqrt(x)
}

func (v1 *Vector) Distance_score(v2 *Vector, metric string) float64 {
	score := math.Inf(1)
	switch metric {
	case "cosine":
		score = 1 - (1+v1.Cosine_similarity(v2))/2
	case "dot_product":
		score = 1 - (1+v1.Normalize().Cosine_similarity(v2.Normalize()))/2
	case "euclidean":
		score = v1.Euclidean_distance(v2)
	}
	return score
}
