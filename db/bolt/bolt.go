package bolt

import (
	"fmt"
	"time"

	"github.com/bayupermadi/dbcheck"
	"github.com/bayupermadi/dbcheck/registry"
	"github.com/boltdb/bolt"
)

type boltdb struct {
	db *bolt.DB
}

func (b *boltdb) Version() error {
	fmt.Println("Bolt does not support version checking.")
	return nil
}

func (b *boltdb) ActiveClient() error {
	dbStats := b.db.Stats()
	fmt.Printf("transaction\n Started transaction: %d \n Open connection: %d\n", dbStats.TxN, dbStats.OpenTxN)

	return nil
}

func (b *boltdb) Health() error {
	tx, err := b.db.Begin(false)
	if err != nil {
		return err
	}
	cursor := tx.Cursor()
	bucket := cursor.Bucket()
	stats := bucket.Stats()
	fmt.Printf("health status\n Number of bucket: %d\n Total Keys: %d\n", stats.BucketN, stats.KeyN)
	return nil
}

func (b *boltdb) Dial(host string) dbcheck.Checker {
	db, err := bolt.Open(host, 0600, &bolt.Options{Timeout: time.Second * 5})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &boltdb{db: db}
}

func init() {
	registry.Register("bolt", &boltdb{})
}
