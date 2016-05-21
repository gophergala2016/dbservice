package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
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
			PropagateSchemas(route)
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
	route := &Route{
		Versions:        make(map[int]*RouteVersion),
		PluginPipelines: make([]*PluginPipeline, 0),
	}
	pipelines := bytes.Split(line, []byte("|"))
	chunks := bytes.Split(pipelines[0], []byte(","))
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
	if len(pipelines) > 0 {
		for i := 1; i < len(pipelines); i++ {
			pipeline := pipelines[i]
			pluginPipeline, err := ParsePluginPipeline(pipeline)
			if err != nil {
				return nil, err
			}
			route.PluginPipelines = append(route.PluginPipelines, pluginPipeline)
		}
	}
	return route, nil
}

func ParsePluginPipeline(content []byte) (*PluginPipeline, error) {
	content = bytes.TrimSpace(content)
	chunks := bytes.Split(content, []byte(" "))
	if len(chunks) == 0 {
		return nil, errors.New("Plugin name not supplied in pipeline")
	}
	pp := &PluginPipeline{Name: string(chunks[0])}
	if len(chunks) == 1 {
		return pp, nil
	}
	jsonContent := bytes.Join(chunks[1:], []byte(" "))
	arg := make(map[string]interface{})
	err := json.Unmarshal(jsonContent, &arg)
	if err != nil {
		return nil, err
	}
	pp.Argument = arg
	return pp, nil
}
func ParseSchema(path string, route *Route) error {
	files, err := filepath.Glob(path + "/schemas/" + route.Name + ".v[0-9]*.schema")
	if err != nil {
		return nil
	}
	for _, file := range files {
		match := versionRegexp.FindAllStringSubmatch(file, -1)
		versionString := match[0][1]
		version, err := strconv.Atoi(versionString)
		if err != nil {
			return err
		}
		err = ParseSchemaVersion(route, file, version)
		if err != nil {
			return err
		}
	}
	defaultPath := path + "/schemas/" + route.Name + ".schema"
	if _, err := os.Stat(defaultPath); err == nil {
		files = append(files, defaultPath)
		err = ParseSchemaVersion(route, defaultPath, 0)
		if err != nil {
			return err
		}
	}
	return nil
}

func ParseSchemaVersion(route *Route, path string, version int) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	schema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(string(content)))
	if err != nil {
		return err
	}
	if route.Versions[version] == nil {
		route.Versions[version] = &RouteVersion{Version: version}
	}
	route.Versions[version].Schema = schema
	return nil
}

var versionRegexp = regexp.MustCompile(".v([0-9]*).(sql|schema)$")

func ParseSqlTemplate(path string, route *Route) error {
	files, err := filepath.Glob(path + "/sql/" + route.Name + ".v[0-9]*.sql")
	if err != nil {
		return err
	}
	for _, file := range files {
		match := versionRegexp.FindAllStringSubmatch(file, -1)
		versionString := match[0][1]
		version, err := strconv.Atoi(versionString)
		if err != nil {
			return err
		}
		err = ParseSqlTemplateVersion(route, file, version)
		if err != nil {
			return err
		}
	}
	defaultPath := path + "/sql/" + route.Name + ".sql"
	if _, err := os.Stat(defaultPath); err == nil {
		files = append(files, defaultPath)
		err = ParseSqlTemplateVersion(route, defaultPath, 0)
		if err != nil {
			return err
		}
	}
	if len(files) == 0 {
		return fmt.Errorf("%v path is missing sql template", route.Name)
	}
	return nil
}

func ParseSqlTemplateVersion(route *Route, path string, version int) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("missing sql template", path)
	}
	tmpl, err := makeTemplate(string(bytes.TrimSpace(content)))
	if err != nil {
		return err
	}
	if route.Versions[version] == nil {
		route.Versions[version] = &RouteVersion{Version: version}
	}
	route.Versions[version].SqlTemplate = tmpl
	return nil
}

func PropagateSchemas(route *Route) {
	versions := make([]int, 0, len(route.Versions))
	for version := range route.Versions {
		if version != 0 {
			versions = append(versions, version)
		}
	}
	sort.Ints(versions)
	var schema *gojsonschema.Schema
	if route.Versions[0] != nil {
		schema = route.Versions[0].Schema
	}
	for i := len(versions) - 1; i >= 0; i-- {
		if route.Versions[versions[i]].Schema == nil {
			route.Versions[versions[i]].Schema = schema
		} else {
			schema = route.Versions[versions[i]].Schema
		}
	}
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
