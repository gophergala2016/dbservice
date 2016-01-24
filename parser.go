package main

import (
	"bytes"
	"io/ioutil"
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
	return &Route{}, nil
}
