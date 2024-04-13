package validator

import (
	"fmt"
	"log"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

const schema = `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "data": {
      "type": "object",
      "properties": {
        "type": {
          "type": "string",
          "const": "endpoints"
        },
        "attributes": {
          "type": "object",
          "properties": {
            "verb": {
              "type": "string",
              "enum": ["GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"]
            },
            "path": {
              "type": "string",
              "pattern": "^\\/.*"
            },
            "response": {
              "type": "object",
              "properties": {
                "code": {
                  "type": "integer",
                  "minimum": 100,
                  "maximum": 599
                },
                "headers": {
                  "type": "object",
                  "additionalProperties": {
                    "type": "string"
                  }
                },
                "body": {
                  "type": "string"
                }
              },
              "required": ["code"]
            }
          },
          "required": ["verb", "path", "response"]
        }
      },
      "required": ["type", "attributes"]
    }
  },
  "required": ["data"]
}`

func Validate(body string) error {
	schemaLoader := gojsonschema.NewStringLoader(schema)
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
