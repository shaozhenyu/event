/**************************************************************************

Copyright:YOYA

Author: shaozhenyu

Date:2017-06-27

Description: kingsoft price decoder

**************************************************************************/

package kingsoft

import (
	"fmt"
	"strconv"

	log "utils/log/logApi"
)

type KingSoft struct {
}

// 解析金山price
func (k *KingSoft) PriceDecoder(price []byte) int64 {
	log.Debugs("In KingSoft PriceDecoder()", "")
	if len(price) > 0 {
		float, err := strconv.ParseFloat(string(price), 32)
		if err != nil {
			log.Errors("In KingSoft PriceDecoder() ParseFloat()", fmt.Sprintf("err:%s", err))
			return 0
		}
		return int64(float * 100) // 金山是以分为单位，我们内部是以分为单位*100
	}

	return 0
}
