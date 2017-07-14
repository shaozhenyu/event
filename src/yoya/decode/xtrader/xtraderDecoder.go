package xtrader

import (
	"crypto/aes"
	"crypto/cipher"
	b64 "encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	log "yoya/utils/log/logApi"
)

const XTRADER_TOKEN = "fd97371e5c1e4796851bd97c9a5c0f9b"

var std = b64.URLEncoding
var url = b64.RawURLEncoding

type Xtrader struct {
}

// 计算xtrader price
func (xtrader *Xtrader) PriceDecoder(price []byte) int64 {
	key, _ := hex.DecodeString(XTRADER_TOKEN)
	content, err := xtraderDecrypt(price, key)
	if err != nil {
		log.Errors("In xtrader PricePares() Decrypt", fmt.Sprintf("err:%s", err))
		return 0
	}
	if len(content) > 0 {
		array := strings.Split(string(content), "_")
		intPrice, err := strconv.ParseInt(array[0], 10, 0)
		if err != nil {
			log.Errors("In xtrader PricePares() parse price", fmt.Sprintf("err:%s", err))
			return 0
		}
		return intPrice * 1000000 / 100 / 1000 // 灵集是以分为单位，我们内部是以元为单位*1000000
	}
	return 0
}

// base64解码
func base64Decode(src []byte) ([]byte, error) {
	if src[len(src)-1] == byte(b64.StdPadding) {
		return std.DecodeString(string(src))
	}
	return url.DecodeString(string(src))
}

// aes解密
func xtraderDecrypt(price, key []byte) ([]byte, error) {
	Decodetxt, err := base64Decode(price)
	if err != nil {
		log.Errors("In xtraderDecrypt() Decode()", fmt.Sprintf("err:%s", err))
		return nil, err
	}
	dest := make([]byte, (len(Decodetxt)/len(key)+1)*len(key))
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		log.Errors("In xtraderDecrypt() NewCipher", fmt.Sprintf("err:%s", err))
		return nil, err
	}
	encrypter := NewECBDecrypter(aesCipher)
	encrypter.CryptBlocks(dest, []byte(Decodetxt))
	return dest, nil
}

// aes加密
type ECB struct {
	b         cipher.Block
	blockSize int
}

func newECB(bk cipher.Block) *ECB {
	return &ECB{
		b:         bk,
		blockSize: bk.BlockSize(),
	}
}

type ecbEncrypter ECB

func NewECBEncrypter(bk cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(bk))
}

func (pter *ecbEncrypter) BlockSize() int { return pter.blockSize }

func (pter *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%pter.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		pter.b.Encrypt(dst, src[:pter.blockSize])
		src = src[pter.blockSize:]
		dst = dst[pter.blockSize:]
	}
}

type ecbDecrypter ECB

func NewECBDecrypter(bk cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(bk))
}

func (pter *ecbDecrypter) BlockSize() int { return pter.blockSize }

func (pter *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%pter.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}

	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}

	for len(src) > 0 {
		pter.b.Decrypt(dst, src[:pter.blockSize])
		src = src[pter.blockSize:]
		dst = dst[pter.blockSize:]
	}
}
