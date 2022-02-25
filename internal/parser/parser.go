package parser

import (
	"github.com/devgenie/plantain/internal/storage"
	tfjson "github.com/hashicorp/terraform-json"
)

type Parser struct {
	Plan tfjson.Plan
	db   *storage.BadgerDB
}

func NewParser() (*Parser, error) {
	db, err := storage.NewBadgerDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	parser := new(Parser)
	parser.db = new(storage.BadgerDB)
	parser.db = db
	return parser, err
}

func (ingester *Parser) Parse(Plan tfjson.Plan) {

}
