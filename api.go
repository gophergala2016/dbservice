package main

import (
	"github.com/gophergala2016/dbserver/plugins"
	"os"
)

type Api struct {
	Version            int
	DeprecatedVersions []int
	MinVersion         int
	Routes             []*Route
	Plugins            map[string]Plugin
}

func (self *Api) IsDeprecated(version int) bool {
	for _, deprecatedVersion := range self.DeprecatedVersions {
		if version == deprecatedVersion {
			return true
		}
	}
	return false
}

func (self *Api) RegisterPlugin(name string, plugin Plugin) {
	if _, err := os.Stat("plugins/" + name + ".toml"); err != nil {
		return
	}
	plugin.ParseConfig("plugins/" + name + ".toml")
	self.Plugins[name] = plugin
}

func (self *Api) GetPlugin(name string) Plugin {
	return self.Plugins[name]
}

type Plugin interface {
	ParseConfig(path string) error
	Process(data map[string]interface{}, arg map[string]interface{}) *plugins.Response
}
