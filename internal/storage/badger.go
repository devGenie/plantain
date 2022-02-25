package storage

import (
	badger "github.com/dgraph-io/badger/v3"
)

type BadgerDB struct {
	db *badger.DB
}

func NewBadgerDB() (*BadgerDB, error) {
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		return nil, err
	}

	badgerDB := new(BadgerDB)
	badgerDB.db = db
	return badgerDB, nil
}

func (badgerDB *BadgerDB) Close() {
	badgerDB.db.Close()
}
