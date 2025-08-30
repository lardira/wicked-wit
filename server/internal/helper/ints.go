package helper

import "math"

type intNumber interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func MinInt[T intNumber](a T, b T) (min T) {
	return T(math.Min(
		float64(a),
		float64(b),
	))
}
