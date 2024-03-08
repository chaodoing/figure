package test

import (
	`testing`
	
	`github.com/chaodoing/figure/encrypt/aes`
)

func TestAes(t *testing.T) {
	key := aes.Key("123456")
	code := aes.EncryptCBC("chaodoing@live.com", key)
	t.Log(code)
	t.Log(aes.DecryptCBC(code, key))
	
	code = aes.EncryptCFB("192.168.cc", key)
	t.Log(code)
	t.Log(aes.DecryptCFB(code, key))
	
	code = aes.EncryptECB("chaodoing@hotmail.com", key)
	t.Log(code)
	t.Log(aes.DecryptECB(code, key))
}
