/**************************************************************************

Copyright:YOYA

Author: shaozhenyu

Date:2017-06-26

Description: api handle

**************************************************************************/

package handle

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"decode"
	log "utils/log/logApi"
	"utils/redis"

	"github.com/valyala/fasthttp"
)

var (
	RdsConn *redis.AdRedis
)

const (
	StrNotice     = "notice"
	StrImpression = "impression"
	StrClick      = "click"
)

const (
	Symbol = "|"
)

type Extend struct {
	AppId       string `json:"appid,omitempty"`       //yoya媒体ID
	TagId       string `json:"tagid,omitempty"`       //yoya广告位ID
	Bundle      string `json:"bundle,omitempty"`      //开发包名
	Os          string `json:"os,omitempty"`          //设备类型
	DeviceType  string `json:"devicetype,omitempty"`  //设备加密类型
	DeviceId    string `json:"deviceid,omitempty"`    //设备ID
	CompanyId   string `json:"companyid,omitempty"`   //广告主ID
	OrderId     string `json:"orderid,omitempty"`     //订单ID
	CampaignId  string `json:"campaignid,omitempty"`  //计划ID
	StrategyId  string `json:"strategyid,omitempty"`  //策略ID
	CreativeId  string `json:"creativeid,omitempty"`  //创意ID
	ImpId       string `json:"impid,omitempty"`       //曝光ID
	ReqId       string `json:"reqid,omitempty"`       //请求ID
	BidId       string `json:"bidid,omitempty"`       //竞价ID
	PersonTagId string `json:"persontagid,omitempty"` //人群标签组ID
	City        string `json:"city,omitempty"`        //经过IP库转换的地区和城市
}

type Message struct {
	AdxId     int64  `json:"adxid,omitempty"`     //ADX在yoya平台的ID
	Price     int64  `json:"price,omitempty"`     //价格
	PriceType string `json:"pricetype,omitempty"` //竞价方式
	Ip        string `json:"ip,omitempty"`        //用户当前真实ip地址
	Status    int    `json:"status"`              //状态
	Type      string `json:"type,omitempty"`      //请求类型
	TimeStamp int64  `json:"timestamp,omitempty"` //状态发生的时间
	Ext       Extend `json:"ext,omitempty"`       //扩展
}

func WinNoticeHandler(ctx *fasthttp.RequestCtx) {
	log.Debugs("In WinNoticeHandler() RequestURI", string(ctx.RequestURI()))
	msg, price, code, err := getQueryArgs(ctx)
	if err != nil {
		log.Errors("In WinNoticeHandler() getQueryArgs", fmt.Sprintf("err:%s", err))
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	msg.Type = StrNotice

	result, err := incrKey(StrNotice + Symbol + msg.Ext.ReqId + Symbol + msg.Ext.BidId + Symbol + msg.Ext.ImpId)
	if err != nil {
		log.Errors("In WinNoticeHandler incrKey", fmt.Sprintf("err:%s", err))
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	if result != 1 {
		log.Debugs("In WinNoticeHandler() ExistsKey", "win already exist")
		msg.Status = 1
	}

	msg.Price = decode.PriceParse(int(msg.AdxId), price)
	member := code + Symbol + strconv.Itoa(int(msg.Price)) + Symbol + strconv.Itoa(msg.Status)
	_, err = RedisRecord(StrNotice, member)
	if err != nil {
		log.Errors("In WinNoticeHandler() RedisRecord error!", fmt.Sprintf("err:%s", err))
	}

	logMsg, err := json.Marshal(msg)
	if err != nil {
		log.Errors("In WinNoticeHandler() json.Marshal error!", "")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	} else {
		log.Infos("In WinNoticeHandler()", string(logMsg))
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func ImpressionHandler(ctx *fasthttp.RequestCtx) {
	log.Debugs("In ImpressionHandler() RequestURI", string(ctx.RequestURI()))
	msg, price, code, err := getQueryArgs(ctx)
	if err != nil {
		log.Errors("In ImpressionHandler() getQueryArgs", fmt.Sprintf("err:%s", err))
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	msg.Price = decode.PriceParse(int(msg.AdxId), price)
	normal, err := ctx.QueryArgs().GetUint(NORMAL)
	if err == nil {
		// 特殊媒体 需要解析价格
		if normal == 0 {
			msg.Type = StrNotice
			result, err := incrKey(StrNotice + Symbol + msg.Ext.ReqId + Symbol + msg.Ext.BidId + Symbol + msg.Ext.ImpId)
			if err != nil {
				log.Errors("In ImpressionHandler() incrKey", fmt.Sprintf("err:%s", err))
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}
			if result != 1 {
				log.Debugs("In ImpressionHandler() ExistsKey", "win already exist")
				msg.Status = 1
			}

			logMsg, err := json.Marshal(msg)
			if err != nil {
				log.Errors("In ImpressionHandler() json.Marshal error!", "")
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			} else {
				log.Infos("In ImpressionHandler()", string(logMsg))
			}
		}
	}

	msg.Type = StrImpression
	//先判断是否有赢得过竞价
	ok, err := RdsConn.Exists(StrNotice + Symbol + msg.Ext.ReqId + Symbol + msg.Ext.BidId + Symbol + msg.Ext.ImpId)
	if err != nil {
		log.Errors("In ImpressionHandler() redis error!", fmt.Sprintf("err:%s", err))
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	if ok {
		result, err := incrKey(StrImpression + Symbol + msg.Ext.ReqId + Symbol + msg.Ext.BidId + Symbol + msg.Ext.ImpId)
		if err != nil {
			log.Errors("In ImpressionHandler() incrKey", fmt.Sprintf("err:%s", err))
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
		if result != 1 {
			log.Debugs("In ImpressionHandler() ExistsKey", "win already exist")
			msg.Status = 1
		}
	} else {
		msg.Status = 2
	}

	member := code + Symbol + strconv.Itoa(int(msg.Price)) + Symbol + strconv.Itoa(msg.Status)
	member = member + Symbol + msg.Ip + Symbol + msg.Ext.City
	_, err = RedisRecord(StrImpression, member)
	if err != nil {
		log.Errors("In ImpressionHandler() RedisRecord error!", fmt.Sprintf("err:%s", err))
	}

	logMsg, err := json.Marshal(msg)
	if err != nil {
		log.Errors("In ImpressionHandler() json.Marshal error!", "")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	} else {
		log.Infos("In ImpressionHandler()", string(logMsg))
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func ClickHandler(ctx *fasthttp.RequestCtx) {
	log.Debugs("In ClickHandler() RequestURI", string(ctx.RequestURI()))
	msg, price, code, err := getQueryArgs(ctx)
	if err != nil {
		log.Errors("In ClickHandler() getQueryArgs", fmt.Sprintf("err:%s", err))
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	ok, err := RdsConn.Exists(StrImpression + Symbol + msg.Ext.ReqId + Symbol + msg.Ext.BidId + Symbol + msg.Ext.ImpId)
	if err != nil {
		log.Errors("In ClickHandler() redis error!", fmt.Sprintf("err:%s", err))
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if ok {
		result, err := incrKey(StrClick + Symbol + msg.Ext.ReqId + Symbol + msg.Ext.BidId + Symbol + msg.Ext.ImpId)
		if err != nil {
			log.Errors("In ClickHandler() incrKey", fmt.Sprintf("err:%s", err))
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
		if result != 1 {
			log.Debugs("In ClickHandler() ExistsKey", "win already exist")
			msg.Status = 1
		}

	} else {
		msg.Status = 2
	}

	msg.Price = decode.PriceParse(int(msg.AdxId), price)
	member := code + Symbol + strconv.Itoa(int(msg.Price)) + Symbol + strconv.Itoa(msg.Status)
	member = member + Symbol + msg.Ip + Symbol + msg.Ext.City
	_, err = RedisRecord(StrClick, member)
	if err != nil {
		log.Errors("In ClickHandler() RedisRecord error!", fmt.Sprintf("err:%s", err))
	}

	msg.Type = StrClick
	logMsg, err := json.Marshal(msg)
	if err != nil {
		log.Errors("In ClickHandler() json.Marshal error!", "")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	} else {
		log.Infos("In ClickHandler()", string(logMsg))
	}

	query := ctx.QueryArgs()
	ldp := query.Peek(LDP)
	spec, specErr := query.GetUint(SPEC)
	if specErr != nil {
		//不是特殊媒体,没有这个参数
		ctx.SetStatusCode(fasthttp.StatusOK)
	} else {
		if spec == 1 {
			// 重定向
			ctx.RedirectBytes(ldp, fasthttp.StatusFound)
		} else if spec == 0 {
			// 代理
			resp, err := http.Get(string(ldp))
			if err != nil {
				log.Errors("In ClickHandler() http.Get() fail!", fmt.Sprintf("%s", err))
				ctx.SetStatusCode(fasthttp.StatusForbidden)
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Errors("In ClickHandler() parse rspBody fail!", fmt.Sprintf("%s", err))
				ctx.SetStatusCode(fasthttp.StatusForbidden)
				return
			}
			ctx.SetContentType("application/xml;charset=utf-8")
			ctx.SetBody(body)
			ctx.SetStatusCode(fasthttp.StatusOK)
		}
	}
}
