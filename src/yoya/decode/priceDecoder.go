/**************************************************************************

Copyright:YOYA

Author: shaozhenyu

Date:2017-06-27

Description: api handle

**************************************************************************/

package decode

import (
	"yoya/decode/adview"
	"yoya/decode/baidu"
	"yoya/decode/huzhong"
	"yoya/decode/kingsoft"
	"yoya/decode/xtrader"
)

//parse price interface
type PriceDecode interface {
	PriceDecoder(price []byte) int64
}

const (
	UnKnow      = iota
	ADIQUITY    // = 1
	INMOBI      // = 2
	NEXAGE      // = 3
	SMAATO      // = 4
	YOUKU       // = 5
	AMOBEE      // = 6
	MOPUB       // = 7
	DOUBLECLICk // = 8
	AXONIX      // = 9
	TAPSENSE    // = 10
	LETV        // = 11
	PPTV        // = 12
	INNERACTIVE // = 13
	IQIYI       // = 14
	MOGO        // = 15
	ADVIEW      // = 16
	WAX         // = 17
	XTRADERID   // = 18
	HUZHONG     // = 19
	KINGSOFT    // = 20
	QICHUANG    // = 21
	BAIDU       // = 22
)

func PriceParse(adx int, price []byte) int64 {
	var decoder PriceDecode
	switch adx {
	case ADVIEW:
		{
			decoder = &adview.Adview{}
		}
	case XTRADERID:
		{
			decoder = &xtrader.Xtrader{}
		}
	case HUZHONG:
		{
			decoder = &huzhong.HuZhong{}
		}
	case KINGSOFT:
		{
			decoder = &kingsoft.KingSoft{}
		}
	case BAIDU:
		{
			decoder = &baidu.Baidu{}
		}
	default:
		return 0
	}

	return decoder.PriceDecoder(price)
}
