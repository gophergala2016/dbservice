package main

type Api struct {
	Version           int
	DeprecatedVersion []int
	MinVersion        int
	Routes            []*Route
}
