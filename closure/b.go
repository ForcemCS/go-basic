package main

import (
	"fmt"
	"sync"
	"time"
)

// NewOrderNumberGenerator 是一个工厂（外层函数）
// 传入参数 prefix：比如 "ORDER" 或 "REFUND"
// 返回值：一个专门用来生成单号的函数
func NewOrderNumberGenerator(prefix string) func() string {
	// --- 下面这些，都是装进“记忆背包”里的私有物品 ---
	var seq int         // 记忆序列号
	var mu sync.Mutex   // 记忆一把专门保护这个序列号的锁
	// ----------------------------------------------

	// 返回真正的打工人（闭包函数）
	return func() string {
		// 1. 拿背包里的锁，锁住（防止并发下单时序号重复）
		mu.Lock()
		defer mu.Unlock()

		// 2. 序列号 +1
		seq++
		
		// 3. 拿到今天的日期
		dateStr := time.Now().Format("20060102")

		// 4. 拼装出最终的单号：前缀-日期-4位补齐的序号 (如 ORDER-20231024-0001)
		return fmt.Sprintf("%s-%s-%04d", prefix, dateStr, seq)
	}
}

func main() {
	// 1. 在程序启动时，我们初始化两个生成器（制造两个带着背包的打工人）
	// 生成订单号的打工人
	generateOrderNo := NewOrderNumberGenerator("ORDER")
	// 生成退款单号的打工人（和上面的完全独立，互不干扰）
	generateRefundNo := NewOrderNumberGenerator("REFUND")

	// 2. 模拟真实业务场景：不停地接单
	fmt.Println(generateOrderNo())  // 输出: ORDER-20231024-0001
	fmt.Println(generateOrderNo())  // 输出: ORDER-20231024-0002
	
	fmt.Println(generateRefundNo()) // 输出: REFUND-20231024-0001
	
	fmt.Println(generateOrderNo())  // 输出: ORDER-20231024-0003
}
