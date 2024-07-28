package memstore

type ValueType int

const (
	ValueTypeString ValueType = iota
	ValueTypeMap
)

type valueString struct {
	val string
}

func (v *valueString) get() string {
	return v.val
}

func (v *valueString) set(arg string) {
	v.val = arg
}

type valueMap struct {
	val map[string]string
}

func (v *valueMap) get(keys ...string) []string {
	if len(keys) <= 0 {
		return v.getAll()
	}

	ret := []string{}

	for _, key := range keys {
		val, ok := v.val[key]
		if !ok {
			if key == "" {
				continue
			}

			val = "-1"
		}
		ret = append(ret, val)
	}

	return ret
}

func (v *valueMap) getAll() []string {
	ret := []string{}

	for key, val := range v.val {
		ret = append(ret, key)
		ret = append(ret, val)
	}

	return ret
}

func (v *valueMap) set(args ...string) int {
	newFieldNum := 0

	for i := 0; i < len(args); i += 2 {
		if args[i] == "" {
			continue
		}

		if _, ok := v.val[args[i]]; !ok {
			newFieldNum++
		}

		v.val[args[i]] = args[i+1]
	}

	return newFieldNum
}
