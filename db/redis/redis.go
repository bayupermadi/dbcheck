package redis

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"github.com/bayupermadi/dbcheck"
	"github.com/bayupermadi/dbcheck/registry"

	"github.com/gomodule/redigo/redis"
)

type rediss struct {
	con redis.Conn
}

func getValue(str string) string {
	split := strings.Split(str, ":")
	return split[1]
}

func getString(info string, prefix string) (string, error) {
	var text string
	reader := bufio.NewReader(bytes.NewBufferString(info))
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, prefix) {
			text = line
			break
		}
	}
	return text, nil
}

func getUsage(info string) (string, error) {
	str, err := getString(info, "used_memory")
	if err != nil {
		return "", err
	}

	usage := getValue(str)
	return usage, nil
}

func getKeys(info string) ([]string, error) {
	strExpired, err := getString(info, "expired_keys")
	if err != nil {
		return nil, err
	}

	strEvicted, err := getString(info, "evicted_keys")
	if err != nil {
		return nil, err
	}

	exp := getValue(strExpired)
	evi := getValue(strEvicted)
	return []string{exp, evi}, nil
}

func (r *rediss) getSlowlogCount() (int, error) {
	count, err := redis.Int(r.con.Do("SLOWLOG", "LEN"))
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *rediss) getSlowLog(count int) error {
	slowlog, err := redis.Values(r.con.Do("SLOWLOG", "GET", count))
	if err != nil {
		fmt.Println(err)
		return err
	}
	// fmt.Println(slowlog[0].([]interface{})[4])
	for i := 0; i < len(slowlog); i++ {

		for j := i; j < len(slowlog[i].([]interface{})); j++ {
			fmt.Println(slowlog[i].([]interface{}))
		}
	}
	return nil
}

func (r *rediss) getMemoryStats() error {
	stats, err := redis.Values(r.con.Do("MEMORY", "STATS"))
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf(" Memory Stats\n  Peak Allocated: %d \n  Total Allocated: %d\n  Startup Allocated: %d\n  Peak Percentage: %s\n  Fragmentation: %s\n", stats[1], stats[3], stats[5], stats[27], stats[29])
	return nil
}

func (r *rediss) redisDoString(command string, args ...interface{}) (string, error) {
	info, err := redis.String(r.con.Do(command, args...))
	if err != nil {
		return "", err
	}
	return info, nil
}

func (r *rediss) redisDoInt(command string, args ...interface{}) (int, error) {
	info, err := redis.Int(r.con.Do(command, args...))
	if err != nil {
		return 0, err
	}
	return info, nil
}

func (r *rediss) Version() error {
	info, err := r.redisDoString("INFO", "SERVER")
	if err != nil {
		return err
	}

	strVersion, err := getString(info, "redis_version")
	if err != nil {
		return err
	}

	strOs, err := getString(info, "os")
	if err != nil {
		return err
	}

	strGcc, err := getString(info, "gcc_version")
	if err != nil {
		return err
	}

	version := getValue(strVersion)
	os := getValue(strOs)
	gcc := getValue(strGcc)

	fmt.Printf("Redis version: %s OS %s gcc_version %s \n", version, os, gcc)

	return nil
}

func (r *rediss) ActiveClient() error {
	info, err := r.redisDoString("INFO", "CLIENTS")
	if err != nil {
		return err
	}

	str, err := getString(info, "connected_clients")
	if err != nil {
		return err
	}

	client := getValue(str)

	fmt.Printf("active_client(s): %s\n", client)

	return nil
}

func (r *rediss) Health() error {
	fmt.Println("health_status:")
	size, err := r.redisDoInt("DBSIZE")
	if err != nil {
		return err
	}

	info, err := r.redisDoString("INFO")
	if err != nil {
		return err

	}

	usage, err := getUsage(info)
	if err != nil {
		return err
	}

	keys, err := getKeys(info)
	if err != nil {
		return err
	}

	countLog, err := r.getSlowlogCount()
	if err != nil {
		return err
	}

	fmt.Printf(" Available Key: %d\n Memory Usage: %s\n Expired Keys: %s\n Evicted Keys: %s\n Slowlog Count: %d\n", size, usage, keys[0], keys[1], countLog)
	// _ = r.getSlowLog(countLog)
	_ = r.getMemoryStats()
	return nil
}

func (r *rediss) Dial(host string) dbcheck.Checker {
	con, err := redis.Dial("tcp", host)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &rediss{con: con}
}

func init() {
	registry.Register("redis", &rediss{})
}
