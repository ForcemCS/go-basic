package main

import (
	"fmt"
	"log"

	// 导入生成的 user 包，注意路径是 go.mod 中定义的模块名 + internal/user
	user "github.com/your-username/protobuf-learning-project/internal/user"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func main() {
	// --- 1. 创建一个复杂的 UserProfile 实例 ---
	// 这里我们会填充所有定义的字段，来演示各个特性
	originalProfile := &user.UserProfile{
		UserId:             "u-12345",
		FullName:           "Alex Smith",
		Status:             user.AccountStatus_STATUS_LATEST_ACTIVE, // 使用别名
		FollowerCount:      1500,
		ProfileVersion:     268435457,         // 一个大数，演示 fixed32
		LoginAttemptsDelta: -2,                // 演示 sint32 处理负数
		LastLoginTime:      timestamppb.Now(), // 使用 well-known type

		// 演示 oneof：我们选择设置 phone_number
		ContactMethod: &user.UserProfile_PhoneNumber{
			PhoneNumber: "+1-800-123-4567",
		},

		// 演示 map：添加一些自定义属性
		Attributes: map[string]string{
			"theme":    "dark",
			"language": "en-US",
		},

		// 演示 wrapper type：创建一个 StringValue 来包装 "Lex"
		Nickname: wrapperspb.String("Lex"),

		// 演示 optional：直接赋值
		Age: 30,
	}

	fmt.Println("----------- Original Go Struct -----------")
	fmt.Printf("%#v\n\n", originalProfile)

	// --- 2. 序列化 (Marshal) ---
	// 将 Go 结构体实例编码成二进制字节流
	binaryData, err := proto.Marshal(originalProfile)
	if err != nil {
		log.Fatalf("Failed to marshal profile: %v", err)
	}

	fmt.Println("----------- Serialization -----------")
	fmt.Printf("Successfully serialized profile into %d bytes of binary data.\n\n", len(binaryData))
	// fmt.Println(binaryData) // 可以取消注释查看二进制数据

	// --- 3. 反序列化 (Unmarshal) ---
	// 假设我们从网络或文件中收到了 binaryData
	// 现在将其解码回一个新的 Go 结构体实例
	deserializedProfile := &user.UserProfile{}
	if err := proto.Unmarshal(binaryData, deserializedProfile); err != nil {
		log.Fatalf("Failed to unmarshal profile: %v", err)
	}

	fmt.Println("----------- Deserialized Go Struct -----------")
	fmt.Printf("%#v\n\n", deserializedProfile)

	// --- 4. 验证和演示 ---
	// 检查反序列化后的数据是否与原始数据一致，并演示如何访问特定字段
	fmt.Println("----------- Verification & Demonstration -----------")

	// 演示如何安全地访问 oneof 字段
	switch c := deserializedProfile.GetContactMethod().(type) {
	case *user.UserProfile_Email:
		fmt.Printf("Contact method is Email: %s\n", c.Email)
	case *user.UserProfile_PhoneNumber:
		fmt.Printf("Contact method is Phone Number: %s\n", c.PhoneNumber)
	default:
		fmt.Println("Contact method is not set.")
	}

	// 演示如何访问 map 字段
	fmt.Printf("User's theme from attributes map: %s\n", deserializedProfile.GetAttributes()["theme"])

	// 演示如何检查 optional 和 wrapper 字段的存在性
	// optional 字段
	if deserializedProfile.GetAge() > 0 { // proto3 中 optional 字段会生成 GetXxx 方法
		fmt.Printf("User's age is set: %d\n", deserializedProfile.GetAge())
	} else {
		fmt.Println("User's age is not set.")
	}

	// wrapper 字段
	if deserializedProfile.GetNickname() != nil {
		fmt.Printf("User's nickname is set: %s\n", deserializedProfile.GetNickname().GetValue())
	} else {
		fmt.Println("User's nickname is not set.")
	}

	// 演示枚举别名
	fmt.Printf("Account status: %s (numeric value: %d)\n",
		deserializedProfile.GetStatus(), deserializedProfile.GetStatus().Number())
}
