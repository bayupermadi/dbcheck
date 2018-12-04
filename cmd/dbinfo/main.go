package main

import (
	"fmt"
	"time"

	_ "github.com/bayupermadi/dbcheck/db/bolt"
	_ "github.com/bayupermadi/dbcheck/db/cassandra"
	_ "github.com/bayupermadi/dbcheck/db/mongo"
	_ "github.com/bayupermadi/dbcheck/db/mysql"
	_ "github.com/bayupermadi/dbcheck/db/psql"
	_ "github.com/bayupermadi/dbcheck/db/redis"
	_ "github.com/bayupermadi/dbcheck/db/sqlite"
	"github.com/bayupermadi/dbcheck/registry"
	"github.com/spf13/viper"
)

func dbInfo(db string, host string, path string) {
	var dial string
	if host == "" {
		dial = path
	} else {
		dial = host
	}
	dialer := registry.Dialers(db)
	if dialer == nil {
		fmt.Printf("(%s) Database not supported\n", db)
		return
	}
	checker := dialer.Dial(dial)
	if checker == nil {
		fmt.Println("Server unreachable")
		return
	}
	if err := checker.Version(); err != nil {
		fmt.Println(err)
		return
	}
	if err := checker.ActiveClient(); err != nil {
		fmt.Println(err)
		return
	}
	if err := checker.Health(); err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	path := ""

	// load configuration file
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	pgMon := viper.GetBool("database.pgsql.enabled")

	if pgMon {
		db := "postgresql"
		host := viper.Get("database.pgsql.uri").(string)
		for {
			dbInfo(db, host, path)

			<-time.After(time.Second * 30)
		}
	}

}
