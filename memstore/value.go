package memstore

type ValueType string

const (
	ValueTypeString ValueType = "string"
	ValueTypeMap    ValueType = "map"
)

type valueType[T any] ValueType

const (
	valueTypeString valueType[valueString] = valueType[valueString](ValueTypeString)
	valueTypeMap    valueType[valueMap]    = valueType[valueMap](ValueTypeMap)
)

type Value[T any] struct {
	types valueType[T]
	value[T]
}

type value[T any] struct {
	internal T
}

type valueString struct {
	val string
}

func (v *valueString) get() string {
	return v.val
}

func (v *valueString) set(val string) {
	v.val = val
}

type valueMap struct {
	val map[string]string
}

func (v *valueMap) get(keys ...string) []string {
	ret := []string{}

	for _, key := range keys {
		v, ok := v.val[key]
		if !ok {
			if key == "" {
				continue
			}

			v = "-1"
		}
		ret = append(ret, v)
	}

	return ret
}

func (v *valueMap) set(args ...string) {
	for i := 0; i < len(args); i += 2 {
		if args[i] == "" {
			continue
		}
		v.val[args[i]] = args[i+1]
	}
}
