package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
)

func ParseRoutes(path string) ([]*Route, error) {
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
			routes = append(routes, route)
		}
	}
	return routes, nil
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
