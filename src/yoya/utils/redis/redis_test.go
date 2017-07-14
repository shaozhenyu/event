/**************************************************************************

Copyright:YOYA

Author: shaozhenyu

Date:2017-06-26

Description: redis test

**************************************************************************/
package redis

import (
	"testing"
	"time"
)

func Test_NewRedis(t *testing.T) {
	r, err := New("localhost:6379", 0, 100)
	if err != nil {
		t.Error(err)
	}

	r.Close()
}

func Test_KV(t *testing.T) {
	r, err := New("localhost:6379", 0, 100)
	if err != nil {
		t.Error(err)
	}

	err = r.Set("test", []byte("test11"))
	if err != nil {
		t.Error(err)
	}
	val, err := r.Get("test")
	if err != nil {
		t.Error(err)
	}

	if string(val) != "test11" {
		t.Error("set value not equal get value")
	}

	err = r.Set("test", []byte("test12"))
	if err != nil {
		t.Error(err)
	}
	val, err = r.Get("test")
	if err != nil {
		t.Error(err)
	}

	if string(val) == "test11" {
		t.Error("set value not equal get value")
	}

	err = r.Delect("test")
	if err != nil {
		t.Error(err)
	}

	val, err = r.Get("test")
	if err == nil {
		t.Error("redis get key error")
	}

	r.Delect("test")

	r.Close()
}

func Test_SetNX(t *testing.T) {
	r, err := New("localhost:6379", 0, 100)
	if err != nil {
		t.Error(err)
	}

	r.Delect("test")
	result, err := r.SetNX("test", []byte("test11"))
	if err != nil {
		t.Error(err)
	}
	if !result {
		t.Error("redis setnx error")
	}

	result, err = r.SetNX("test", []byte("test22"))
	if err != nil {
		t.Error(err)
	}
	if result {
		t.Error("redis setnx error")
	}

	r.Delect("test")
	r.Close()
}

func Test_SetExpire(t *testing.T) {
	r, err := New("localhost:6379", 0, 100)
	if err != nil {
		t.Error(err)
	}

	r.Delect("test")
	err = r.SetExpire("test", []byte("test11"), 5)
	duration, err := r.conn.TTL("test").Result()
	if err != nil {
		t.Error(err)
	}

	if int(duration.Seconds()) != 5 {
		t.Error("set key expiration error")
	}

	r.Delect("test")
	err = r.SetExpire("test", []byte("test11"), 1)
	_, err = r.Get("test")
	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Millisecond * 1001)
	_, err = r.Get("test")
	if err == nil {
		t.Error("redis setexpire error")
	}

	r.Close()
}

func Test_HM(t *testing.T) {
	r, err := New("localhost:6379", 0, 100)
	if err != nil {
		t.Error(err)
	}

	r.Delect("test")
	err = r.HMSet("test", "k1", "v1", "k2", "v2")
	if err != nil {
		t.Error(err)
	}

	values, err := r.HMGet("test", "k1", "k2")
	trueVal := []interface{}{"v1", "v2"}
	if trueVal[0] != values[0] && trueVal[1] != values[1] {
		t.Error("redis HMGET error")
	}

	r.Delect("test")
	r.Close()
}

func Test_FlushDB(t *testing.T) {
	r, err := New("localhost:6379", 0, 100)
	if err != nil {
		t.Error(err)
	}

	err = r.Set("test", []byte("test11"))
	if err != nil {
		t.Error(err)
	}
	_, err = r.Get("test")
	if err != nil {
		t.Error(err)
	}

	err = r.FlushDB()
	if err != nil {
		t.Error(err)
	}

	_, err = r.Get("test")
	if err == nil {
		t.Error("redis flushdb error")
	}

	r.Close()
}
