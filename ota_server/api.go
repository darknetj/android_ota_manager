package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/coreos/go-semver/semver"
)

// Api spec
// -------------------------------------------------------------------

var ApiCurrentVersion = "2.0.0"

var ApiVersions = SemVerMiddleware{
	MinVersion: "1.0.0",
	MaxVersion: ApiCurrentVersion,
}

// Version 1.0.0 API response format
type ApiResponseLegacy struct {
	Id     string      `json:"id"`
	Result interface{} `json:"result"`
	Error  interface{} `json:"error"`
}

// Version 2.0.0+ API response format
type ApiResponse struct {
	Data    ApiResource `json:"data"`
	Errors  []ApiError  `json:"errors"`
	JsonApi JsonApi     `json:"jsonapi"`
	Links   ParentLink  `json:"links"`
}

type ApiListResponse struct {
	Data    []ApiResource `json:"data"`
	Errors  []ApiError    `json:"errors"`
	JsonApi JsonApi       `json:"jsonapi"`
	Links   Links         `json:"links"`
}

type JsonApi struct {
	Version string `json:"version"`
}

type Links struct {
	Self string `json:"self"`
}

type ParentLink struct {
	Parent interface{} `json:"parent"`
}

type ApiResource struct {
	Type       string      `json:"type"`
	Id         string      `json:"id"`
	Attributes interface{} `json:"attributes"`
	Links      Links       `json:"links"`
}

type ApiError struct {
	Id     string `json:"name"`
	Status string `json:"error"`
	Code   string `json:"ip_address"`
	Title  string `json:"title"`
	Detail string `json:"title"`
	Source string `json:"source"`
}

// Api Server
// -------------------------------------------------------------------

func server(port string, staticDir string) {
	log.Println("--- Started OTA Server on port", port, "---")

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	api.Use(corsMiddleware())

	restRoutes := generateRoutes()
	log.Println(restRoutes)
	router, err := rest.MakeRouter(restRoutes...)

	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))
	log.Fatal(http.ListenAndServe(cat(":", port), nil))
}

func exampleVersioned(w rest.ResponseWriter, req *rest.Request) {
	version := req.Env["VERSION"].(*semver.Version)
	if version.Major == 2 {
		// http://en.wikipedia.org/wiki/Second-system_effect
		w.WriteJson(map[string]string{
			"Body": "Hello broken World!",
		})
	} else {
		w.WriteJson(map[string]string{
			"Body": "Hello World!",
		})
	}
}

// Api response generators
// -------------------------------------------------------------------

func apiUrl(req *rest.Request) string {
	return cat(host, "/api", req.URL.RequestURI())
}

func apiParentUrl(req *rest.Request) string {
	var path []string
	a := strings.Split(req.URL.RequestURI(), "/")
	_, path = a[len(a)-1], a[:len(a)-1]
	parent := strings.Join(path,"/")
	return cat(host, "/api", parent)
}

func apiResponse(req *rest.Request, data ApiResource) ApiResponse {
	return ApiResponse{
		Data:   data,
		Errors: nil,
		JsonApi: JsonApi{
			Version: ApiCurrentVersion,
		},
		Links: ParentLink{
			Parent: apiParentUrl(req),
		},
	}
}

func apiListResponse(req *rest.Request, data []ApiResource) ApiListResponse {
	return ApiListResponse{
		Data:   data,
		Errors: nil,
		JsonApi: JsonApi{
			Version: ApiCurrentVersion,
		},
		Links: Links{
			Self: apiUrl(req),
		},
	}
}

func Resource(id int64, resourceType string, path string) ApiResource {
	resourceId := fmt.Sprintf("%v", id)
	return ApiResource{
		Type: resourceType,
		Id:   resourceId,
		Links: Links{
			Self: cat(host, "/api/", ApiCurrentVersion, path, "/", resourceId),
		},
	}
}
