package main

import (
	"fmt"
	"greetings"
	"log"
)

func main() {
	// message, err := greetings.Greeting("world")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(message)

	names := []string{"张三", "李四", "王五"}
	messages, err := greetings.Greetings(names)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(messages)
}
