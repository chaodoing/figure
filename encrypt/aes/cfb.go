package aes

import (
	`crypto/aes`
	`crypto/cipher`
	`crypto/rand`
	`io`
)

// EncryptCFB 使用AES的CFB模式对原始数据进行加密
// origData: 需要加密的原始数据
// key: 用于加密的密钥
// 返回值: 加密后的数据
func EncryptCFB(origData []byte, key []byte) (encrypted []byte) {
	// 使用密钥创建AES加密器
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err) // 如果创建加密器失败，则直接panic
	}
	// 创建存储加密结果的字节数组，长度为一个块大小加上原始数据的长度
	encrypted = make([]byte, aes.BlockSize+len(origData))
	iv := encrypted[:aes.BlockSize] // 分配空间给初始化向量
	
	// 从随机数源读取初始化向量
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err) // 如果读取初始化向量失败，则直接panic
	}
	// 创建CFB加密器，并使用初始化向量
	stream := cipher.NewCFBEncrypter(block, iv)
	// 使用加密器对原始数据进行XOR运算，加密原始数据
	stream.XORKeyStream(encrypted[aes.BlockSize:], origData)
	return encrypted // 返回加密后的数据
}

// DecryptCFB 使用AES的CFB模式解密密文
// 参数：
//   encrypted []byte: 待解密的密文
//   key []byte: 解密使用的密钥
// 返回值：
//   decrypted []byte: 解密后的明文
func DecryptCFB(encrypted []byte, key []byte) (decrypted []byte) {
	// 使用密钥创建AES密码器
	block, _ := aes.NewCipher(key)
	// 检查密文长度是否足够
	if len(encrypted) < aes.BlockSize {
		panic("ciphertext too short")
	}
	// 从密文中提取初始化向量IV
	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]
	// 创建CFB解密器，并使用提取的IV
	stream := cipher.NewCFBDecrypter(block, iv)
	// 使用解密器对密文进行解密
	stream.XORKeyStream(encrypted, encrypted)
	return encrypted
}
