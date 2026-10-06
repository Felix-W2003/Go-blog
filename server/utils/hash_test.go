package utils

import "testing"

// TestJtiSHA256 验证黑名单 key 的摘要函数。
//
// 这个函数的输出直接决定「写入黑名单」和「查询黑名单」能否对上：
// 一旦算法被改动（例如有人换成 MD5），下面的 SHA256 标准测试向量会立刻失败。
func TestJtiSHA256(t *testing.T) {
	t.Run("确定性", func(t *testing.T) {
		if JtiSHA256("abc") != JtiSHA256("abc") {
			t.Error("相同输入应得到相同摘要，否则黑名单的写入与查询会错位")
		}
	})

	t.Run("输出为64个十六进制字符", func(t *testing.T) {
		if n := len(JtiSHA256("abc")); n != 64 {
			t.Errorf("摘要长度 = %d, 期望 64", n)
		}
	})

	t.Run("符合SHA256标准测试向量", func(t *testing.T) {
		const want = "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"
		if got := JtiSHA256("abc"); got != want {
			t.Errorf("JtiSHA256(\"abc\") = %s, 期望 %s", got, want)
		}
	})

	t.Run("空串也会得到固定摘要", func(t *testing.T) {
		// 空 jti 的摘要是固定常量，这意味着「空 jti」会映射到同一个黑名单 key。
		// 调用方应避免拿空 jti 去写黑名单（JoinInBlacklist 的入参来自 Redis，正常不会为空）。
		const want = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
		if got := JtiSHA256(""); got != want {
			t.Errorf("JtiSHA256(\"\") = %s, 期望 %s", got, want)
		}
	})

	t.Run("不同输入不碰撞", func(t *testing.T) {
		if JtiSHA256("abc") == JtiSHA256("abd") {
			t.Error("不同输入不应得到相同摘要")
		}
	})
}

// TestBcrypt 验证密码哈希与校验的往返，以及异常输入的处理。
func TestBcrypt(t *testing.T) {
	const password = "P@ssw0rd-测试"

	hash := BcryptHash(password)
	if hash == "" {
		t.Fatal("BcryptHash 返回了空串")
	}
	if hash == password {
		t.Fatal("哈希值不应等于明文密码")
	}

	if !BcryptCheck(password, hash) {
		t.Error("正确密码应校验通过")
	}
	if BcryptCheck("wrong-password", hash) {
		t.Error("错误密码不应校验通过")
	}
	if BcryptCheck(password, "not-a-bcrypt-hash") {
		t.Error("非法哈希值不应校验通过")
	}
}

// TestMD5V 钉住上传文件名的摘要算法。
//
// 注意：MD5 在本项目里只用于生成上传文件名（utils/upload/local.go、qiniu.go），
// 属于非安全用途；密码等安全场景请使用 BcryptHash。
func TestMD5V(t *testing.T) {
	const want = "900150983cd24fb0d6963f7d28e17f72" // md5("abc") 的标准值
	if got := MD5V([]byte("abc")); got != want {
		t.Errorf("MD5V(\"abc\") = %s, 期望 %s", got, want)
	}
}
