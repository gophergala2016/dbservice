package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"strings"
	"text/template"
)

type Route struct {
	Name        string
	Method      string
	Path        string
	Collection  bool
	Custom      bool
	Schema      *gojsonschema.Schema
	SqlTemplate *template.Template
}

func (self *Route) validate(params map[string]interface{}) (string, error) {
	documentLoader := gojsonschema.NewGoLoader(params)
	result, err := self.Schema.Validate(documentLoader)
	if err != nil {
		return "", err
	}
	if !result.Valid() {
		errors := make(map[string]string)
		for _, resultErr := range result.Errors() {
			errors[resultErr.Field()] = resultErr.Description()
		}
		errorsJson, err := json.Marshal(errors)
		if err != nil {
			return "", err
		}
		return string(errorsJson), nil
	}
	return "", nil
}

func (self *Route) Sql(params map[string]interface{}) (string, error) {
	var out bytes.Buffer
	response, err := self.validate(params)
	if err != nil {
		return "", err
	}
	if response != "" {
		return response, nil
	}
	if !self.Custom {
		if self.Collection {
			out.Write([]byte("select array_to_json(array_agg(row_to_json(t))) as value from ("))
		} else {
			out.Write([]byte("select row_to_json(t) as value from ("))
		}
	}
	err = self.SqlTemplate.Execute(&out, params)
	if err != nil {
		return "", err
	}
	if self.Custom {
		return out.String(), nil
	}
	out.Write([]byte(") t"))
	return out.String(), nil
}

func quoteString(value interface{}) string {
	stringValue := fmt.Sprintf("%v", value)
	stringValue = strings.Replace(stringValue, "'", "''", -1)
	return "'" + stringValue + "'"
}

func makeTemplate(t string) (*template.Template, error) {
	funcMap := template.FuncMap{
		"quote": quoteString,
	}
	return template.New("").Funcs(funcMap).Parse(t)

}
