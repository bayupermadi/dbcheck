package main

import (
	"fmt"

	_ "github.com/bayupermadi/dbcheck/db/bolt"
	_ "github.com/bayupermadi/dbcheck/db/cassandra"
	_ "github.com/bayupermadi/dbcheck/db/mongo"
	_ "github.com/bayupermadi/dbcheck/db/mysql"
	_ "github.com/bayupermadi/dbcheck/db/psql"
	_ "github.com/bayupermadi/dbcheck/db/redis"
	_ "github.com/bayupermadi/dbcheck/db/sqlite"
	"github.com/bayupermadi/dbcheck/registry"
	"github.com/spf13/viper"
	"github.com/wjaoss/aws-wrapper/session"
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
	// db := flag.String("db", "", "Specify your database server. Supported databases (key): redis, mongo, postgresql, mysql, cassandra, bolt, sqlite ")
	// host := flag.String("host", "", "Specify your database connection URI depending your server")
	// path := flag.String("path", "", "Specify your database path (used for bolt and sqlite)")
	// flag.Parse()
	// pathDB := []string{"sqlite", "bolt"}
	// hostDB := []string{"mysql", "postgresql", "mongo", "redis", "cassandra"}
	// for _, v := range pathDB {
	// 	if v == *db && strings.Contains(os.Args[3], "host") {
	// 		fmt.Printf("%s need a argument path (found %s)\n", *db, os.Args[3])
	// 		return
	// 	}
	// }

	// for _, v := range hostDB {
	// 	if v == *db && strings.Contains(os.Args[3], "path") {
	// 		fmt.Printf("%s need a argument host (found %s)\n", *db, os.Args[3])
	// 		return
	// 	}
	// }
	path := ""

	// load configuration file
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	// open connection to aws
	if viper.GetBool("app.aws.enabled") {
		awsKeyID := viper.Get("app.aws.credential.id-key").(string)
		awsSecretKey := viper.Get("app.aws.credential.secret-key").(string)
		awsRegion := viper.Get("app.aws.credential.region").(string)

		session.SetConfiguration(awsKeyID, awsSecretKey, awsRegion)
	}

	pgMon := viper.GetBool("database.pgsql.enabled")

	if pgMon {
		db := "postgresql"
		host := viper.Get("database.pgsql.uri").(string)
		fmt.Println(db, host)
		dbInfo(db, host, path)
	}

}
