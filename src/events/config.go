/**************************************************************************

Copyright:YOYA

Author: shaozhenyu

Date:2017-06-26

Description: event config

**************************************************************************/

package events

import (
	log "utils/log/logApi"
)

type EventConfig struct {
	Host string
	Port int
}

type RedisConfig struct {
	IndexAddr     string
	DuplicateAddr string
}

type Config struct {
	Event EventConfig
	Redis RedisConfig
	Log   log.LogConfig
}
