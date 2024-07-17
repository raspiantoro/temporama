package resp

import (
	"fmt"
	"strconv"
)

type ValueNodeType string

const (
	ValueNodeTypeBulkString   ValueNodeType = "$"
	ValueNodeTypeArray        ValueNodeType = "*"
	ValueNodeTypeSimpleString ValueNodeType = "+"
	ValueNodeTypeSimpleError  ValueNodeType = "-"
	ValueNodeTypeIntegers     ValueNodeType = ":"
)

type ValueNode struct {
	types ValueNodeType
	val   string
	nodes []ValueNode
}

func WithValue(val string) func(node *ValueNode) {
	return func(node *ValueNode) {
		node.val = val
	}
}

func NewValueNode(types ValueNodeType, opts ...func(node *ValueNode)) ValueNode {
	node := ValueNode{types: types}

	for _, opt := range opts {
		opt(&node)
	}

	return node
}

func (v *ValueNode) Value(val string) {
	v.val = val
}

func (v *ValueNode) Append(node ValueNode) {
	v.nodes = append(v.nodes, node)
}

func (v *ValueNode) Marshal() []byte {
	result := []byte{}

	switch v.types {
	case ValueNodeTypeBulkString:
		result = v.marshalBulkString()
	case ValueNodeTypeArray:
		result = v.marshalArray()
	case ValueNodeTypeSimpleString:
		result = v.marshalSimpleString()
	case ValueNodeTypeIntegers:
		result = v.marshalIntegers()
	case ValueNodeTypeSimpleError:
		result = v.marshalSimpleError()
	}

	return result
}

func (v *ValueNode) marshalSimpleError() []byte {
	result := []byte{}
	result = append(result, []byte(fmt.Sprintf("%s%s\r\n", v.types, v.val))...)
	return result
}

func (v *ValueNode) marshalIntegers() []byte {
	result := []byte{}
	val, err := strconv.ParseInt(v.val, 10, 64)
	if err != nil {
		// write simple error
		result = append(result, []byte("-ERROR: value is not integers\r\n")...)
		return result
	}

	result = append(result, []byte(fmt.Sprintf("%s%d\r\n", v.types, val))...)
	return result
}

func (v *ValueNode) marshalSimpleString() []byte {
	result := []byte{}
	result = append(result, []byte(fmt.Sprintf("%s%s\r\n", v.types, v.val))...)
	return result
}

func (v *ValueNode) marshalBulkString() []byte {
	result := []byte{}

	result = append(result, []byte(v.types)...)

	// -1 indicates the value is nil
	if v.val == "-1" {
		result = append(result, []byte(v.val)...)
	} else {
		result = append(result, []byte(fmt.Sprintf("%d\r\n%s", len(v.val), v.val))...)
	}

	result = append(result, []byte("\r\n")...)

	return result
}

func (v *ValueNode) marshalArray() []byte {
	result := []byte{}

	result = append(result, []byte(fmt.Sprintf("%s%d\r\n", v.types, len(v.nodes)))...)

	for _, node := range v.nodes {
		result = append(result, node.Marshal()...)
	}

	return result
}
