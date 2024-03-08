package encrypt

import (
	`bytes`
	`crypto/aes`
	`crypto/cipher`
	`encoding/base64`
	`errors`
	`fmt`
)

type Aes struct {
	IVString  string `json:"IVString"`
	SecretKey string `json:"secretKey"`
}

func AES(SecretKey string) Aes {
	return Aes{
		IVString:  "A-16-Byte-String",
		SecretKey: SecretKey,
	}
}

func (a Aes) ToString(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func (a Aes) AESEncrypt(data string) (baseString string, err error) {
	block, err := aes.NewCipher([]byte(a.SecretKey))
	if err != nil {
		return "", err
	}
	if data == "" {
		return "", errors.New("plain Content Empty")
	}
	var crypted []byte
	ecb := cipher.NewCBCEncrypter(block, []byte(a.IVString))
	content := []byte(data)
	content = a.PKCS5Padding(content, block.BlockSize())
	crypted = make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)
	baseString = base64.StdEncoding.EncodeToString(crypted)
	return
}

func (a Aes) AESDecrypt(content string) (data []byte, err error) {
	crypt, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return
	}
	block, err := aes.NewCipher([]byte(a.SecretKey))
	if err != nil {
		fmt.Println("Aes")
		return nil, err
	}
	if len(crypt) == 0 {
		return nil, errors.New("plain content empty")
	}
	ecb := cipher.NewCBCDecrypter(block, []byte(a.IVString))
	decrypted := make([]byte, len(crypt))
	ecb.CryptBlocks(decrypted, crypt)
	return a.PKCS5Trimming(decrypted), nil
}

func (a Aes) PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func (a Aes) PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}
