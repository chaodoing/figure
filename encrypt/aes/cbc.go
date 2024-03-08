package aes

import (
	`crypto/aes`
	`crypto/cipher`
)

// EncryptCBC 使用AES的CBC模式对原始数据进行加密
// origData: 需要加密的原始数据
// key: 用于加密的密钥
// 返回值: 加密后的数据
func EncryptCBC(origData []byte, key []byte) (encrypted []byte) {
	// 创建AES加密器，要求密钥长度为16, 24或32字节
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize() // 获取块大小
	// 对原始数据进行PKCS5填充以满足加密要求
	origData = pkcs5Padding(origData, blockSize)
	// 创建CBC模式的加密器
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	// 创建用于存放加密结果的字节数组
	encrypted = make([]byte, len(origData))
	// 对数据进行加密
	blockMode.CryptBlocks(encrypted, origData)
	return encrypted
}

// DecryptCBC 使用AES-CBC模式解密数据
// encrypted：被加密的数据字节切片
// key：解密所需的秘钥字节切片
// 返回值 decrypted：解密后的数据字节切片
func DecryptCBC(encrypted []byte, key []byte) (decrypted []byte) {
	// 创建AES分组加密器
	block, _ := aes.NewCipher(key)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 创建CBC解密模式
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	// 准备解密后的数据切片
	decrypted = make([]byte, len(encrypted))
	// 执行解密操作
	blockMode.CryptBlocks(decrypted, encrypted)
	// 移除填充的补全码
	decrypted = pkcs5UnPadding(decrypted)
	return decrypted
}
