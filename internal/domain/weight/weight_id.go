package weight

import (
	"errors"
	"strings"
)

type WeightID string

var (
	ErrEmptyWeightID = errors.New("weight ID cannot be empty")
)

func NewWeightID(value string) (WeightID, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", ErrEmptyWeightID
	}

	return WeightID(trimmed), nil
}

func (w WeightID) String() string {
	return string(w)
}

func (w WeightID) IsEmpty() bool {
	return string(w) == ""
}
