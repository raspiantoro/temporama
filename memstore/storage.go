package memstore

import (
	"errors"
)

var (
	ErrNilEntries = errors.New("nil entries")
)

type Storage struct {
	entries map[uint32]EntryNode
}

func (s Storage) Delete(hashKey uint32, key string) {
	var entry EntryNode

	entry, ok := s.entries[hashKey]
	if !ok {
		return
	}

	newEntry := entry.Delete(key, nil)
	if newEntry == nil {
		delete(s.entries, hashKey)
	} else {
		s.entries[hashKey] = *newEntry
	}
}

func (s Storage) Set(valueType ValueType, hashKey uint32, key string, args ...string) (newFieldNum int, err error) {
	entry := s.entries[hashKey]

	newFieldNum, err = entry.Set(valueType, key, args...)
	if err != nil {
		return 0, err
	}

	s.entries[hashKey] = entry

	return
}

func (s Storage) Get(valueType ValueType, hashKey uint32, key string, args ...string) (any, error) {
	entry, ok := s.entries[hashKey]
	if !ok {
		return nil, ErrNilEntries
	}

	return entry.Get(valueType, key, args...)
}
