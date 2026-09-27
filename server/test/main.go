package main

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// 在 Handler 中
func GetSessionCookie(c *gin.Context) {
	cookie, err := c.Cookie("session")
	if err != nil {
		c.JSON(400, gin.H{"error": "未找到 Cookie"})
		return
	}

	fmt.Println("Cookie Value:", cookie)
	// MTU5MjM4NzY1NHxEdi1CQkFFQ0IwQU1U...
}

func DecryptSession(cookieValue string, secretKey []byte) (map[interface{}]interface{}, error) {
	// 1. 创建 CookieStore
	store := sessions.NewCookieStore(secretKey)

	// 2. 解码 Cookie 值
	data, err := base64.URLEncoding.DecodeString(cookieValue)
	if err != nil {
		return nil, fmt.Errorf("base64 decode failed: %v", err)
	}

	// 3. 解析 Session
	session := sessions.NewSession(store, "session")
	session.ID = string(data[:32]) // Session ID 可能在前32字节

	// 4. 解密数据（需要访问内部方法）
	// 注意：gorilla/sessions 没有直接提供解密方法
	// 需要从 Cookie 中读取并解密

	return nil, fmt.Errorf("需要手动实现解密")
}
func main() {
	// 1. 从配置文件获取密钥
	secretKey := []byte("620WFwf0111")

	// 2. 从浏览器复制 Cookie 值
	cookieValue := "MTc4ODE2NjIwN3xEWDhFQVFMX2dBQUJFQUVRQUFEX2lmLUFBQU1HYzNSeWFXNW5EQTBBQzJWNGNHbHlaVjkwYVcxbEJXbHVkRFkwQkFZQV9OVXFndFlHYzNSeWFXNW5EQk1BRVhabGNtbG1hV05oZEdsdmJsOWpiMlJsQm5OMGNtbHVad3dJQUFZek5EWTNNakFHYzNSeWFXNW5EQWNBQldWdFlXbHNCbk4wY21sdVp3d2JBQmwzWVc1bmFtbHVjblZwUUhKbGJHbGhZMmhwYm1FdVkyOXR8zIjrBKgXSROspqbDuKgHzjmfQ5Zx9iZRDcIyb27EY5Q="

	// 3. 解密
	data, err := DecryptSession(cookieValue, secretKey)
	if err != nil {
		log.Fatal(err)
	}

	// 4. 查看结果
	fmt.Printf("Session 数据:\n")
	for key, value := range data {
		fmt.Printf("  %s: %v (类型: %T)\n", key, value, value)
	}
}
