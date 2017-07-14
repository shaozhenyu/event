/**************************************************************************

Copyright:YOYA

Author: shaozhenyu

Date:2017-06-27

Description: event redis handler

**************************************************************************/

package handle

import (
	"fmt"
	"time"

	log "yoya/utils/log/logApi"
)

const TimeZone = 28800 // 8h**60m*60s
const TimeFormat = "2006-01-02 15:04:05"

func getExpire(midnight int64) int64 {
	tm := time.Now().Unix()
	return midnight - tm
}

func setTodayExpire(key string) error {
	date := time.Unix(time.Now().Unix(), 0).Format("2006-01-02")
	tm, _ := time.Parse(TimeFormat, date+" "+"23:59:59")
	expire := getExpire(tm.Unix() - TimeZone)
	return RdsConn.Expire(key, expire)
}

//redis中 设置key的值，如果key已经存在，自加1。
//如果key不存在。则设置key的值为1，并设置过期时间到当天24点
func incrKey(key string) (int64, error) {
	value, err := RdsConn.Incr(key)
	if err != nil {
		return 0, err
	}

	if value == 1 {
		err = setTodayExpire(key)
		if err != nil {
			return 0, err
		}
	}

	return value, nil
}

//redis记录当此请求
func RedisRecord(reqType string, member string) (int64, error) {
	log.Debugs("In RedisRecord() ADD TO REDIS", fmt.Sprintf("%s:%s", reqType, member))
	today := time.Unix(time.Now().Unix(), 0).Format("20060102")
	key := today + "-" + reqType
	exists, err := RdsConn.Exists(key)
	if err != nil {
		return 0, err
	}
	//设置过期时间
	if !exists {
		value, err := RdsConn.HINCRBY(key, member, 1)
		if err != nil {
			return 0, err
		}

		err = setTodayExpire(key)
		if err != nil {
			return 0, err
		}
		return value, nil
	}
	return RdsConn.HINCRBY(key, member, 1)
}
