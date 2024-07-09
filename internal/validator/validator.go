package validator

import (
	_ "embed"
	"fmt"
	"log"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

//go:embed schema.json
var schema string

func Validate(body string) error {
	schemaLoader := gojsonschema.NewStringLoader(schema)
	documentLoader := gojsonschema.NewStringLoader(body)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		log.Fatalf("error validating schema: %v", err)
	}

	if !result.Valid() {
		var errors []string
		for _, desc := range result.Errors() {
			errors = append(errors, desc.Field()+": "+desc.Description())
		}

		return fmt.Errorf(strings.Join(errors, ", "))
	}

	return nil
}
