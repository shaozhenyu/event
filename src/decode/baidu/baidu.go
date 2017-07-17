/**************************************************************************

Copyright:YOYA

Author: shaozhenyu

Date:2017-07-05

Description: baidu decode

**************************************************************************/

package baidu

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"

	log "utils/log/logApi"
)

const (
	InitSliceSize = 16
	SignatureSize = 4
	BlockSize     = 20
)

var encryption_key_v = []byte{
	0x01, 0x71, 0x55, 0xc4, 0x00, 0xdb, 0xd6, 0xb6,
	0x60, 0xaa, 0x4c, 0x1f, 0x00, 0xdb, 0xd6, 0xb6,
	0x60, 0xaa, 0x6d, 0xfa, 0x00, 0xdb, 0xd6, 0xb6,
	0x60, 0xaa, 0x7e, 0x0f, 0x94, 0x73, 0x6b, 0x85,
}

type Baidu struct {
}

func (bd *Baidu) PriceDecoder(price []byte) int64 {
	log.Debugs("In BaiduPriceParse()", "")
	cipherText, err := base64.RawURLEncoding.DecodeString(string(price))
	if err != nil {
		log.Errors("In Baidu PriceDecoder() base64.RawURLEncoding", fmt.Sprintf("err:%s", err))
		return 0
	}
	clearTextLength := len(cipherText) - InitSliceSize - SignatureSize
	if clearTextLength < 0 {
		log.Errors("In Baidu PriceDecoder() cleartextLength < 0", "")
		return 0
	}

	text := cipherText[0:16]
	clearText := make([]byte, clearTextLength)
	clearTextBegin := 0
	cipherTextBegin := len(text)
	cipherTextEnd := cipherTextBegin + cap(clearText)

	encryptionPad := computeHmac1(text, string(encryption_key_v))
	for i := 0; i < BlockSize && cipherTextBegin < cipherTextEnd; {
		clearText[clearTextBegin] = cipherText[cipherTextBegin] ^ encryptionPad[i]
		i++
		clearTextBegin++
		cipherTextBegin++
	}

	intPrice := ntohll(clearText)
	return intPrice
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
		log.Errors("In Baidu  ntohll() binary.Read failed", fmt.Sprintf("err:%s", err))
		return 0
	}
	return int64(value)
}
