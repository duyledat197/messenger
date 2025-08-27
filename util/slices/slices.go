package slices

// Transform applies a transformer function to each item and returns the results.
func Transform[T any, S any](items []T, transformer func(T) S) []S {
	result := make([]S, len(items))
	for i, item := range items {
		result[i] = transformer(item)
	}
	return result
}

// Find returns a slice of elements that match the predicate, limited by the given limit.
func Find[T any](items []T, predicate func(T) bool, limit int) []T {
	var result []T
	for _, item := range items {
		if predicate(item) {
			result = append(result, item)
		}
		if len(result) == limit {
			break
		}
	}
	return result
}

// Take returns a new slice containing the first n elements of the input slice.
// If n is less than or equal to zero, an empty slice is returned.
// If n is greater than the length of the slice, the entire slice is returned.
//
// Example:
//
//	nums := []int{1, 2, 3, 4, 5}
//	firstTwo := Take(nums, 2) // firstTwo -> [1, 2]
func Take[T any](items []T, n int) []T {
	if n <= 0 {
		return []T{}
	}
	if n >= len(items) {
		return items
	}
	return items[:n]
}
