package wa

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"

	mgo "github.com/globalsign/mgo"
)

func ConnectionCluster() (*mgo.Session, error, string) {
	tlsConfig := &tls.Config{}
	dialInfo := &mgo.DialInfo{
		Addrs: []string{
			os.Getenv("DB_REPLICA_1"),
			os.Getenv("DB_REPLICA_2"),
			os.Getenv("DB_REPLICA_3"),
		},
		Database: os.Getenv("DB_NAME"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Source:   "admin",
	}

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)

		if err != nil {
			fmt.Println("Some error happens:", err)
			panic("Couldn't establish connection DB")
		}

		return conn, err
	}

	session, err := mgo.DialWithInfo(dialInfo)

	return session, err, dialInfo.Database
}

func ConnectionSingle() (*mgo.Session, error, string) {
	session, err := mgo.Dial(os.Getenv("DB_URL"))
	return session, err, os.Getenv("DB_NAME")
}

// ConnectionOpen Open connection into mongoDB and returning the session
func ConnectionOpen() (*mgo.Session, string) {
	env := os.Getenv("ENV")

	var session *mgo.Session
	var err error
	var dbName string
	if env == "production" {
		session, err, dbName = ConnectionCluster()
	} else {
		session, err, dbName = ConnectionSingle()
	}

	if err != nil {
		fmt.Println("Some error happens:", err)
		panic("Couldn't establish connection DB")
	}

	session.SetMode(mgo.Monotonic, true)

	return session, dbName
}

// ConnectionClose Close the connection mongodb with supplying session to close
func ConnectionClose(session *mgo.Session) {
	session.Close()
}

func DBName() string {
	return os.Getenv("DB_NAME")
}
