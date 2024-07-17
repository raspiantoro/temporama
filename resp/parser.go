package resp

import (
	"io"
	"log"
	"strconv"
	"strings"
)

func parse(sr io.Reader) ValueNode {
	prefix := make([]byte, 1)
	_, err := sr.Read(prefix)
	if err != nil {
		return ValueNode{}
	}

	switch ValueNodeType(prefix) {
	case ValueNodeTypeBulkString:
		return parseBulkString(sr)
	case ValueNodeTypeArray:
		return parseArray(sr)
	case "\r\n":
		return parse(sr)
	case "\r":
		return parse(sr)
	case "\n":
		return parse(sr)
	default:
		log.Print("unknown command prefix: ", string(prefix))
	}

	return ValueNode{}
}

func parseArray(sr io.Reader) ValueNode {
	length := parseInteger(sr)

	node := ValueNode{
		types: ValueNodeTypeArray,
		nodes: make([]ValueNode, length),
	}

	for i := int64(0); i < length; i++ {
		node.nodes[i] = parse(sr)
	}

	return node
}

func parseInteger(sr io.Reader) int64 {
	// read 3 byte, include \r\n
	length := make([]byte, 3)
	_, err := sr.Read(length)
	if err != nil {
		log.Println(err)
		return 0
	}

	trimLength := strings.Trim(string(length), "\r\n")
	trimLength = strings.Trim(trimLength, "\r")
	trimLength = strings.Trim(trimLength, "\n")

	ln, err := strconv.ParseInt(trimLength, 10, 64)
	if err != nil {
		log.Println("err while parse int: ", err)
		return 0
	}

	return ln
}

func parseBulkString(sr io.Reader) ValueNode {
	length := parseInteger(sr)
	buf := bufCreate(length)

	node := ValueNode{
		types: ValueNodeTypeBulkString,
	}

	_, err := sr.Read(buf)
	if err != nil {
		log.Println(err)
		return ValueNode{}
	}

	str := strings.Trim(string(buf), "\r\n")
	str = strings.Trim(str, "\r")
	str = strings.Trim(str, "\n")

	node.val = str

	return node
}

func bufCreate(length int64) []byte {
	// +2 include \r\n
	return make([]byte, length+2)
}
