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
	var entry EntryNode

	entry = s.entries[hashKey]
	// if !ok {
	// 	entry = EntryNode{
	// 		key: key,
	// 	}
	// 	s.entries[hashKey] = entry
	// }

	newFieldNum, err = entry.Set(valueType, key, args...)
	if err != nil {
		return 0, err
	}

	s.entries[hashKey] = entry

	return
}

// func (s Storage) setString(key, arg string) error {
// 	var val valueString
// 	_, ok := s.entries[key]
// 	if ok {
// 		val, ok = s.entries[key].(valueString)
// 		if !ok {
// 			return errors.New("invalid operations for the value type")
// 		}

// 		val.set(arg)
// 		s.entries[key] = val

// 		return nil
// 	}

// 	val = valueString{
// 		val: arg,
// 	}

// 	s.entries[key] = val
// 	return nil
// }

// func (s Storage) setMap(key string, args ...string) (int, error) {
// 	var (
// 		newFieldNum int
// 		val         valueMap
// 	)

// 	_, ok := s.entries[key]
// 	if ok {
// 		val, ok = s.entries[key].(valueMap)
// 		if !ok {
// 			return newFieldNum, errors.New("invalid operations for the value type")
// 		}

// 		newFieldNum = val.set(args...)
// 		s.entries[key] = val

// 		return newFieldNum, nil
// 	}

// 	val = valueMap{
// 		val: make(map[string]string),
// 	}

// 	newFieldNum = val.set(args...)
// 	s.entries[key] = val

// 	return newFieldNum, nil
// }

func (s Storage) Get(valueType ValueType, hashKey uint32, key string, args ...string) (any, error) {
	// switch valueType {
	// case ValueTypeString:
	// 	return s.getString(key)
	// case ValueTypeMap:
	// 	return s.getMap(key, args...)
	// }

	var entry EntryNode

	entry, ok := s.entries[hashKey]
	if !ok {
		return nil, ErrNilEntries
	}

	return entry.Get(valueType, key, args...)
}

// func (s Storage) getString(key string) (string, error) {
// 	container, ok := s.entries[key]
// 	if !ok {
// 		return "", ErrNilEntries
// 	}

// 	valContainer, ok := container.(valueString)
// 	if !ok {
// 		return "", errors.New("invalid operations for the value type")
// 	}

// 	val := valContainer.get()
// 	return val, nil
// }

// func (s Storage) getMap(key string, args ...string) ([]string, error) {
// 	container, ok := s.entries[key]
// 	if !ok {
// 		return nil, ErrNilEntries
// 	}

// 	valContainer, ok := container.(valueMap)
// 	if !ok {
// 		return nil, errors.New("invalid operations for the value type")
// 	}

// 	val := valContainer.get(args...)
// 	return val, nil
// }
