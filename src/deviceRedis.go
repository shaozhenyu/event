package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"utils/redis"
)

const (
	DefaultRdsAddr = "localhost:6379"
)

const (
	RedisDeviceClearText = "devicecleartext" //明文
	RedisDeviceMd5       = "devicemd5"
	RedisDeviceSha1      = "devicesha1"
)

var (
	rdsAddr  = flag.String("rdsaddr", DefaultRdsAddr, "The address of connect redis")
	filepath = flag.String("filepath", "", "The filepath to read device")
)

func main() {
	flag.Parse()
	if len(*filepath) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	RdsConn, err := redis.New(*rdsAddr, 0, 100)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(RdsConn)

	f, err := os.Open(*filepath)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	err = handle(RdsConn, f)
	if err != nil {
		log.Fatalln(err)
	}
}

func handle(conn *redis.AdRedis, f *os.File) error {
	buf := bufio.NewReader(f)
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		handleOne(conn, string(line))
	}
	return nil
}

func handleOne(conn *redis.AdRedis, devid string) error {
	devMd5 := getMd5(devid)
	devSha1 := getSha1(devid)

	_, err := conn.HSETNX(RedisDeviceClearText, devid, 1)
	if err != nil {
		return err
	}

	_, err = conn.HSETNX(RedisDeviceMd5, devMd5, 1)
	if err != nil {
		return err
	}

	_, err = conn.HSETNX(RedisDeviceSha1, devSha1, 1)
	if err != nil {
		return err
	}

	return nil
}

func getMd5(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func getSha1(str string) string {
	h := sha1.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}
