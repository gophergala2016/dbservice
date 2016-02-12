package main

import (
	"github.com/xeipuuv/gojsonschema"
	"testing"
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
	tmpl, err := makeTemplate(`select * from users where id={{.id}}`)
	if err != nil {
		t.Errorf("Expected not to get error, but got: %v", err)
	}
	route := &Route{
		Versions: map[int]*RouteVersion{
			0: {
				Version:     0,
				SqlTemplate: tmpl,
				Schema:      schema,
			},
		},
	}
	params := make(map[string]interface{})
	params["id"] = 23
	sql, err := route.Sql(params, 0)
	if err != nil {
		t.Errorf("Expected not to get error, but got: %v", err)
	}
	expected := "with response_table as (select * from users where id=23) select row_to_json(t) as value from (select * from response_table) t"
	if sql != expected {
		t.Errorf("Expected sql:\n%v, but got:\n%v\n", expected, sql)
	}
}
