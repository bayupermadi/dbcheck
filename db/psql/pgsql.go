package psql

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/bayupermadi/dbcheck"
	"github.com/bayupermadi/dbcheck/registry"
	"github.com/lib/pq"
)

type psql struct {
	DB     *sql.DB
	DBName string
}

func (p *psql) Version() error {
	var version string
	_ = p.DB.QueryRow("SELECT version()").Scan(&version)

	fmt.Println(version)
	return nil
}

func (p *psql) ActiveClient() error {
	var count int
	var host string
	err := p.DB.QueryRow("SELECT count(0) FROM pg_stat_activity where state='active' ").Scan(&count)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = p.DB.QueryRow("select inet_server_addr()").Scan(&host)
	if err != nil {
		fmt.Println(err)
		return err
	}

	info := fmt.Sprintf("active_client(s): %d", count)
	fmt.Println(info)

	// if viper.GetBool("app.aws.service.cloudwatch.enabled") == true {
	// 	tools.CW("PostgreSQL", "Hosts", "Count", float64(count), "Active Connection", host)
	// }

	return nil
}

func (p *psql) Health() error {
	var size int

	err := p.DB.QueryRow("select pg_database_size('" + p.DBName + "') as size;").Scan(&size)
	if err != nil {
		return err
	}

	info := fmt.Sprintf("health_status: \n Database Information \n DBName: %s\t DBSize: %d\n", p.DBName, size)
	fmt.Print(info)

	// if viper.GetBool("app.aws.service.cloudwatch.enabled") == true {
	// 	tools.CW("PostgreSQL", "DB Name", "Megabytes", float64(size), "DB Size", p.DBName)
	// }

	if err := p.getTableSize(); err != nil {
		return err
	}

	return nil
}

func (p *psql) getTables() (map[string][]string, error) {
	rows, err := p.DB.Query("select schemaname as schema, relname as table from pg_statio_all_tables where schemaname not in ('pg_catalog', 'pg_toast', 'information_schema')")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	tables := make(map[string][]string)
	for rows.Next() {
		var schema, table string
		err := rows.Scan(&schema, &table)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		tables[schema] = append(tables[schema], table)
	}

	if len(tables) < 1 {
		return nil, errors.New("Could not find table")
	}

	return tables, nil
}

func (p *psql) getTableSize() error {
	tables, err := p.getTables()
	if err != nil {
		return err
	}
	fmt.Println(" Table Information")
	for k, v := range tables {
		var tableSize, indexSize string
		fmt.Printf("  > Schema: %s\n", k)
		if len(v) < 1 {
			return errors.New("Schema has no table")
		}

		for _, val := range v {
			qTable := fmt.Sprintf("SELECT pg_total_relation_size('%s.%s') as tableSize", k, val)
			qIndex := fmt.Sprintf("SELECT pg_indexes_size('%s.%s') as indexSize", k, val)
			err := p.DB.QueryRow(qTable).Scan(&tableSize)
			if err != nil {
				fmt.Println(err)
				return err
			}
			err = p.DB.QueryRow(qIndex).Scan(&indexSize)
			if err != nil {
				fmt.Println(err)
				return err
			}
			fmt.Printf("     Table: %s\n      Table Size: %s\n      Index Size: %s\n", val, tableSize, indexSize)

			// if viper.GetBool("app.aws.service.cloudwatch.enabled") == true {
			// 	tools.CW("PostgreSQL", "DB Name", "Megabytes", float64(size), "DB Size", p.DBName)
			// }
		}
	}
	return nil
}

func (p *psql) Dial(host string) dbcheck.Checker {
	db, err := sql.Open("postgres", host)
	if err != nil {
		return nil
	}

	str, _ := pq.ParseURL(host)
	arr := strings.Split(str, " ")
	dbname := strings.Split(arr[0], "=")[1]

	return &psql{DB: db, DBName: dbname}
}

func init() {
	registry.Register("postgresql", &psql{})
}
