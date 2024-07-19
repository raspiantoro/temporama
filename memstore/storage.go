package memstore

import (
	"errors"
)

var (
	ErrNilEntries = errors.New("nil entries")
)

var (
	storage Storage
)

func init() {
	storage = Storage{
		entries: make(map[string]any),
	}
}

func Get(valueType ValueType, key string, args ...string) (any, error) {
	return storage.Get(valueType, key, args...)
}

func Set(valueType ValueType, key string, args ...string) (int, error) {
	return storage.Set(valueType, key, args...)
}

func Delete(key string) string {
	delete(storage.entries, key)
	return "1"
}

type Storage struct {
	entries map[string]any
}

func (s Storage) Set(valueType ValueType, key string, args ...string) (newFieldNum int, err error) {
	switch valueType {
	case ValueTypeString:
		err = s.setString(key, args[0])
	case ValueTypeMap:
		newFieldNum, err = s.setMap(key, args...)
	}

	return
}

func (s Storage) setString(key, arg string) error {
	var val valueString
	_, ok := s.entries[key]
	if ok {
		val, ok = s.entries[key].(valueString)
		if !ok {
			return errors.New("invalid operations for the value type")
		}

		val.set(arg)
		s.entries[key] = val

		return nil
	}

	val = valueString{
		val: arg,
	}

	s.entries[key] = val
	return nil
}

func (s Storage) setMap(key string, args ...string) (int, error) {
	var (
		newFieldNum int
		val         valueMap
	)

	_, ok := s.entries[key]
	if ok {
		val, ok = s.entries[key].(valueMap)
		if !ok {
			return newFieldNum, errors.New("invalid operations for the value type")
		}

		newFieldNum = val.set(args...)
		s.entries[key] = val

		return newFieldNum, nil
	}

	val = valueMap{
		val: make(map[string]string),
	}

	newFieldNum = val.set(args...)
	s.entries[key] = val

	return newFieldNum, nil
}

func (s Storage) Get(valueType ValueType, key string, args ...string) (any, error) {
	switch valueType {
	case ValueTypeString:
		return s.getString(key)
	case ValueTypeMap:
		return s.getMap(key, args...)
	}

	return nil, nil
}

func (s Storage) getString(key string) (string, error) {
	container, ok := s.entries[key]
	if !ok {
		return "", ErrNilEntries
	}

	valContainer, ok := container.(valueString)
	if !ok {
		return "", errors.New("invalid operations for the value type")
	}

	val := valContainer.get()
	return val, nil
}

func (s Storage) getMap(key string, args ...string) ([]string, error) {
	container, ok := s.entries[key]
	if !ok {
		return nil, ErrNilEntries
	}

	valContainer, ok := container.(valueMap)
	if !ok {
		return nil, errors.New("invalid operations for the value type")
	}

	val := valContainer.get(args...)
	return val, nil
}
