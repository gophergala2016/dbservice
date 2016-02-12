package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"strings"
	"text/template"
)

type Route struct {
	Name       string
	Method     string
	Path       string
	Collection bool
	Custom     bool
	Versions   map[int]*RouteVersion
}

type RouteVersion struct {
	Version     int
	Schema      *gojsonschema.Schema
	SqlTemplate *template.Template
}

func (self *Route) validate(params map[string]interface{}, version int) (string, error) {
	route := self.Versions[version]
	if route == nil {
		return "", fmt.Errorf("Route version %v missing from %v route", version, self.Name)
	}
	if route.Schema == nil {
		return "", nil
	}
	documentLoader := gojsonschema.NewGoLoader(params)
	result, err := route.Schema.Validate(documentLoader)
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

func (self *Route) Sql(params map[string]interface{}, version int) (string, error) {
	route := self.Versions[version]
	if route == nil {
		return "", fmt.Errorf("Route version %v missing from %v route", version, self.Name)
	}
	var out bytes.Buffer
	response, err := self.validate(params, version)
	if err != nil {
		return "", err
	}
	if response != "" {
		return response, errors.New("schema validation failed")
	}
	if !self.Custom {
		out.Write([]byte("with response_table as ("))
	}
	err = route.SqlTemplate.Execute(&out, params)
	if err != nil {
		return "", err
	}
	if self.Custom {
		return out.String(), nil
	}
	if self.Collection {
		out.Write([]byte(") select array_to_json(array_agg(row_to_json(t))) as value from (select * from response_table) t"))
	} else {
		out.Write([]byte(") select row_to_json(t) as value from (select * from response_table) t"))
	}
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
