package memstore

import (
	"errors"
)

type EntryNode struct {
	// will hold unhashed key
	key string
	val any
	// child is needed when there is hash collision within the key
	// it will hold the value with the same hashkey
	child    *EntryNode
	assigned bool
}

func NewEntryNode(key string, val any) EntryNode {
	return EntryNode{
		key: key,
		val: val,
	}
}

func (e *EntryNode) Append(child *EntryNode) {
	e.child = child
}

func (e *EntryNode) Child() *EntryNode {
	return e.child
}

// need to iterate child where there is collision within the key
func (e *EntryNode) Next() (*EntryNode, error) {
	if e.child == nil {
		return nil, errors.New("empty child")
	}

	return e.child, nil
}

func (e *EntryNode) Delete(key string, parent *EntryNode) (entry *EntryNode) {
	if e.key == key {
		e.val = nil
		e.key = ""

		child := e.child

		e.child = nil

		if child != nil {
			if parent != nil {
				parent.child = child
				return parent
			}

			return child
		} else {
			return parent
		}
	}

	if e.child != nil {
		child := e.child
		node := child.Delete(key, e)
		return node
	}

	return nil
}

func (e *EntryNode) Set(valueType ValueType, key string, args ...string) (newFieldNum int, err error) {
	var node *EntryNode

	if !e.assigned {
		e.key = key
		e.assigned = true

		switch valueType {
		case ValueTypeString:
			err = e.setString(args[0])
		case ValueTypeMap:
			newFieldNum, err = e.setMap(args...)
		}

		return newFieldNum, err
	}

	if e.key == key {
		node = e
	} else {
		node, _ = e.Next()
		if node != nil {
			return node.Set(valueType, key, args...)
		}
	}

	if node == nil {
		node = new(EntryNode)
		node.key = key
		e.child = node
	}

	switch valueType {
	case ValueTypeString:
		err = node.setString(args[0])
	case ValueTypeMap:
		newFieldNum, err = node.setMap(args...)
	}

	return newFieldNum, err
}

func (e *EntryNode) setString(arg string) error {
	if e.val != nil {
		val, ok := e.val.(valueString)
		if !ok {
			return errors.New("invalid operations for the value type")
		}

		val.set(arg)
	}

	e.val = valueString{
		val: arg,
	}

	return nil
}

func (e *EntryNode) setMap(args ...string) (int, error) {
	var (
		newFieldNum int
	)

	if e.val != nil {
		val, ok := e.val.(valueMap)
		if !ok {
			return newFieldNum, errors.New("invalid operations for the value type")
		}

		newFieldNum = val.set(args...)

		e.val = val

		return newFieldNum, nil
	}

	val := valueMap{
		val: make(map[string]string),
	}

	newFieldNum = val.set(args...)

	e.val = val
	return newFieldNum, nil
}

func (e *EntryNode) Get(valueType ValueType, key string, args ...string) (any, error) {
	switch valueType {
	case ValueTypeString:
		return e.getString(key)
	case ValueTypeMap:
		return e.getMap(key, args...)
	}

	return nil, nil
}

func (e *EntryNode) getString(key string) (string, error) {
	container := e.val
	if container == nil {
		return "", ErrNilEntries
	}

	valContainer, ok := container.(valueString)
	if !ok {
		return "", errors.New("invalid operations for the value type")
	}

	val := valContainer.get()
	return val, nil
}

func (e *EntryNode) getMap(key string, args ...string) ([]string, error) {
	container := e.val
	if container == nil {
		return nil, ErrNilEntries
	}

	valContainer, ok := container.(valueMap)
	if !ok {
		return nil, errors.New("invalid operations for the value type")
	}

	val := valContainer.get(args...)
	return val, nil
}
