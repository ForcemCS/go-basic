.PHONY: all proto run clean

# 定义 protoc 命令的快捷方式
PROTOC = protoc

# 定义 Go 插件和 proto 文件的路径
PROTOC_GEN_GO = protoc-gen-go
PROTO_FILES = $(wildcard proto/*.proto)

all: proto run

# proto: 编译 .proto 文件生成 Go 代码
# --proto_path: 指定 .proto 文件的搜索路径
# --go_out: 指定生成 Go 文件的输出目录
# --go_opt=paths=source_relative: 这是一个重要选项，让生成的文件路径与源文件保持一致，而不是生成冗长的目录结构
proto:
	@echo "Generating Go code from Protobuf definitions..."
	$(PROTOC) --proto_path=proto --go_out=. --go_opt=paths=source_relative $(PROTO_FILES)
	@echo "Done."

# run: 运行主程序
run:
	@echo "Running the main application..."
	go run cmd/main.go

# clean: 清理生成的文件
clean:
	@echo "Cleaning up generated files..."
	rm -f internal/user/*.pb.go