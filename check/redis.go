package check

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

type rediss struct {
	Host string
}

func (r rediss) Version() (string, error) {
	con, err := redis.Dial("tcp", r.Host)
	if err != nil {
		return "", err
	}

	defer con.Close()

	version, err := redis.String(con.Do("INFO"))
	if err != nil {
		fmt.Println("getting info", err)
		return "", nil
	}

	return version, nil
}

func NewRedis(host string) VersionChecker {
	return rediss{
		Host: host,
	}
}
