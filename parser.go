package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"strings"
)

func ParseRoutes(path string) (*Api, error) {
	content, err := ioutil.ReadFile(path + "/routes")
	if err != nil {
		return nil, err
	}
	routes := make([]*Route, 0, 0)
	lines := bytes.Split(content, []byte("\n"))
	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) != 0 {
			route, err := ParseRoute(line)
			if err != nil {
				return nil, err
			}
			err = ParseSchema(path, route)
			if err != nil {
				return nil, err
			}
			err = ParseSqlTemplate(path, route)
			if err != nil {
				return nil, err
			}
			routes = append(routes, route)
		}
	}
	return &Api{Routes: routes}, nil
}

func ParseRoute(line []byte) (*Route, error) {
	route := &Route{}
	chunks := bytes.Split(line, []byte(","))
	urlParams := bytes.Split(chunks[0], []byte(" "))
	route.Method = strings.ToUpper(string(urlParams[0]))
	route.Path = string(urlParams[1])
	for i, chunk := range chunks {
		if i != 0 {
			chunkParts := bytes.Split(chunk, []byte(":"))
			if len(chunkParts) != 2 {
				return nil, fmt.Errorf("unexpected route parameters: %v", string(line))
			}
			name := string(bytes.TrimSpace(chunkParts[0]))
			value := string(bytes.TrimSpace(chunkParts[1]))
			if value[0] == '\'' && value[len(value)-1] == '\'' {
				value = value[1 : len(value)-1]
			}
			if name == "name" {
				route.Name = value
			}
			if name == "collection" && value == "true" {
				route.Collection = true
			}
			if name == "custom" && value == "true" {
				route.Custom = true
			}
		}
	}
	return route, nil
}

func ParseSchema(path string, route *Route) error {
	content, err := ioutil.ReadFile(path + "/schemas/" + route.Name + ".schema")
	if err != nil {
		return nil
	}
	schema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(string(content)))
	if err != nil {
		return err
	}
	route.Schema = schema
	return nil
}

func ParseSqlTemplate(path string, route *Route) error {
	content, err := ioutil.ReadFile(path + "/sql/" + route.Name + ".sql")
	if err != nil {
		return fmt.Errorf("%v path is missing sql template", route.Name)
	}
	tmpl, err := makeTemplate(string(bytes.TrimSpace(content)))
	if err != nil {
		return err
	}
	route.SqlTemplate = tmpl
	return nil
}

type DbConfig struct {
	User     string
	Password string
	Database string
	Host     string
	Port     int
	SslMode  string
}

func ParseDbConfig(path string) (*DbConfig, error) {
	conf := &DbConfig{}
	content, err := ioutil.ReadFile(path + "/config.toml")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error while reading config.toml configuration: %v", err))
	}
	_, err = toml.Decode(string(content), &conf)
	if err != nil {
		return nil, err
	}
	if conf.User == "" {
		conf.User = "postgres"
	}
	if conf.Host == "" {
		conf.Host = "127.0.0.1"
	}
	if conf.Port == 0 {
		conf.Port = 5432
	}
	if conf.SslMode == "" {
		conf.SslMode = "disable"
	}
	return conf, nil

}
