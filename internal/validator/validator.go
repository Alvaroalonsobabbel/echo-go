package validator

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

func Validate(body string) error {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	schemaPath := filepath.Join(basepath, "..", "schemas", "endpoint_schema.json")
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + schemaPath)
	documentLoader := gojsonschema.NewStringLoader(body)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("error validating schema: %s", err)
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
