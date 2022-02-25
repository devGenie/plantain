package storage

import (
	"bytes"
	"encoding/gob"

	badger "github.com/dgraph-io/badger/v3"
	tfjson "github.com/hashicorp/terraform-json"
)

type BadgerDB struct {
}

func (badgerDB *BadgerDB) openDB() (*badger.DB, error) {
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewBadgerDB() (*BadgerDB, error) {

	badgerDB := new(BadgerDB)
	return badgerDB, nil
}

func (badgerDB *BadgerDB) Read(tfPlanHash string) (*tfjson.Plan, error) {
	var badgerData []byte
	db, err := badgerDB.openDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var planBytes bytes.Buffer
	decorder := gob.NewDecoder(&planBytes)
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(tfPlanHash))
		if err != nil {
			return err
		}
		badgerData, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	_, err = planBytes.Write(badgerData)
	if err != nil {
		return nil, err
	}

	tfPlan := new(tfjson.Plan)
	err = decorder.Decode(tfPlan)
	if err != nil {
		return nil, err
	}
	return tfPlan, nil
}

func (badgerDB *BadgerDB) Write(tfPlanHash string, plan tfjson.Plan) error {
	db, err := badgerDB.openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	var planBytes bytes.Buffer
	byteEncoder := gob.NewEncoder(&planBytes)
	err = byteEncoder.Encode(plan)
	if err != nil {
		return err
	}

	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(tfPlanHash), planBytes.Bytes())
		return err
	})

	return err
}
