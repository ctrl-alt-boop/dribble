package util

type Pair[F, S any] struct {
	Left  F
	Right S
}

func Zip[F, S any](left []F, right []S) []Pair[F, S] {
	if len(left) != len(right) {
		panic("slices must have the same length")
	}

	result := make([]Pair[F, S], len(left))
	for i := range left {
		result[i] = Pair[F, S]{left[i], right[i]}
	}
	return result
}

func Sum(numbers ...int) int {
	sum := 0
	for _, number := range numbers {
		sum += number
	}
	return sum
}
