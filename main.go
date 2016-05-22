package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gophergala2016/dbserver/plugins/jwt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
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

var apiVersionRegexp = regexp.MustCompile(".v([0-9]*)\\+json$")

func handler(api *Api, route *Route, version int) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var err error
		apiVersion := version
		headerVersion := r.Header.Get("api-version")
		if apiVersion == 0 && headerVersion != "" {
			apiVersion, err = strconv.Atoi(headerVersion)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Unknown api version: %v", r.Header.Get("api-version"))
				return
			}
		}
		if apiVersion == 0 && r.Header.Get("Accept") != "" {
			matches := apiVersionRegexp.FindAllStringSubmatch(r.Header.Get("Accept"), -1)
			if len(matches) > 0 && len(matches[0]) > 1 {
				apiVersion, err = strconv.Atoi(matches[0][1])
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Println("Unknown api version: %v", matches[0][1])
				}
			}
		}
		if apiVersion == 0 {
			apiVersion = api.Version
		}
		urlParams := make(map[string]interface{})
		for _, urlParam := range ps {
			urlParams[urlParam.Key] = urlParam.Value
		}
		params, err := getRequestParams(r, urlParams)
		sql, err := route.Sql(params, apiVersion)
		if err != nil && sql != "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, sql)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		log.Println(sql)
		rows, err := db.Query(sql)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		defer rows.Close()
		var jsonValue string
		w.Header().Set("X-Api-Version", strconv.Itoa(apiVersion))
		if api.IsDeprecated(apiVersion) {
			w.Header().Set("X-Api-Deprecated", "true")
		}
		if rows.Next() {
			err := rows.Scan(&jsonValue)
			if err != nil {
				if route.Collection {
					jsonValue = "[]"
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
				return
			}
		}
		if len(route.PluginPipelines) > 0 {
			goThroughPipelines(api, jsonValue, route.PluginPipelines, w)
		} else {
			fmt.Fprint(w, jsonValue)
		}
	}
}

var db *sql.DB

func main() {
	api, err := ParseRoutes(".")
	if err != nil {
		log.Fatal(err)
	}
	//Plugins
	api.RegisterPlugin("jwt", &jwt.JWT{})
	//Plugins
	db, err = GetDbConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	router := httprouter.New()
	for _, route := range api.Routes {
		if route.Method == "GET" {
			router.GET(route.Path, handler(api, route, 0))
			if api.Version > 0 {
				for i := api.MinVersion; i <= api.Version; i++ {
					router.GET("/v"+strconv.Itoa(i)+route.Path, handler(api, route, i))
				}
			}
		}
		if route.Method == "POST" {
			router.POST(route.Path, handler(api, route, 0))
			if api.Version > 0 {
				for i := api.MinVersion; i <= api.Version; i++ {
					router.POST("/v"+strconv.Itoa(i)+route.Path, handler(api, route, i))
				}
			}
		}
		if route.Method == "PUT" {
			router.PUT(route.Path, handler(api, route, 0))
			if api.Version > 0 {
				for i := api.MinVersion; i <= api.Version; i++ {
					router.PUT("/v"+strconv.Itoa(i)+route.Path, handler(api, route, i))
				}
			}
		}
		if route.Method == "DELETE" {
			router.DELETE(route.Path, handler(api, route, 0))
			if api.Version > 0 {
				for i := api.MinVersion; i <= api.Version; i++ {
					router.DELETE("/v"+strconv.Itoa(i)+route.Path, handler(api, route, i))
				}
			}
		}
	}
	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func goThroughPipelines(api *Api,
	jsonValue string,
	pluginPipelines []*PluginPipeline,
	w http.ResponseWriter) error {

	data := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonValue), &data)
	if err != nil {
		return err
	}
	for _, pp := range pluginPipelines {
		plugin := api.GetPlugin(pp.Name)
		if plugin == nil {
			return errors.New(fmt.Sprintf("Plugin missing: %v", pp.Name))
		}
		response := plugin.Process(data, pp.Argument)
		if response.Headers != nil {
			for name, values := range response.Headers {
				for _, value := range values {
					w.Header().Set(name, value)
				}
			}
		}
		if response.ResponseCode != 0 {
			w.WriteHeader(response.ResponseCode)
			if response.Error != "" {
				fmt.Fprint(w, response.Error)
			}
			return nil
		}
		data = response.Data
	}
	return nil
}
