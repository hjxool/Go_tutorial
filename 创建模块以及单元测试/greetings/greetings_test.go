// _test.go _window.go是特殊后缀不会打包进package
package greetings

import (
	// "regexp"
	"strings"
	"testing"
)

// 必须要以Test开头命名 且后面的第一个字母还必须是大写 才能识别
// 必须要传入*testing.T/B参数
// 测试函数不能有任何返回值！ 否则go test 会报错或拒绝承认它是测试函数
func TestGreetingName(t *testing.T) {
	name := "测试"
	// regexp 还有Compile 但更常用MustCompile
	// Compile 如果正则语法写错了 会返回一个 err 程序不会崩溃 需要写代码去处理错误
	// MustCompile 如果正则语法写错了 会直接触发 panic 导致程序崩溃 为了省去写 if err != nil 的麻烦 所以习惯用 MustCompile
	// want := regexp.MustCompile(`\b` + name + `\b`) // 不能用双引号 否则会吞掉反斜杠 ``是原生字面量

	msg, err := Greeting(name)
	// if !want.MatchString(msg) || err != nil {
	// 	// 结果不包含name值 或 报错了
	// 	// %#q 中 #号 表示用 Go 语法的字面量形式输出字符串 如fmt.Printf("%#q", "Hello\nWorld") → "Hello\nWorld"
	// 	t.Errorf(`Greeting(测试) = %q, %v, want匹配到 %#q, nil`, msg, err, want)
	// }

	// \b 只能检查[a-zA-Z0-9_]和空格 因此输入中文无法通过验证 用 strings 库更高效
	if !strings.Contains(msg, name) || err != nil {
		t.Errorf(`Greeting(测试) = %q, %v, 期望包含 %q 且无错误`, msg, err, name)
	}
}

func TestGreetingEmpty(t *testing.T) {
	msg, err := Greeting("")
	if msg != "" || err == nil {
		// 结果不为空 或 没有返回err
		t.Errorf(`Greeting("") = %q, %v, 期望结果 "", error`, msg, err)
	}
}

// go test 执行测试
// go test -v 获得详细输出
