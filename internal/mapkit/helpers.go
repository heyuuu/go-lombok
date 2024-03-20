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

func First[K cmp.Ordered, V any](m map[K]V) (key K, value V, ok bool) {
	if key, ok = FirstKey(m); ok {
		return key, m[key], true
	}
	return
}

func Each[K comparable, V any](m map[K]V, handler func(K, V)) {
	for k, v := range m {
		handler(k, v)
	}
}

func EachEx[K comparable, V any](m map[K]V, handler func(K, V) error) error {
	for k, v := range m {
		err := handler(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetValue[K comparable, V any](m map[K]V, key K) V {
	v, _ := Get(m, key)
	return v
}

func Get[K comparable, V any](m map[K]V, key K) (V, bool) {
	if m == nil {
		var temp V
		return temp, false
	} else {
		v, ok := m[key]
		return v, ok
	}
}

func Add[K comparable, V any](m *map[K]V, key K, value V) bool {
	if *m == nil {
		*m = map[K]V{key: value}
		return true
	} else if _, exists := (*m)[key]; exists {
		return false
	} else {
		(*m)[key] = value
		return true
	}
}
func Set[K comparable, V any](m *map[K]V, key K, value V) {
	if *m == nil {
		*m = map[K]V{key: value}
	} else {
		(*m)[key] = value
	}
}

func Clean[K comparable, V any](m map[K]V) {
	if m == nil {
		return
	}
	for k := range m {
		delete(m, k)
	}
}
