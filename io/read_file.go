package io

import (
	"fmt"
	"io"
	"strings"
)

func ReadFile() {
	// 1. 模拟一个数据源（实现了 io.Reader 接口）
	data := "The design philosophy of the Go language: simplicity, efficiency, and reuse."
	reader := strings.NewReader(data)

	// 2. 关键点：在循环外部创建一个“容器”（缓冲区）
	// 这里我们只给 8 个字节的空间，模拟多次读取的过程
	buf := make([]byte, 8)

	fmt.Printf("开始读取，初始缓冲区地址: %p\n", &buf[0])
	fmt.Println("------------------------------------")

	for {
		// 3. 将容器 buf 传给 Read 方法
		// Read 不会创建新切片，而是尝试填满你给它的 buf
		n, err := reader.Read(buf)

		if n > 0 {
			// 4. 处理读到的数据
			// 注意：必须使用 buf[:n]，因为 buf 后面的部分可能还留着上次读取的旧数据
			fmt.Printf("读取字节数: %d | 内容: [%s] | 当前缓冲区首地址: %p\n",
				n, string(buf[:n]), &buf[0])
		}

		// 5. 判断是否读完（EOF）或报错
		if err == io.EOF {
			fmt.Println("------------------------------------")
			fmt.Println("读取完毕：已到达文件末尾 (EOF)")
			break
		}
		if err != nil {
			fmt.Printf("读取出错: %v\n", err)
			break
		}
	}
}
