package main

type Api struct {
	Version            int
	DeprecatedVersions []int
	MinVersion         int
	Routes             []*Route
}

func (self *Api) IsDeprecated(version int) bool {
	for _, deprecatedVersion := range self.DeprecatedVersions {
		if version == deprecatedVersion {
			return true
		}
	}
	return false
}
