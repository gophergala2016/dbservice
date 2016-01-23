package main

import (
	"bytes"
	"github.com/xeipuuv/gojsonschema"
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

func (self *Route) Sql(params map[string]interface{}) (string, error) {
	var out bytes.Buffer
	if !self.Custom {
		if self.Collection {
			out.Write([]byte("select array_to_json(array_agg(row_to_json(t))) as value from ("))
		} else {
			out.Write([]byte("select row_to_json(t) as value from ("))
		}
	}
	err := self.SqlTemplate.Execute(&out, params)
	if err != nil {
		return "", err
	}
	if self.Custom {
		return out.String(), nil
	}
	out.Write([]byte(") t"))
	return out.String(), nil
}
