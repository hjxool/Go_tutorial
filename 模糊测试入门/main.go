package main

import "fmt"

func Reverse(s string) string {
	b := []byte(s) // 将string转换为字节切片类型 将每一个字符放入切片
	for i, j := 0, len(b)-1; i < len(b)/2; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b) // 再转换回字符串
}

func main() {
	input := "The quick brown fox jumped over the lazy dog"
	rev := Reverse(input)
	doubleRev := Reverse(rev)
	fmt.Printf("原字符串: %q\n", input)
	fmt.Printf("第一次转换: %q\n", rev)
	fmt.Printf("第二次转换: %q\n", doubleRev)
}
