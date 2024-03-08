package aes

import (
	`bytes`
)

// generateKey 函数根据输入的 key 生成一个长度为 16 的新 key。
// 参数 key：输入的原始 key，长度不限。
// 返回值 genKey：生成的长度为 16 的新 key。
func generateKey(key []byte) (genKey []byte) {
	// 初始化生成的 key，长度为 16，内容复制自原始 key 的前 16 个字节。
	genKey = make([]byte, 16)
	copy(genKey, key)
	
	// 遍历原始 key 的剩余部分，依次与生成的 key 的每个字节进行异或操作。
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

// pkcs5Padding 对密文进行 PKCS#5 填充。
// 参数：
//   ciphertext []byte：待填充的密文。
//   blockSize int：块大小，表示填充的单位长度。
// 返回值：
//   []byte：经过 PKCS#5 填充后的密文。
func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	// 计算需要填充的字节长度
	padding := blockSize - len(ciphertext)%blockSize
	// 生成填充的字节序列
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	// 将填充序列添加到密文末尾
	return append(ciphertext, padtext...)
}

// pkcs5UnPadding 去除 PKCS#5 填充的数据
// origData: 待处理的原始字节切片数据
// 返回值: 去除填充后的原始数据字节切片
func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)                // 获取原始数据长度
	unpadding := int(origData[length-1])   // 从最后一位获取填充的字节数
	return origData[:(length - unpadding)] // 返回去除填充后的数据
}
