package readfile

import "os"

// 读取文件
func ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}
