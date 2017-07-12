/**************************************************************************

Copyright:YOYA

Author: shaozhenyu

Date:2017-06-23

Description:  function of redis

**************************************************************************/

package redis

import (
	"errors"
	"time"

	"gopkg.in/redis.v5"
)

type AdRedis struct {
	conn *redis.Client
}

var (
	ErrRedisNotFound = errors.New("redis not found")
)

func New(addr string, db, maxPoolSize int) (*AdRedis, error) {
	opt := redis.Options{}
	opt.Addr = addr
	opt.DB = db
	opt.PoolSize = maxPoolSize

	conn := redis.NewClient(&opt)
	err := conn.Ping().Err()

	return &AdRedis{conn: conn}, err
}

func (this *AdRedis) Set(key string, value []byte) error {
	return this.conn.Set(key, value, 0).Err()
}

func (this *AdRedis) Get(key string) ([]byte, error) {
	s, err := this.conn.Get(key).Result()
	if err == redis.Nil {
		err = ErrRedisNotFound
	}
	return []byte(s), err
}

func (this *AdRedis) SetExpire(key string, value []byte, expire int64) error {
	return this.conn.Set(key, value, time.Duration(expire)*time.Second).Err()
}

func (this *AdRedis) SetNX(key string, value []byte) (bool, error) {
	return this.conn.SetNX(key, value, 0).Result()
}

func (this *AdRedis) SetNXExpire(key string, value []byte, expire int64) (bool, error) {
	return this.conn.SetNX(key, value, time.Duration(expire)*time.Second).Result()
}

func (this *AdRedis) Delect(key string) error {
	return this.conn.Del(key).Err()
}

func (this *AdRedis) Incr(key string) (int64, error) {
	return this.conn.Incr(key).Result()
}

func (this *AdRedis) IncrBy(key string, delta int64) (int64, error) {
	return this.conn.IncrBy(key, delta).Result()
}

func (this *AdRedis) Expire(key string, expire int64) error {
	return this.conn.Expire(key, time.Duration(expire)*time.Second).Err()
}

func (this *AdRedis) HSET(key string, field string, value interface{}) error {
	return this.conn.HSet(key, field, value).Err()
}

func (this *AdRedis) HGET(key string, field string) ([]byte, error) {
	s, err := this.conn.HGet(key, field).Result()
	if err == redis.Nil {
		err = ErrRedisNotFound
	}
	return []byte(s), err
}

func (this *AdRedis) HINCRBY(key string, field string, incr int64) (int64, error) {
	return this.conn.HIncrBy(key, field, incr).Result()
}

func generateMap(strs []string) (map[string]string, error) {
	fields := map[string]string{}
	if len(strs)%2 != 0 {
		return nil, errors.New("HMSet fields must be even numbers")
	}
	for i := 0; i < len(strs); i = i + 2 {
		fields[strs[i]] = strs[i+1]
	}
	return fields, nil
}

func (this *AdRedis) HMSet(key string, strs ...string) error {
	fields, err := generateMap(strs)
	if err != nil {
		return err
	}
	return this.conn.HMSet(key, fields).Err()
}

func (this *AdRedis) HMGet(key string, fields ...string) ([]interface{}, error) {
	s, err := this.conn.HMGet(key, fields...).Result()
	if err == redis.Nil {
		err = ErrRedisNotFound
	}
	return s, err
}

func (this *AdRedis) Exists(key string) (bool, error) {
	result, err := this.conn.Exists(key).Result()
	if err == redis.Nil {
		return false, ErrRedisNotFound
	}
	return result, err
}

func (this *AdRedis) FlushDB() error {
	return this.conn.FlushDb().Err()
}

func (this *AdRedis) Close() error {
	return this.conn.Close()
}
