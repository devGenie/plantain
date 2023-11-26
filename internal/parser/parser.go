package parser

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"

	"github.com/devgenie/plantain/internal/storage"
	"github.com/dgraph-io/badger/v3"
	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
)

type Parser struct {
	plan      tfjson.Plan
	db        *storage.BadgerDB
	tf        *tfexec.Terraform
	toDelete  []*tfjson.ResourceChange
	toCreate  []*tfjson.ResourceChange
	toReplace []*tfjson.ResourceChange
	toUpdate  []*tfjson.ResourceChange
}

func NewParser(tf *tfexec.Terraform) (*Parser, error) {
	db, err := storage.NewBadgerDB()
	if err != nil {
		return nil, err
	}

	parser := new(Parser)
	parser.db = new(storage.BadgerDB)
	parser.db = db
	parser.tf = tf
	return parser, err
}

func (parser *Parser) Parse(planFilePath string) error {
	//get checksum of a plan file and compare it to previously indexed plans
	planFile, err := os.Open(planFilePath)
	if err != nil {
		return err
	}
	defer planFile.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, planFile); err != nil {
		return err
	}

	sha256Hash := hash.Sum(nil)
	checksum := hex.EncodeToString(sha256Hash)
	log.Println("Plan file checksum", checksum)

	badgerData, err := parser.db.Read(checksum)
	if err == badger.ErrKeyNotFound {
		pln, err := parser.tf.ShowPlanFile(context.Background(), planFilePath)
		if err != nil {
			return err
		}
		err = parser.db.Write(checksum, *pln)
		if err != nil {
			return err
		}
		badgerData, err = parser.db.Read(checksum)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	parser.plan = *badgerData
	parser.serialize(&parser.plan)
	return nil
}

func (parser *Parser) serialize(plan *tfjson.Plan) {
	log.Printf("Serializing %d resources \n", len(plan.ResourceChanges))
	for i := 0; i < len(plan.ResourceChanges); i++ {
		resource := plan.ResourceChanges[i]
		log.Println("serializing:", resource.Address)
		actions := resource.Change.Actions
		switch {
		case actions.Create():
			parser.toCreate = append(parser.toCreate, resource)
		case actions.Delete():
			parser.toDelete = append(parser.toDelete, resource)
		case actions.Replace():
			parser.toReplace = append(parser.toReplace, resource)
		case actions.Update():
			parser.toUpdate = append(parser.toUpdate, resource)
		}
	}
	log.Printf("Successfully serialized %d resource(s) \n", len(plan.ResourceChanges))
}
