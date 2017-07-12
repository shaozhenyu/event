/**************************************************************************

Copyright:YOYA

Author: shaozhenyu

Date:2017-06-22

Description:  bidder event main

**************************************************************************/

package main

import (
	"fmt"

	"events"
	"handle"
	log "utils/log/logApi"
	"utils/redis"
	"utils/yaml"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

var (
	configFile = "config.yml"
)

func main() {
	conf := &events.Config{}
	err := yaml.Read(configFile, conf)
	if err != nil {
		log.Errors("In main() yaml.Read", fmt.Sprintf("err:%s", err))
		return
	}

	log.InitLog(&conf.Log)

	handle.RdsConn, err = redis.New(conf.Redis.IndexAddr, 0, 100)
	if err != nil {
		log.Errors("In main() redis.New", fmt.Sprintf("err:%s", err))
		return
	}

	router := fasthttprouter.New()
	Register(router)

	if err := fasthttp.ListenAndServe(fmt.Sprintf(":%d", conf.Event.Port), router.Handler); err != nil {
		log.Errors("In main() ListenAndServe", fmt.Sprintf("err:%s", err))
	}
}

func Register(router *fasthttprouter.Router) {
	router.GET("/notice", handle.WinNoticeHandler)
	router.GET("/impression", handle.ImpressionHandler)
	router.GET("/click", handle.ClickHandler)
}
