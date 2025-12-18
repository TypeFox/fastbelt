package regexp

import "sort"

/**
 * BinarySearch_NextState performs a binary search on the lookup table to find the next state for the given rune.
 * @param r The input rune to search for.
 * @param lookup A slice of runes representing the start and end of ranges ({start0, end0, start1, end1, ...}).
 * @param next A slice of integers representing the next states corresponding to the ranges in lookup.
 * @return The next state if found, otherwise -1.
 */
func BinarySearch_NextState(r rune, lookup []rune, next []int) int {
	searchIndex := sort.Search(len(next), func(i int) bool {
		return lookup[i*2] > r
	}) - 1
	if searchIndex > -1 && lookup[searchIndex*2] <= r && r <= lookup[searchIndex*2+1] {
		return next[searchIndex]
	}
	return -1
}
