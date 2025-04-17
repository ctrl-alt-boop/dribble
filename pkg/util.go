package shared

func SliceTransform[T any, U any](slice []T, selector func(T) U) []U {
	results := make([]U, len(slice))
	for i, value := range slice {
		results[i] = selector(value)
	}
	return results
}
