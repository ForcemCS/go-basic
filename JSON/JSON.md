## json.Marshal

`json.Marshal` 是将 Go 语言中的**结构化对象**（如 `struct`、`map`）转化为**字节序列（JSON 字符串）**的过程。

### 1. 深入理解 `json.Marshal`

想象一下，你在 Go 代码里有一个整齐的房间（结构体），里面放着各种家具（字段）。如果你想把这个房间“寄送”给网络另一端的浏览器，你不能直接搬房间，你需要把它“拆解并打包”成一个扁平的纸箱（JSON 字节流）。

- **输入：** Go 的变量（结构体、切片、映射等）。
- **输出：** `[]byte`（符合 JSON 规范的字节数组）和 `error`。

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Config struct {
	DatabasePath string `json:"db_path"`
	SecretKey    string `json:"secret_key,omitempty`
	MaxConns     int    `json:"max_connections"`

	Metadata struct {
		CreatedAt time.Time `json:"created_at"`
		Version   string    `json:"version"`
	} `json:"metadata"`
}

func main() {
	// 这里的原始数据通常是外部输入的
	inputData := []byte(`{
		"db_path": "C:\\Users\\ForceCS\\Desktop\\go_project\\duckdb\\data.db",
		"max_connections": 100,
		"metadata": {
			"version": "v1.0.0"
		}
	}`)

	var cfg Config
	// 将字节流转换为 Go 程序可以操作的结构体
	if err := json.Unmarshal(inputData, &cfg); err != nil {
		log.Fatal("解析配置失败: %v", err)

	}

	fmt.Println("成功加载配置。数据库路径: %s\n", cfg.DatabasePath)
	// 生产实践：在这里进行逻辑处理
	fmt.Printf("成功加载配置。数据库路径: %s\n", cfg.DatabasePath)
	cfg.MaxConns = 200 // 动态修改配置
	cfg.Metadata.CreatedAt = time.Now()
	cfg.SecretKey = "super-secret-token" // 之前为空，现在赋值

	// Marshal (序列化) ---
	// 目的：将修改后的对象转换回字节流，以便保存到文件或发送到其他系统

	// 在生产调试或保存配置文件时，通常使用 MarshalIndent
	// 它会带缩进，方便人类阅读
	outputBytes, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		log.Fatalf("生成配置失败: %v", err)
	}
	// 将结果输出（模拟保存到文件）
	fmt.Println("\n--- 生成的新配置文件内容 ---")
	fmt.Println(string(outputBytes))

	// 生产中的实际保存逻辑：
	// os.WriteFile("config.json", outputBytes, 0644)
}

```

