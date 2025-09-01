package weight

import "errors"

type WeightUnit string

const (
	WeightUnitKg WeightUnit = "kg"
	WeightUnitLb WeightUnit = "lb"
)

var (
	ErrInvalidWeightUnit = errors.New("invalid weight unit")
)

func NewWeightUnit(value string) (WeightUnit, error) {
	unit := WeightUnit(value)
	if !unit.IsValid() {
		return "", ErrInvalidWeightUnit
	}
	return unit, nil
}

func (w WeightUnit) String() string {
	return string(w)
}

func (w WeightUnit) IsValid() bool {
	switch w {
	case WeightUnitKg, WeightUnitLb:
		return true
	default:
		return false
	}
}