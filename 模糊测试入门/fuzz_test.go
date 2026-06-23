package main

import (
	"testing"
	"unicode/utf8"
)

// 单元测试得手动设计输入 模糊测试由api提供
// 单元测试以TestXxx开头 并采用*testing.T参数
func TestReverse(t *testing.T) {
	// struct{ in, want string } 是匿名结构体 等价于 struct{ in string, want string } 相同类型可以合并连写
	// {"!12345", "54321!"} 等价于 {in: "!12345", want: "54321!"} 因为go结构体字面量允许省略字段名 只要按字段顺序提供值即可
	testcases := []struct{ in, want string }{
		{"Hello, world", "dlrow ,olleH"},
		{" ", " "},
		{"!12345", "54321!"},
	}
	for _, v := range testcases {
		rev := Reverse(v.in)
		if rev != v.want {
			// t.Errorf 会自动在输出末尾添加换行 所以不需要 \n
			t.Errorf("转换值: %q, 期望: %q", rev, v.want)
		}
	}
}

// 模糊测试以FuzzXxx开头 并采用*testing.F参数
func FuzzReverse(f *testing.F) {
	testcases := []string{"Hello, world", " ", "!12345"}
	for _, v := range testcases {
		// 模糊测试会根据语料库自动变异输入 所要提供的初始样本类似黑盒测试等价类划分 但不需要覆盖所有类别 模糊测试会自动扩展
		f.Add(v) // 往语料库添加初始输入
	}

	// f.Fuzz 传入目标函数 会随着生成的变异输入不断执行
	// t *testing.T 后的参数必须与 f.Add 一一对应 如f.Add(10, "abc")则f.Fuzz(func(t *testing.T, n int, s string)
	// 如果f.Add("abc") f.Add(123)添加不同类型的语料 则报错
	f.Fuzz(func(t *testing.T, origin string) {
		rev := Reverse(origin)
		doubleRev := Reverse(rev)

		// 因为无法确认输入 因此不检查输出 而是检查结果是否符合某些特性 这是模糊测试的核心
		if origin != doubleRev {
			t.Errorf("之前: %q, 之后: %q", origin, doubleRev)
		}
		// 也会变异出中文
		if utf8.ValidString(origin) && !utf8.ValidString(rev) {
			// 原字符串是有效 UTF‑8 且 转换后不是有效 UTF‑8
			t.Errorf("转换为无效UTF-8字符串: %q", rev)
		}
	})
}

// 补充：基准测试以BenchmarkXxx开头 并采用*testing.B参数

// 如果只希望运行某个具体的测试 go test -run=FuzzReverse
// 注意：go test 会执行所有目录下所有测试 但对于模糊测试只执行到添加语料库阶段 不会执行f.Fuzz 加-fuzz=Fuzz标志才真正执行模糊测试
// 官方建议执行go test -fuzz=Fuzz前 先go test一下 以确认种子输入有没有问题
