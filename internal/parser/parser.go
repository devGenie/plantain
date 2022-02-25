package parser

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"

	"github.com/devgenie/plantain/internal/storage"
	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
)

type Parser struct {
	Plan tfjson.Plan
	db   *storage.BadgerDB
	tf   *tfexec.Terraform
}

func NewParser(tf *tfexec.Terraform) (*Parser, error) {
	db, err := storage.NewBadgerDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

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

	pln, err := parser.tf.ShowPlanFile(context.Background(), planFilePath)
	if err != nil {
		return err
	}
	log.Println(pln.FormatVersion)
	return nil
}
