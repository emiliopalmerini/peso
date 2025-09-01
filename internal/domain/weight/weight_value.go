package weight

import (
	"errors"
	"fmt"
)

type WeightValue float64

const (
	MinWeightValue = 10.0
	MaxWeightValue = 500.0
)

var (
	ErrWeightTooLow  = errors.New("weight must be at least 10kg")
	ErrWeightTooHigh = errors.New("weight must be at most 500kg")
	ErrWeightInvalid = errors.New("weight must be positive")
)

func NewWeightValue(value float64) (WeightValue, error) {
	if value <= 0 {
		return 0, ErrWeightInvalid
	}
	
	if value < MinWeightValue {
		return 0, ErrWeightTooLow
	}
	
	if value > MaxWeightValue {
		return 0, ErrWeightTooHigh
	}
	
	return WeightValue(value), nil
}

func (w WeightValue) Float64() float64 {
	return float64(w)
}

func (w WeightValue) String() string {
	return fmt.Sprintf("%.1f", float64(w))
}

func (w WeightValue) IsZero() bool {
	return w == 0
}

func (w WeightValue) Subtract(other WeightValue) WeightValue {
	return WeightValue(float64(w) - float64(other))
}

func (w WeightValue) Add(other WeightValue) (WeightValue, error) {
	result := float64(w) + float64(other)
	return NewWeightValue(result)
}