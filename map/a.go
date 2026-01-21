package main

import "fmt"

func main() {
	m1 := make(map[string]int)
	m1["a"] = 1
	m1["b"] = 2
	fmt.Printf("%v\n", m1)
	for k, v := range m1 {
		println(k, v)
	}

	m2 := map[string]int{
		"a": 3,
		"b": 4,
	}
	fmt.Printf("%v\n", m2)
	delete(m2, "a")
	fmt.Printf("%v\n", m2)

	var m3 map[string]int //此时的默认值是nil,此时没有分配存储空间，所以不能执行赋值
	fmt.Printf("%v", m3)

}
