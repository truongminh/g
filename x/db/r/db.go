package r

import (
	"github.com/golang/glog"
	r "gopkg.in/dancannon/gorethink.v2"
)

type RethinkDB struct {
	Session *r.Session
	dbName  string
}

func NewRethinkDB(address string, dbName string, tag string) *RethinkDB {
	var session, err = r.Connect(r.ConnectOpts{Address: address})
	if err != nil {
		glog.Fatal(err)
	}
	glog.Infof("RethinkDB online %v/%v for %v", address, dbName, tag)
	return &RethinkDB{
		Session: session,
		dbName:  dbName,
	}
}

func (db *RethinkDB) Table(name string) r.Term {
	return r.DB(db.dbName).Table(name)
}

func (db *RethinkDB) QueryBuilder() r.Term {
	if db == nil {
		glog.Fatal("nil db instance")
	}
	return r.DB(db.dbName)
}

func (db *RethinkDB) IsErrEmpty(err error) bool {
	return err == r.ErrEmptyResult
}

func (db *RethinkDB) MustGetTable(name string) r.Term {
	var cursor, err = r.DB(db.dbName).TableList().Run(db.Session)
	if err != nil {
		panic(err)
	}

	var names []string
	if err := cursor.All(&names); err != nil {
		panic(err)
	}

	for _, table := range names {
		if table == name {
			return db.Table(name)
		}
	}

	// create table
	{
		var _, err = r.DB(db.dbName).TableCreate(name).RunWrite(db.Session)
		if err != nil {
			panic(err)
		}
		return db.Table(name)
	}

}
