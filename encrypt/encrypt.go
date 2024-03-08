package encrypt

import (
	`crypto/md5`
	`crypto/sha1`
	`crypto/sha256`
	`crypto/sha512`
	`encoding/base64`
	`fmt`
	
	`github.com/google/uuid`
)

// Md5 用于计算输入字符串的MD5值
// 参数：
//   value string - 需要计算MD5值的字符串
// 返回值：
//   string - 计算得到的MD5字符串
func Md5(value string) string {
	// 使用md5包计算字符串的哈希值
	h := md5.Sum([]byte(value))
	// 将哈希值转换为十六进制字符串
	md5String := fmt.Sprintf("%x", h)
	return md5String
}

// SHA1 通过SHA1算法计算给定字符串的哈希值，并以Base64编码格式返回
// 参数:
//   s string - 需要进行SHA1哈希计算的原始字符串
// 返回值:
//   p string - 计算得到的哈希值，以Base64编码格式呈现
func SHA1(s string) (p string) {
	// 创建SHA1哈希对象
	o := sha1.New()
	// 向哈希对象中写入字符串的字节序列
	o.Write([]byte(s))
	// 将哈希对象的最终结果编码为Base64字符串
	p = base64.StdEncoding.EncodeToString(o.Sum(nil))
	return
}

// SHA256 使用SHA256算法计算给定字符串的哈希值，并以Base64编码格式返回。
//
// 参数:
//   s string - 需要进行SHA256哈希计算的原始字符串。
//
// 返回值:
//   p string - 计算得到的哈希值，以Base64编码格式表示。
func SHA256(s string) (p string) {
	// 创建SHA256哈希对象
	o := sha256.New()
	// 向哈希对象中写入字符串的字节序列
	o.Write([]byte(s))
	// 将哈希对象的最终结果编码为Base64字符串
	p = base64.StdEncoding.EncodeToString(o.Sum(nil))
	return
}

// SHA512 用于计算给定字符串的 SHA512 散列值，并以 base64 编码格式返回。
// 参数:
//   s - 需要计算散列值的字符串。
// 返回值:
//   p - 输入字符串 s 的 SHA512 散列值，以 base64 编码格式表示。
func SHA512(s string) (p string) {
	// 创建一个新的 SHA512 散列对象
	o := sha512.New()
	// 向散列对象中写入字符串 s 的字节表示
	o.Write([]byte(s))
	// 将散列对象的最终结果编码为 base64 字符串
	p = base64.StdEncoding.EncodeToString(o.Sum(nil))
	return
}

// UUID 生成一个唯一的UUID字符串。
//
// 参数: 无
//
// 返回值: 一个字符串，代表一个唯一的UUID。
func UUID() string {
	return uuid.New().String()
}
