package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func getRequestParams(r *http.Request, urlParams map[string]interface{}) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}
	for k, v := range r.Form {
		if len(v) >= 1 {
			params[k] = v[0]
		}
	}
	if r.Header.Get("Content-Type") == "application/json" {
		decoder := json.NewDecoder(r.Body)
		requestBodyMap := make(map[string]interface{})
		err = decoder.Decode(&requestBodyMap)
		if err != nil {
			return nil, err
		}
		for k, v := range requestBodyMap {
			params[k] = v
		}
	}
	for k, v := range urlParams {
		params[k] = v
	}
	return params, nil
}

func handler(route *Route) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		urlParams := make(map[string]interface{})
		for _, urlParam := range ps {
			urlParams[urlParam.Key] = urlParam.Value
		}
		params, err := getRequestParams(r, urlParams)
		sql, err := route.Sql(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, sql)
	}
}

var db *sql.DB

func main() {
	routes, err := ParseRoutes(".")
	if err != nil {
		log.Fatal(err)
	}
	db, err = GetDbConnection()
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
