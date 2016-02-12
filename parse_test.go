package main

import (
	"testing"
)

func TestParseRoutes(t *testing.T) {
	api, err := ParseRoutes("testapp")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(api.Routes) != 3 {
		t.Errorf("Expected to get 3 routes, but got %v", len(api.Routes))
	}
	if api.Routes[0].Name != "get_users" {
		t.Errorf("Expected to get 'get_users' route name, but got: %v", api.Routes[0].Name)
	}
	if api.Routes[0].Path != "/users" {
		t.Errorf("Expected to get /users path, but got: %v", api.Routes[0].Path)
	}
	if api.Routes[0].Collection != true {
		t.Errorf("Expected to get path collection to be true, but got false")
	}
	if api.Routes[0].Versions[0].Schema != nil {
		t.Errorf("Expected to get no route schema, but got")
	}
	if api.Routes[1].Versions[0].Schema == nil {
		t.Errorf("Expected to get route schema, but got nil")
	}
	if api.Version != 5 {
		t.Errorf("Expected to get api version 5, but got %v", api.Version)
	}
	if api.MinVersion != 3 {
		t.Errorf("Expected to get api version 3, but got %v", api.MinVersion)
	}
	if len(api.DeprecatedVersions) != 3 {
		t.Errorf("Expected to get 3 deprecated versions, but got %v", len(api.DeprecatedVersions))
	}
	if api.Routes[0].Versions[4] == nil {
		t.Error("Expected to get version 4 of get_users route, but got nil")
	}
	if api.Routes[1].Versions[4] == nil || api.Routes[1].Versions[4].Schema == nil {
		t.Errorf("Expected to get create_user route v4 schema, but got nil")
	}
}
