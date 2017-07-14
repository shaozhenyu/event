package huzhong

import (
	"crypto/rc4"
	b64 "encoding/base64"
	"fmt"
	"strconv"

	log "yoya/utils/log/logApi"
)

const HUZHONG_TOKEN = "98466a7b749a6b4509c61c65651bf58b6ad730f2"

var std = b64.URLEncoding
var url = b64.RawURLEncoding

type HuZhong struct {
}

// 解析互众price
func (hz *HuZhong) PriceDecoder(price []byte) int64 {
	log.Debugs("In HuZhong PriceDecoder()", "")
	content, err := base64Decode(price)
	if err != nil {
		log.Errors("In HuZhong PriceDecoder() base64Decode()", fmt.Sprintf("err:%s", err))
		return 0
	}
	rc4Price := rc4Decode(content)
	if len(rc4Price) > 0 {
		float, err := strconv.ParseFloat(string(rc4Price), 32)
		if err != nil {
			log.Errors("In HuZhong PriceDecoder() ParseFloat()", fmt.Sprintf("err:%s", err))
			return 0
		}
		return int64(float * 1000000 / 100 / 1000) // 互众是以分为单位，我们内部是以元为单位*1000000
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

func rc4Decode(price []byte) []byte {
	rc4obj, _ := rc4.NewCipher([]byte(HUZHONG_TOKEN))
	text := make([]byte, len(price))
	rc4obj.XORKeyStream(text, price)
	return text
}
