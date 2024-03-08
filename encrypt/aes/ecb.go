package aes

import (
	`crypto/aes`
	`encoding/base64`
)

// EncryptECB 使用ECB模式对原始数据进行加密
// origData: 需要加密的原始数据
// key: 加密使用的密钥
// 返回值: 加密后的数据
func EncryptECB(ciphertext, chars string) (value string) {
	var encrypted []byte
	origData := []byte(ciphertext)
	key := []byte(chars)
	// 使用提供的密钥生成AES加密器
	cipher, _ := aes.NewCipher(generateKey(key))
	// 计算加密后数据的长度
	length := (len(origData) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, origData)
	// 对原始数据进行填充，保证数据长度是AES块大小的整数倍
	pad := byte(len(plain) - len(origData))
	for i := len(origData); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	// 分组加密数据
	for bs, be := 0, cipher.BlockSize(); bs <= len(origData); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}
	return base64.StdEncoding.EncodeToString(encrypted)
}

// DecryptECB 使用ECB模式解密数据
// encrypted: 需要解密的密文
// key: 解密使用的密钥
// 返回值 decrypted: 解密后的明文
func DecryptECB(ciphertext, chars string) (value string) {
	var decrypted []byte
	encrypted, _ := base64.StdEncoding.DecodeString(ciphertext)
	key := []byte(chars)
	// 使用提供的密钥生成AES加密器
	cipher, _ := aes.NewCipher(generateKey(key))
	decrypted = make([]byte, len(encrypted))
	
	// 按块解密密文
	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}
	
	// 移除可能存在的填充字符
	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}
	
	return string(decrypted[:trim])
}
