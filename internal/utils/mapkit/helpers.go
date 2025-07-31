package mapkit

import (
	"cmp"
	"slices"
)

func SortedKeys[K cmp.Ordered, V any](m map[K]V) []K {
	keys := Keys(m)
	slices.Sort(keys)
	return keys
}

func SortedValues[K cmp.Ordered, V any](m map[K]V) []V {
	keys := SortedKeys(m)
	values := make([]V, len(keys))
	for i, key := range keys {
		values[i] = m[key]
	}
	return values
}

func FirstKey[K cmp.Ordered, V any](m map[K]V) (firstKey K, ok bool) {
	for key, _ := range m {
		if !ok {
			firstKey, ok = key, true
		} else if key < firstKey {
			firstKey = key
		}
	}
	return
}
