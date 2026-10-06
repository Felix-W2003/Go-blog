package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// TestParseDuration 覆盖 ParseDuration 的正常输入、边界输入与非法输入。
//
// 重点是把「0s / 0d 能解析成功，但返回的值是 0」这个易踩的行为钉住：
// go-redis 的 Set 只在过期时间 > 0 时才带 EX/PX 参数，传 0 表示「永不过期」。
// 所以调用方（如 JoinInBlacklist、SetRedisJWT）必须自己判断返回值是否为正数，
// 不能只判断 err 是否为 nil。
func TestParseDuration(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    time.Duration
		wantErr bool
	}{
		// ---- 正常输入 ----
		{name: "天", input: "7d", want: 7 * 24 * time.Hour},
		{name: "小时", input: "2h", want: 2 * time.Hour},
		{name: "分钟", input: "30m", want: 30 * time.Minute},
		{name: "秒", input: "45s", want: 45 * time.Second},
		{name: "组合", input: "1d2h30m", want: 26*time.Hour + 30*time.Minute},

		// ---- 边界：能解析成功，但值是 0，调用方必须自己拦住 ----
		{name: "零点", input: "0s", want: 0},
		{name: "零天", input: "0d", want: 0},

		// ---- 前后空白会被 TrimSpace 去掉 ----
		{name: "前后空格", input: "  7d  ", want: 7 * 24 * time.Hour},

		// ---- 非法输入：必须报错，绝不能静默返回 0 ----
		{name: "空串", input: "", wantErr: true},
		{name: "纯空格", input: "   ", wantErr: true},
		{name: "纯数字无单位", input: "7", wantErr: true},
		{name: "中文单位", input: "7天", wantErr: true},
		{name: "单位乱序", input: "1m1d", wantErr: true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := ParseDuration(c.input)

			if c.wantErr {
				if err == nil {
					t.Fatalf("ParseDuration(%q) 期望报错，实际返回 (%v, nil)", c.input, got)
				}
				return
			}

			if err != nil {
				t.Fatalf("ParseDuration(%q) 不该报错，实际 err = %v", c.input, err)
			}
			if got != c.want {
				t.Errorf("ParseDuration(%q) = %v, 期望 %v", c.input, got, c.want)
			}
		})
	}
}

// TestParseRefreshExp 验证「传入绝对过期时间，返回相对剩余时长」这一契约。
//
// 这个函数内部用 time.Parse 解析 time.Time.String() 的输出，
// layout 里同时出现 -0700 和 MST，时区偏移是最容易出错的地方，
// 因此用 55m ~ 65m 的容差兜住：一旦偏移算错（例如差了 8 小时）会立刻失败。
func TestParseRefreshExp(t *testing.T) {
	t.Run("未来一小时应返回约一小时", func(t *testing.T) {
		exp := jwt.NewNumericDate(time.Now().Add(time.Hour))

		if got := ParseRefreshExp(exp); got < 55*time.Minute || got > 65*time.Minute {
			t.Errorf("ParseRefreshExp(now+1h) = %v, 期望落在 55m ~ 65m 之间", got)
		}
	})

	t.Run("已经过期的时间应返回负数", func(t *testing.T) {
		exp := jwt.NewNumericDate(time.Now().Add(-time.Hour))

		if got := ParseRefreshExp(exp); got >= 0 {
			t.Errorf("ParseRefreshExp(now-1h) = %v, 期望为负数", got)
		}
	})
}
