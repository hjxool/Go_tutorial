package greetings

import (
	"errors"
	"fmt"
	"math/rand"
)

func Greeting(name string) (string, error) {
	if name == "" {
		return "", errors.New("空名称")
	}
	message := fmt.Sprintf(randomStr(), name)
	return message, nil
}

// 当包增加新功能时不要删除原有的 而是添加新方法
func Greetings(names []string) (map[string]string, error) {
	messages := make(map[string]string) // 更常用make创建因为可以初始化参数及容量
	for _, name := range names {
		message, err := Greeting(name)
		if err != nil {
			return nil, err
		}
		messages[name] = message
	}
	return messages, nil
}

func randomStr() string {
	strs := []string{
		"你好，%v，欢迎",
		"很高心见到你，%v",
		"hi，%v，很高心遇见你",
	}
	return strs[rand.Intn(len(strs))]
}
