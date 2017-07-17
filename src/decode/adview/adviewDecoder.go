package adview

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"

	log "utils/log/logApi"
)

// ekey：Zp9QpIWkRqbNExs3cDJgxNrGM0fOJC8T
// ikey：fJNBPNm4ToyA7DP4c2ERBfKEKdMTbTc4
// 加密后的价格为：GjXscFgBAABfdiAxHDBKPh1YxUagX-zjhGEnvA

const (
	InitSliceSize = 16
	SignatureSize = 4
	BlockSize     = 20
)

const EKEY = "Zp9QpIWkRqbNExs3cDJgxNrGM0fOJC8T"

type Adview struct {
}

func (ad *Adview) PriceDecoder(price []byte) int64 {
	log.Debugs("In AdviewPriceParse()", "")
	cipherText, err := base64.RawURLEncoding.DecodeString(string(price))
	if err != nil {
		log.Errors("In Adview PriceDecoder() base64.RawURLEncoding", fmt.Sprintf("err:%s", err))
		return 0
	}
	clearTextLength := len(cipherText) - InitSliceSize - SignatureSize
	if clearTextLength < 0 {
		log.Errors("In Adview PriceDecoder() cleartextLength < 0", "")
		return 0
	}

	text := cipherText[0:16]
	clearText := make([]byte, clearTextLength)
	clearTextBegin := 0
	cipherTextBegin := len(text)
	cipherTextEnd := cipherTextBegin + cap(clearText)

	encryptionPad := computeHmac1(text, EKEY)
	for i := 0; i < BlockSize && cipherTextBegin < cipherTextEnd; {
		clearText[clearTextBegin] = cipherText[cipherTextBegin] ^ encryptionPad[i]
		i++
		clearTextBegin++
		cipherTextBegin++
	}

	intPrice := ntohll(clearText)
	return intPrice * 100 / 1000 // adview是以元为单位*10000，我们内部是以元为单位*1000000
}

func computeHmac1(message []byte, secret string) []byte {
	key := []byte(secret)
	h := hmac.New(sha1.New, key)
	h.Write(message)
	return h.Sum(nil)
}

// 将一个无符号长整形数从网络字节顺序转换为主机字节顺序
func ntohll(price []byte) int64 {
	var value uint64
	buf := bytes.NewReader(price)
	err := binary.Read(buf, binary.BigEndian, &value)
	if err != nil {
		log.Errors("In Adview  ntohll() binary.Read failed", fmt.Sprintf("err:%s", err))
		return 0
	}
	return int64(value)
}
