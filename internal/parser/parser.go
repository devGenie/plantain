package parser

import tfjson "github.com/hashicorp/terraform-json"

type Parser struct {
	Plan tfjson.Plan
}

func (ingester *Parser) Parse(Plan tfjson.Plan) {

}
