package main

type Api struct {
	ApiVersion           int
	DeprecatedApiVersion []int
	MinApiVersion        int
	Routes               []*Route
}
