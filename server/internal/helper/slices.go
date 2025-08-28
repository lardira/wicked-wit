package helper

import "math/rand"

func RandomItem[T any](ts []T) T {
	return ts[rand.Intn(len(ts))]
}
