package cassandra

import (
	"fmt"
	"os/exec"

	"github.com/bayupermadi/dbcheck"
	"github.com/bayupermadi/dbcheck/registry"
	"github.com/gocql/gocql"
)

type cassandra struct {
	session *gocql.Session
}

func (c *cassandra) Version() error {
	q := c.session.Query("select cql_version from system.local;").Iter()
	var cqlVersion string
	q.Scan(&cqlVersion)

	fmt.Printf("Cassandra version %s CQL version %s \n", q.Host().Version(), cqlVersion)
	return nil
}

func (c *cassandra) ActiveClient() error {
	fmt.Println("Cassandra does not provided active_client checking.")
	return nil
}

func (c *cassandra) Health() error {
	//getting information from nodetool command line
	fmt.Println("health information: ")
	output, err := exec.Command("nodetool", "info").CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

func (c *cassandra) Dial(host string) dbcheck.Checker {
	cluster := gocql.NewCluster(host)
	if cluster == nil {
		fmt.Println("no cluster")
		return nil
	}

	session, err := cluster.CreateSession()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &cassandra{session: session}
}

func init() {
	registry.Register("cassandra", &cassandra{})
}
