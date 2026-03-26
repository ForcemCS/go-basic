package main

import "fmt"

// 1. 老板函数（外层函数），它的返回值是一个“打工人（函数）”
func newCounterWorker() func() int {
	// 这就是那个高科技计数器，它是老板办公室里的东西（局部变量）
	count := 0 

	// 2. 老板招了一个打工人（匿名函数），并把他派出去
	return func() int {
		// 【魔法发生的地方】：打工人把老板办公室的 count 装进了自己的“记忆背包”里！
		count++       // 每次调用，就操作自己背包里的 count
		return count
	}
}

func main() {
	// 老板派出了 1 号打工人。此时老板函数执行完毕退出了。
	// 但是！1 号打工人的背包里，永远留存着那个 count！
	worker1 := newCounterWorker()

	fmt.Println(worker1()) // 输出: 1
	fmt.Println(worker1()) // 输出: 2
	fmt.Println(worker1()) // 输出: 3

	// 老板又派出了 2 号打工人。这是一个全新的人，背着一个全新的、从 0 开始的背包。
	worker2 := newCounterWorker()
	fmt.Println(worker2()) // 输出: 1
	
	// 1 号打工人的数据依然安全
	fmt.Println(worker1()) // 输出: 4 
}
