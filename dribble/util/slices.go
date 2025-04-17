package util

type Pair[F, S any] struct {
	First  F
	Second S
}

func Zip[F, S any](first []F, second []S) []Pair[F, S] {
	if len(first) != len(second) {
		panic("slices must have the same length")
	}

	result := make([]Pair[F, S], len(first))
	for i := range first {
		result[i] = Pair[F, S]{first[i], second[i]}
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
