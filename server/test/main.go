package main

import (
	"encoding/json"
	"fmt"
)

// Sex 自定义int枚举
type Sex int

const (
	Male   Sex = iota // 0 男
	Female            // 1 女
)

// ---------------- ① MarshalJSON：Go → JSON 输出 ----------------
// 实现Marshaler接口，控制【输出成json长什么样】
func (s Sex) MarshalJSON() ([]byte, error) {
	var txt string
	switch s {
	case Male:
		txt = "男"
	case Female:
		txt = "女"
	default:
		txt = "未知"
	}
	// json.Marshal会给字符串加上双引号，变成json合法 "男"
	return json.Marshal(txt)
}

// ---------------- ② UnmarshalJSON：JSON → Go 读进来 ----------------
// 必须指针接收者！要修改变量
func (s *Sex) UnmarshalJSON(data []byte) error {
	// data是json原始字节，比如 []byte(`"男"`)，自带双引号
	var inputStr string
	// 把带双引号的json，解析成普通go字符串："男" → 男
	err := json.Unmarshal(data, &inputStr)
	if err != nil {
		return err
	}

	// 根据字符串转成数字枚举
	switch inputStr {
	case "男":
		*s = Male
	case "女":
		*s = Female
	default:
		// 非法值
		*s = -1
	}
	return nil // 成功必须return nil
}

func main() {
	// type User struct {
	// 	Name string `json:"name"`
	// 	Sex  Sex    `json:"sex"`
	// }

	// // ========== 测试1：序列化 Go对象转json字符串 ==========
	// u1 := User{
	// 	Name: "小明",
	// 	Sex:  Male, // 内部值是数字0
	// }
	// b, _ := json.Marshal(u1)
	// fmt.Println("序列化结果：", string(b))
	// // 如果没有自定义接口：{"name":"小明","sex":0}
	// // ✅自定义后输出：{"name":"小明","sex":"男"}

	// // ========== 测试2：反序列化 json字符串转Go对象 ==========
	// jsonText := `{"name":"小红","sex":"女"}`
	// var u2 User
	// _ = json.Unmarshal([]byte(jsonText), &u2)
	// fmt.Printf("反序列化结果：name=%s, sex内部数字=%d\n", u2.Name, u2.Sex)
	// // sex内部存的是数字1(Female)，但json写的是"女"

	var s1 = "127.0.0.1"
	var s2 = 8080
	var s3 = fmt.Sprint("%s:%d", s1, s2)
	fmt.Println(s3)
}
