package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 数据库中存储的密码哈希
	storedHash := "$2a$10$/C/H3O/n/ekNcOZZdyFIpezN6R3B3C0z8YNsgBUhawzCCNU3kbQDm"
	
	// 用户输入的密码
	inputPassword := "123456"
	
	// 验证密码
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(inputPassword))
	if err != nil {
		fmt.Printf("密码验证失败: %v\n", err)
	} else {
		fmt.Println("密码验证成功")
	}
	
	// 生成新的密码哈希用于比较
	newHash, err := bcrypt.GenerateFromPassword([]byte(inputPassword), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("生成新哈希失败: %v\n", err)
	} else {
		fmt.Printf("新生成的密码哈希: %s\n", string(newHash))
	}
} 