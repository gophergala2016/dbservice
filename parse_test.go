package main

import (
	"testing"
)

func TestParseRoutes(t *testing.T) {
	routes, err := ParseRoutes("testapp")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(routes) != 3 {
		t.Errorf("Expected to get 3 routes, but got %v", len(routes))
	}
	if routes[0].Name != "get_users" {
		t.Errorf("Expected to get 'get_users' route name, but got: %v", routes[0].Name)
	}
	if routes[0].Path != "/users" {
		t.Errorf("Expected to get /users path, but got: %v", routes[0].Path)
	}
	if routes[0].Collection != true {
		t.Errorf("Expected to get path collection to be true, but got false")
	}
}
