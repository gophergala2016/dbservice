package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func handler(route *Route) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		params := make(map[string]interface{})
		sql, err := route.Sql(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, sql)
	}
}

func main() {
	routes, err := ParseRoutes(".")
	if err != nil {
		log.Fatal(err)
	}
	router := httprouter.New()
	for _, route := range routes {
		if route.Method == "GET" {
			router.GET(route.Path, handler(route))
		}
		if route.Method == "POST" {
			router.POST(route.Path, handler(route))
		}
		if route.Method == "PUT" {
			router.PUT(route.Path, handler(route))
		}
		if route.Method == "DELETE" {
			router.DELETE(route.Path, handler(route))
		}
	}
	log.Fatal(http.ListenAndServe(":8080", router))
}
