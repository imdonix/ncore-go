package ncore

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

var unitSize = map[string]float64{
	"B":   1,
	"KiB": 1024,
	"MiB": math.Pow(1024, 2),
	"GiB": math.Pow(1024, 3),
	"TiB": math.Pow(1024, 4),
}

type Size struct {
	bytes float64
	unit  string
}

func NewSize(s string) (Size, error) {
	parts := strings.Split(s, " ")
	if len(parts) != 2 {
		return Size{}, fmt.Errorf("invalid size format: %s", s)
	}
	val, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return Size{}, err
	}
	unit := parts[1]
	multiplier, ok := unitSize[unit]
	if !ok {
		return Size{}, fmt.Errorf("invalid unit: %s", unit)
	}
	return Size{bytes: val * multiplier, unit: unit}, nil
}

var units = []string{"B", "KiB", "MiB", "GiB", "TiB"}

func NewSizeFromBytes(bytes float64) Size {
	unit := "B"
	for _, u := range units {
		multiplier := unitSize[u]
		if bytes/multiplier <= 1000 {
			unit = u
			break
		}
	}

	if bytes/unitSize[unit] > 1000 {
		unit = "TiB"
	}

	return Size{bytes: bytes, unit: unit}
}

func (s Size) String() string {
	multiplier := unitSize[s.unit]
	return fmt.Sprintf("%.2f %s", s.bytes/multiplier, s.unit)
}

func (s Size) Bytes() int64 {
	return int64(s.bytes)
}

func (s Size) Add(other Size) Size {
	return NewSizeFromBytes(s.bytes + other.bytes)
}
