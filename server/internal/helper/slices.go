package helper

import "math/rand"

func randomIndex[T any](ts []T) int {
	if len(ts) < 1 {
		return 0
	}

	return rand.Intn(len(ts))
}

func RandomItem[T any](ts []T) T {
	return ts[randomIndex(ts)]
}

func RandomSubset[T any](ts []T, n int) []T {
	if n > len(ts) || n < 1 {
		return nil
	}

	indexes := make(map[int]bool)
	output := make([]T, 0, n)

	for len(output) < n {
		i := randomIndex(ts)
		_, ok := indexes[i]
		if ok {
			continue
		}

		output = append(output, ts[i])
		indexes[i] = true
	}

	return output
}
