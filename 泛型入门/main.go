package main

import "fmt"

func main() {
	ints := map[string]int64{
		"first":  34,
		"second": 12,
	}
	floats := map[string]float64{
		"first":  35.98,
		"second": 26.99,
	}

	fmt.Printf("普通求和: %v 和 %v\n", SumInts(ints), SumFloats(floats))
	// 泛型示例
	fmt.Printf("泛型求和: %v 和 %v\n", SumIntsOrFloats[string, int64](ints), SumIntsOrFloats[string, float64](floats))
	// 也可以不传入类型参数 由编译器推断类型 但注意对于不接收参数的泛型函数 就需要传入类型参数
	fmt.Printf("自动推断类型的泛型求和: %v 和 %v\n", SumIntsOrFloats(ints), SumIntsOrFloats(floats))

	fmt.Printf("使用类型约束和泛型的求和: %v 和 %v\n",
		SumNumbers(ints),
		SumNumbers(floats))
}

func SumInts(m map[string]int64) int64 {
	var sum int64
	for _, v := range m {
		sum += v
	}
	return sum
}

func SumFloats(m map[string]float64) float64 {
	var sum float64
	for _, v := range m {
		sum += v
	}
	return sum
}

// comparable 不是具体类型 是类型约束 表示该类型的值必须支持 比较运算符
func SumIntsOrFloats[K comparable, V int64 | float64](m map[K]V) V {
	var sum V
	for _, v := range m {
		sum += v
	}
	return sum
}

// 类型约束
type Number interface {
	int64 | float64
}

func SumNumbers[K comparable, V Number](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}
