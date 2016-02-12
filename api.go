package main

type Api struct {
	Version            int
	DeprecatedVersions []int
	MinVersion         int
	Routes             []*Route
}
