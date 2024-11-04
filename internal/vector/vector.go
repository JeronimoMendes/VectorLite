package vector

type Vector struct {
	Values []float64
}

func NewVector(values ...float64) *Vector {
	return &Vector{Values: values}
}
