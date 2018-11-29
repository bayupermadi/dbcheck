package mysql

import (
	"database/sql"
	"fmt"

	"github.com/bayupermadi/dbcheck"
	"github.com/bayupermadi/dbcheck/registry"
	_ "github.com/go-sql-driver/mysql"
)

type mysql struct {
	DB *sql.DB
}

type activeClient struct {
	user    string `json:"USER"`
	host    string `json:"HOST"`
	db      string `json:"DB"`
	command string `json:"COMMAND"`
	state   string `json:"STATE"`
	info    string `json:"INFO"`
	time    string `json:"TIME_MS"`
}

func (m *mysql) Version() error {
	rows, err := m.DB.Query("SHOW VARIABLES LIKE '%version%'")
	if err != nil {
		fmt.Println("query ", err)
		return nil
	}

	var version, variable, value string
	version = "MySql "
	for rows.Next() {
		err := rows.Scan(&variable, &value)
		if err != nil {
			fmt.Println("scan", err)
			return err
		}
		version += value + " "
	}

	fmt.Println(version)

	return nil
}

func (m *mysql) ActiveClient() error {
	rows, err := m.DB.Query("SELECT USER, HOST, DB, COMMAND, STATE, INFO, TIME_MS FROM INFORMATION_SCHEMA.PROCESSLIST")
	if err != nil {
		return err
	}

	activeClients := []activeClient{}
	for rows.Next() {
		var c activeClient
		_ = rows.Scan(&c.user, &c.host, &c.db, &c.command, &c.state, &c.info, &c.time)

		activeClients = append(activeClients, c)
	}

	info := fmt.Sprintf("active_client(s): %d \n", len(activeClients))
	for _, v := range activeClients {
		if v.state == "" {
			v.state = "-"
		}
		if v.time == "" {
			v.time = "0"
		}
		if v.info == "" {
			v.info = "-"
		}
		info += " ____________\n User: " + v.user + "\n Host: " + v.host + "\n Command: " + v.command + "\n State: " + v.state + "\n DB Name: " + v.db + "\n Running Queries: " + v.info + "\n Time: " + v.time + "\n"
	}
	fmt.Println(info)
	return nil
}

func (m *mysql) Health() error {
	fmt.Println("health_status: \n")

	return nil
}

func (m *mysql) Dial(host string) dbcheck.Checker {
	db, err := sql.Open("mysql", host)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &mysql{DB: db}
}

func init() {
	registry.Register("mysql", &mysql{})
}
