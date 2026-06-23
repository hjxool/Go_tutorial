package main

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

func ReverseHasBug(s string) string {
	b := []byte(s) // 将string转换为字节切片类型 将每一个字符放入切片
	for i, j := 0, len(b)-1; i < len(b)/2; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i] // 类似JS中 [a, b] = [b, a]解构赋值
	}
	return string(b) // 再转换回字符串
}
func ReverseHasBug2(s string) string {
	fmt.Printf("输入: %q\n", s)
	// 因为存在汉字 因此不应按byte遍历 而是按rune
	r := []rune(s)
	fmt.Printf("[]rune后: %q\n", r)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func Reverse(s string) (string, error) {
	if !utf8.ValidString(s) {
		// 不是有效 UTF-8 的字符
		return s, errors.New("输入不是有效UTF-8字符")
	}
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r), nil
}

func main() {
	input := "The quick brown fox jumped over the lazy dog"
	rev, revErr := Reverse(input)
	doubleRev, doubleRevErr := Reverse(rev)
	fmt.Printf("原字符串: %q\n", input)
	fmt.Printf("第一次转换: %q, err: %v\n", rev, revErr)
	fmt.Printf("第二次转换: %q, err: %v\n", doubleRev, doubleRevErr)
}

// go test -run=FuzzReverse/2f... 运行特定语料库条目
