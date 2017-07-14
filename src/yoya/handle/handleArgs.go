/**************************************************************************

Copyright:YOYA

Author: shaozhenyu

Date:2017-06-29

Description:  function of marshal request to json

**************************************************************************/

package handle

import (
	"errors"
	"reflect"
	"strings"

	"yoya/utils/aes"

	"github.com/valyala/fasthttp"
)

const (
	ADX       = "adx"
	PRICE     = "price"
	EXT       = "ext"
	NORMAL    = "normal"
	LOCATION  = "location"
	LDP       = "ldp"
	SPEC      = "specmedium"
	Time      = "time"
	PRICETYPE = "pricetype"
)

func test() []byte {
	rr := "1|1|5|2|1|1||44|1|1|1|1|12|1|1|bc"

	ciphertext, _ := aes.Encrypter([]byte(rr))
	return ciphertext
}

func stringMarshal(strs []string, ext *Extend) error {
	v := reflect.ValueOf(ext).Elem()
	if len(strs) != v.NumField() {
		return errors.New("request ext error")
	}
	for i := 0; i < v.NumField(); i++ {
		v.Field(i).SetString(strs[i])
	}

	return nil
}

func getQueryArgs(ctx *fasthttp.RequestCtx) (*Message, []byte, string, error) {

	message := &Message{
		Ip: ctx.RemoteIP().String(),
	}

	args := ctx.QueryArgs()

	adxId, err := args.GetUint(ADX)
	if err != nil {
		return nil, nil, "", err
	}
	message.AdxId = int64(adxId)

	price := args.Peek(PRICE)
	if len(price) == 0 {
		return nil, nil, "", errors.New("price is null")
	}

	timeStamp, err := args.GetUint(Time)
	if err != nil {
		return nil, nil, "", err
	}
	message.TimeStamp = int64(timeStamp)

	message.PriceType = string(args.Peek(PRICETYPE))
	if len(message.PriceType) == 0 {
		return nil, nil, "", errors.New("pricetype is null")
	}

	ext := args.Peek(EXT)
	ext = test()
	if len(ext) == 0 {
		return nil, nil, "", errors.New("req is null")
	}
	extStr, err := aes.Decrypter(ext)
	if err != nil {
		return nil, nil, "", err
	}

	extSli := strings.Split(string(extStr), "|")
	err = stringMarshal(extSli, &message.Ext)
	if err != nil {
		return nil, nil, "", err
	}

	if len(message.Ext.BidId) == 0 || len(message.Ext.ImpId) == 0 || len(message.Ext.ReqId) == 0 || len(message.Ext.DeviceId) == 0 {
		return nil, nil, "", errors.New("request ext get error")
	}

	return message, price, string(extStr), nil
}
