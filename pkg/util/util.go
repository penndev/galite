package util

import (
	"fmt"
	"os"
	"strconv"
)

type Number interface {
	int | int16 | int32 | int64 | uint | uint16 | uint32 | uint64
}

func StrConv[R Number](str string) (R, error) {
	var zero R
	switch any(zero).(type) {
	case int:
		n, err := strconv.Atoi(str)
		return any(n).(R), err
	case int16:
		n, err := strconv.ParseInt(str, 10, 16)
		return any(int16(n)).(R), err
	case int32:
		n, err := strconv.ParseInt(str, 10, 32)
		return any(int32(n)).(R), err
	case int64:
		n, err := strconv.ParseInt(str, 10, 64)
		return any(n).(R), err
	case uint:
		n, err := strconv.ParseUint(str, 10, 64)
		return any(uint(n)).(R), err
	case uint16:
		n, err := strconv.ParseUint(str, 10, 16)
		return any(uint16(n)).(R), err
	case uint32:
		n, err := strconv.ParseUint(str, 10, 32)
		return any(uint32(n)).(R), err
	case uint64:
		n, err := strconv.ParseUint(str, 10, 64)
		return any(n).(R), err
	default:
		return zero, fmt.Errorf("unsupported type")
	}
}

// 字符串切片批量转换为 []R
func StrArrConv[R Number](strs []string) ([]R, error) {
	out := make([]R, 0, len(strs))
	for _, s := range strs {
		v, err := StrConv[R](s)
		if err != nil {
			return nil, fmt.Errorf("转换失败 '%s': %w", s, err)
		}
		out = append(out, v)
	}
	return out, nil
}

func GetEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
