package sqlite

import (
	"database/sql"
	"fmt"

	//sqlite3 driver
	"github.com/bayupermadi/dbcheck"
	"github.com/bayupermadi/dbcheck/registry"
	_ "github.com/mattn/go-sqlite3"
)

type sqlite struct {
	db *sql.DB
}

func (s *sqlite) Version() error {
	var version string
	err := s.db.QueryRow("select sqlite_version() as version;").Scan(&version)
	if err != nil {
		return err
	}
	fmt.Printf("SQLite version %s\n", version)
	return nil
}

func (s *sqlite) ActiveClient() error {
	fmt.Println("Sqlite does not provided active_client checking")
	return nil
}

func (s *sqlite) Health() error {
	var pageSize, pageCount int
	if err := s.db.QueryRow("select page_size as pageSize, page_count as pageCount from pragma_page_size, pragma_page_count;").Scan(&pageSize, &pageCount); err != nil {
		return err
	}
	fmt.Printf("health_status: \n pragma_page_size: %d\n pragma_page_count: %d\n", pageSize, pageCount)
	return nil
}

func (s *sqlite) Dial(path string) dbcheck.Checker {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &sqlite{db: db}
}

func init() {
	registry.Register("sqlite", &sqlite{})
}
