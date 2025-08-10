package workflow

import (
	"iter"
	"sync"
)

type Map[K comparable, V any] struct {
	mutex  sync.RWMutex
	values map[K]V
	parent *Map[K, V]
}

func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		values: make(map[K]V),
	}
}

func MapFrom[K comparable, V any](values map[K]V) *Map[K, V] {
	return &Map[K, V]{
		values: values,
	}
}

func (m *Map[K, V]) Set(k K, v V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.values[k] = v
}

func (m *Map[K, V]) Get(k K) (V, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	v, ok := m.values[k]
	for !ok && m.parent != nil {
		v, ok = m.parent.Get(k)
	}

	return v, ok
}

func (m *Map[K, V]) Clone() *Map[K, V] {
	return &Map[K, V]{
		values: make(map[K]V),
		parent: m,
	}
}

func (m *Map[K, V]) Iter() iter.Seq2[K, V] {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return func(yield func(k K, v V) bool) {
		for k, v := range m.values {
			if !yield(k, v) {
				return
			}
		}

		if m.parent != nil {
			for k, v := range m.parent.Iter() {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}
