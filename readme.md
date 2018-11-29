# Database Checking
A tool for database checking. Version checking, active client checking and health checking.
## Getting Started
### Prerequisites
Before use this tool, you need some packages installed on your system.
- [Mysql Driver](http://github.com/go-sql-driver/mysql)
- [Postgresql Driver](https://github.com/lib/pq)
- [Redigo-redis](https://github.com/gomodule/redigo)
- [Mgo](https://github.com/globalsign/mgo)
- [GoCQL](https://github.com/gocql/gocql)
- [Bolt](https://github.com/boltdb/bolt)
- [SQLite Driver](https://github.com/mattn/go-sqlite3)

### Using this tool
Get this tool by execute this command
```
go get github.com/bayupermadi/dbcheck
```

With go, you can simply get all required packages by its tool with a single command.
First go to root project:
```
cd $(go env GOPATH)/src/github.com/bayupermadi/dbcheck

```
then hit: 
```
go get -v -u
```

To build this tool, you can follow this command below

```
cd $(go env GOPATH)/src/github.com/bayupermadi/dbcheck/cmd/dbinfo
go build
```

### Available Commands
> dbname: `mysql`, `postgresql`, `redis`, `mongo`, `cassandra`, `bolt`, `sqlite`

Command for Mysql, Postgresql, Redis, Mongo, Cassandra
```
./dbinfo --db [dbname] --host [connection_string]
```

Command for Bolt and SQLite
```
./dbinfo --db [dbname] --host [path_to_your_db]
```
