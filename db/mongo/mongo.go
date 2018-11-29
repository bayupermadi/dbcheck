package mongo

import (
	"fmt"

	"github.com/bayupermadi/dbcheck"
	"github.com/bayupermadi/dbcheck/registry"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type mongo struct {
	session *mgo.Session
}

type mongoHealth struct {
	db          string `bson:"db""`
	collections int    `bson:"collections"`
	views       int    `bson:"views"`
	objects     int    `bson:"objects"`
	avgObjSize  int    `bson:"avgObjSize"`
	dataSize    int    `bson:"dataSize"`
	storageSize int    `bson:"storageSize"`
	numExtents  int    `bson:"numExtents"`
	indexes     int    `bson:"indexes"`
	indexSize   int    `bson:"indexSize"`
	fsUsedSize  int    `bson:"fsUsedSize"`
	fsTotalSize int    `bson:"fsTotalSize"`
}

func (m *mongo) Version() error {
	buildInfo, err := m.session.BuildInfo()
	if err != nil {
		fmt.Println("getting build info", err)
		return err
	}

	version := fmt.Sprintf(" MongoDB\n db version %s \n git version %s \n", buildInfo.Version, buildInfo.GitVersion)

	fmt.Println(version)
	return nil
}

func (m *mongo) ActiveClient() error {
	sessionCP := m.session.Copy()
	defer sessionCP.Close()

	var res bson.M
	if err := sessionCP.DB("test").Run("serverStatus", &res); err != nil {
		return err
	}

	client := res["globalLock"].(bson.M)["activeClients"].(bson.M)["total"]

	fmt.Printf(" active_client(s): %d \n", client)

	return nil
}

func (m *mongo) Health() error {
	sessionCP := m.session.Copy()
	defer sessionCP.Close()

	res := bson.M{}
	if err := sessionCP.DB("test").Run("dbstats", &res); err != nil {
		return err
	}

	fmt.Printf(" health_status: \n")
	fmt.Printf("  DB: %s\n  Collection: %d\n  Storage Size: %f\n  Indexes: %d\n  Data Size: %f\n", res["db"], res["collections"], res["storageSize"], res["indexes"], res["dataSize"])

	return nil
}

func (m *mongo) Dial(host string) dbcheck.Checker {
	session, err := mgo.Dial(host)
	if err != nil {
		fmt.Println("mongo conn", err)
		return nil
	}

	return &mongo{session: session}
}

func init() {
	registry.Register("mongo", &mongo{})
}
