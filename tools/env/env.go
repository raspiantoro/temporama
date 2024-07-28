package env

import (
	"errors"
	"os"
	"strconv"
)

func GetUint32(name string) (uint32, error) {
	s := os.Getenv(name)
	if s == "" {
		return 0, errors.New("empty value")
	}

	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return uint32(v), nil
}
