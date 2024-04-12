package validator

import (
	"fmt"
	"log"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

func Validate(body string) error {
	schemaLoader := gojsonschema.NewReferenceLoader("file://internal/validator/endpoint_schema.json")
	documentLoader := gojsonschema.NewStringLoader(body)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		log.Fatalf("error validating schema: %v", err)
	}

	if result.Valid() {
		return nil
	} else {
		var errors []string
		for _, desc := range result.Errors() {
			errors = append(errors, desc.Field()+": "+desc.Description())
		}
		return fmt.Errorf(strings.Join(errors, ", "))
	}
}
