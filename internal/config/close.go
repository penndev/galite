package config

import "log"

type closeInterface interface {
	Close() error
}

var closeList []closeInterface

// 执行程序清理
// defer关闭前执行的close操作
func Defer() {
	for _, r := range closeList {
		if err := r.Close(); err != nil {
			log.Printf("关闭资源出错: %v\n", err)
		}
	}
}
