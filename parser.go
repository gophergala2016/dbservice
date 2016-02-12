package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"strconv"
	"strings"
)

func ParseRoutes(path string) (*Api, error) {
	content, err := ioutil.ReadFile(path + "/routes")
	if err != nil {
		return nil, err
	}
	api := &Api{Routes: make([]*Route, 0, 0)}
	lines := bytes.Split(content, []byte("\n"))
	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) != 0 {
			ok, err := ParseApiSettings(api, line)
			if ok {
				continue
			}
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
			api.Routes = append(api.Routes, route)
		}
	}
	return api, nil
}

func ParseApiSettings(api *Api, line []byte) (bool, error) {
	var err error
	if bytes.HasPrefix(line, []byte("api_version")) {
		line = bytes.TrimPrefix(line, []byte("api_version:"))
		line = bytes.TrimSpace(line)
		api.Version, err = strconv.Atoi(string(line))
		if err != nil {
			return false, err
		}
		return true, nil
	} else if bytes.HasPrefix(line, []byte("deprecated_api_version")) {
		line = bytes.TrimPrefix(line, []byte("deprecated_api_version:"))
		line = bytes.TrimSpace(line)
		line = bytes.TrimPrefix(line, []byte("["))
		line = bytes.TrimSuffix(line, []byte("]"))
		chunks := bytes.Split(line, []byte(","))
		versions := make([]int, 0, 0)
		for _, chunk := range chunks {
			if bytes.Contains(chunk, []byte("-")) {
				rng := bytes.Split(chunk, []byte("-"))
				if len(rng) != 2 {
					return false, fmt.Errorf("Expected to get deprecated api version range, but got: %v", chunk)
				}
				from, err := strconv.Atoi(string(rng[0]))
				if err != nil {
					return false, err
				}
				to, err := strconv.Atoi(string(rng[1]))
				if err != nil {
					return false, err
				}
				if from > to {
					return false, fmt.Errorf("Got invalid deprecated api version range: %v", chunk)
				}
				if from == to {
					versions = append(versions, from)
				} else {
					for version := from; version <= to; version++ {
						versions = append(versions, version)
					}
				}
			} else {
				version, err := strconv.Atoi(string(chunk))
				if err != nil {
					return false, err
				}
				versions = append(versions, version)
			}
			api.DeprecatedVersions = versions
		}
		return true, nil
	} else if bytes.HasPrefix(line, []byte("min_api_version")) {
		line = bytes.TrimPrefix(line, []byte("min_api_version:"))
		line = bytes.TrimSpace(line)
		api.MinVersion, err = strconv.Atoi(string(line))
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
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
