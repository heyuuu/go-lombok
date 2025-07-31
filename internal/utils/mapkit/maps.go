package mapkit

// common
type MapEntry[K comparable, V any] struct {
	Key   K
	Value V
}

// Map
type Map[K comparable, V any] interface {
	Exists(K) bool
	Get(K) (V, bool)
	GetOrDefault(K, V) V
	Set(K, V)
	Del(K)
	Keys() []K
	Values() []V
	Size() int
}

// StableMap 有序map，保持插入顺序; 零值可用
type StableMap[K comparable, V any] struct {
	keys    map[K]int
	entries []MapEntry[K, V]
}

var _ Map[string, any] = (*StableMap[string, any])(nil)

func (m *StableMap[K, V]) idx(key K) (idx int, ok bool) {
	if m.keys == nil {
		return 0, false
	}
	idx, ok = m.keys[key]
	return
}

func (m *StableMap[K, V]) Len(key K) int {
	return len(m.entries)
}

func (m *StableMap[K, V]) Exists(key K) bool {
	_, ok := m.idx(key)
	return ok
}

func (m *StableMap[K, V]) Get(key K) (V, bool) {
	if idx, ok := m.idx(key); ok {
		return m.entries[idx].Value, true
	} else {
		var tmp V
		return tmp, false
	}
}

func (m *StableMap[K, V]) GetOrDefault(key K, defaultValue V) V {
	if idx, ok := m.idx(key); ok {
		return m.entries[idx].Value
	}
	return defaultValue
}

func (m *StableMap[K, V]) Set(key K, value V) {
	if idx, ok := m.idx(key); ok {
		m.entries[idx].Value = value
	} else {
		newIdx := len(m.entries)
		if m.keys == nil {
			m.keys = make(map[K]int)
		}
		m.keys[key] = newIdx
		m.entries = append(m.entries, MapEntry[K, V]{Key: key, Value: value})
	}
}

func (m *StableMap[K, V]) Del(key K) {
	if idx, ok := m.idx(key); ok {
		delete(m.keys, key)

		if idx < len(m.entries)-1 {
			copy(m.entries[idx:], m.entries[idx+1:])
		}
		m.entries = m.entries[:len(m.entries)-1]
	}
}

func (m *StableMap[K, V]) Keys() []K {
	if len(m.entries) == 0 {
		return nil
	}

	keys := make([]K, len(m.entries))
	for idx, entry := range m.entries {
		keys[idx] = entry.Key
	}
	return keys
}

func (m *StableMap[K, V]) Values() []V {
	if len(m.entries) == 0 {
		return nil
	}

	values := make([]V, len(m.entries))
	for idx, entry := range m.entries {
		values[idx] = entry.Value
	}
	return values
}

func (m *StableMap[K, V]) Size() int {
	return len(m.entries)
}
