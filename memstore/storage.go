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

func Set(valueType ValueType, key string, args ...string) error {
	return storage.Set(valueType, key, args...)
}

type Storage struct {
	entries map[string]any
}

func (s Storage) validate(valueType ValueType, key string) error {
	container, ok := s.entries[key]
	if !ok {
		return nil
	}

	var valid bool
	switch valueType {
	case ValueTypeString:
		_, valid = container.(Value[valueString])
	case ValueTypeMap:
		_, valid = container.(Value[valueMap])
	default:
		return errors.New("unknown value type")
	}

	if !valid {
		return errors.New("invalid operations for the value type")
	}

	return nil
}

func (s Storage) Set(valueType ValueType, key string, args ...string) error {
	err := s.validate(valueType, key)
	if err != nil {
		return err
	}

	switch valueType {
	case ValueTypeString:
		s.entries[key] = s.setString(args[0])
	case ValueTypeMap:
		s.entries[key] = s.setMap(args...)
	}

	return nil
}

func (s Storage) setString(arg string) Value[valueString] {
	val := Value[valueString]{
		types: valueTypeString,
		value: value[valueString]{
			internal: valueString{},
		},
	}

	val.internal.set(arg)

	return val
}

func (s Storage) setMap(args ...string) Value[valueMap] {
	val := Value[valueMap]{
		types: valueTypeMap,
		value: value[valueMap]{
			internal: valueMap{
				val: make(map[string]string),
			},
		},
	}

	val.internal.set(args...)

	return val
}

func (s Storage) Get(valueType ValueType, key string, args ...string) (any, error) {
	err := s.validate(valueType, key)
	if err != nil {
		return nil, err
	}

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
	valContainer := container.(Value[valueString])
	val := valContainer.internal.get()
	return val, nil
}

func (s Storage) getMap(key string, args ...string) ([]string, error) {
	container, ok := s.entries[key]
	if !ok {
		return nil, ErrNilEntries
	}

	valContainer := container.(Value[valueMap])
	val := valContainer.internal.get(args...)
	return val, nil
}
