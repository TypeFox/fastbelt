package automatons

// Helper function for max
func max(a, b rune) rune {
	if a > b {
		return a
	}
	return b
}

// Helper function for min
func min(a, b rune) rune {
	if a < b {
		return a
	}
	return b
}

func CopyMap[K comparable, V any](m map[K]V) map[K]V {
	cp := make(map[K]V)
	for k, v := range m {
		cp[k] = v
	}
	return cp
}
