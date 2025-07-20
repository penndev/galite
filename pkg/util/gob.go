package util

import (
	"bytes"
	"encoding/gob"
)

// 将结构体编译为gob字节流
func Encode(data any) (bytes.Buffer, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		return bytes.Buffer{}, err
	}
	return buf, nil
}

// 从字节流恢复结构体数据
//
// - 参数
//   - []byte bytes.NewBuffer
//   - string bytes.NewBufferString
func Decode(buf *bytes.Buffer, target any) error {
	dec := gob.NewDecoder(buf)
	return dec.Decode(target)
}
