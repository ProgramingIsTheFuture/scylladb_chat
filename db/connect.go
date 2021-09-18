package db

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
)

var Session gocqlx.Session

func ConnectScylladb(host string) gocqlx.Session {
	cluster := gocql.NewCluster(host)

	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		panic("Failed wrapping the session")
	}

	err = createKeySpace(session)
	if err != nil {
		panic("Failed creating the keyspace")
	}

	err = createTables(session)
	if err != nil {
		panic("Failed creating the tables")
	}

	Session = session

	fmt.Println("Db connected")

	return session
}

func createKeySpace(session gocqlx.Session) error {
	err := session.ExecStmt(`
	CREATE KEYSPACE IF NOT EXISTS pixelchart 
		with replication = {'class': 'SimpleStrategy', 'replication_factor': 1}
	`)

	if err != nil {
		return err
	}

	return nil
}

func createTables(session gocqlx.Session) error {

	err := session.ExecStmt(`
		CREATE TABLE IF NOT EXISTS pixelchart.message (
			id uuid PRIMARY KEY,	
			message text,
			sender text,
			room uuid,
		)
	`)

	if err != nil {
		return err
	}

	err = session.ExecStmt(`
		CREATE TABLE IF NOT EXISTS pixelchart.room (
			id uuid PRIMARY KEY,	
			name text,
		)
	`)

	if err != nil {
		return err
	}

	return nil
}
