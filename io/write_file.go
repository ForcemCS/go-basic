package io

import (
	"bufio"
	"fmt"
	"os"
)

// 适合一次性读写大文件
func WriteFile() {
	//fout, err是if语句的局部变量
	//返回的是os.File,它的Read 或 Write 方法时，每一次调用都会触发一次系统调用（System Call）。
	//我们可以使用bufio,在内存中设置一个缓冲区，来减少系统调用的次数，从而提高文件读写的效率。
	if fout, err := os.OpenFile("data/verse.txt", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0o666); err != nil {
		fmt.Printf("Open file failed : %s\n", err.Error())
	} else {
		defer fout.Close() //没有发生err的时候才能调用close
		fout.WriteString("纳兰性德\n")
		fout.WriteString("明月多情应笑我")
		fout.WriteString("\n")
		fout.WriteString("笑我如今")

	}
}

func WriteFileWithBufio() {
	// 1. 打开文件（逻辑与之前一致）
	fout, err := os.OpenFile("data/verse.txt", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0o666)
	if err != nil {
		fmt.Printf("Open file failed : %s\n", err.Error())
		return // 发生错误直接返回
	}
	// 确保在函数结束时关闭文件句柄
	defer fout.Close()

	// 2. 创建一个带缓冲的 Writer
	// 默认缓冲区大小通常是 4096 字节 (4KB)
	writer := bufio.NewWriter(fout)

	// 3. 将内容写入缓冲区
	// 注意：此时数据只是写到了内存里，并没有真正写入硬盘的文件中
	writer.WriteString("纳兰性德\n")
	writer.WriteString("明月多情应笑我\n")
	writer.WriteString("笑我如今")

	// 4. 【非常关键】刷新缓冲区
	// 将内存缓冲区中剩下的所有数据强制推送到硬盘
	// 如果不调用 Flush，由于缓冲区可能还没满，最后一部分数据会丢失
	err = writer.Flush()
	if err != nil {
		fmt.Printf("Flush failed: %s\n", err.Error())
	} else {
		fmt.Println("文件写入成功（已缓冲）")
	}
}
