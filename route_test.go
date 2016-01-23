package main

import (
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"testing"
	"text/template"
)

func TestRouteSql(t *testing.T) {
	schemaContent := `
{
  "title": "Get user",
  "description": "Get user",
  "type": "object",
  "properties": {
    "id": {
      "description": "User id",
      "type": "integer"
    }
  },
  "required": ["id"]
}
`
	schema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(string(schemaContent)))
	if err != nil {
		t.Errorf("Expected not to get error, but got: %v", err)
	}
	tmpl, err := template.New("").Parse(`select * from users where id={{.id}}`)
	if err != nil {
		t.Errorf("Expected not to get error, but got: %v", err)
	}
	route := &Route{
		SqlTemplate: tmpl,
		Schema:      schema,
	}
	params := make(map[string]interface{})
	params["id"] = 23
	sql, err := route.Sql(params)
	if err != nil {
		t.Errorf("Expected not to get error, but got: %v", err)
	}
	fmt.Println(sql)
}
