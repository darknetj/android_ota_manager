package main

import (
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/coreos/go-semver/semver"
)

// Semantic Versioning
// --------------------------------------------------------------------
type SemVerMiddleware struct {
	MinVersion string
	MaxVersion string
}

func (mw *SemVerMiddleware) MiddlewareFunc(handler rest.HandlerFunc) rest.HandlerFunc {
	minVersion, err := semver.NewVersion(mw.MinVersion)
	if err != nil {
		panic(err)
	}

	maxVersion, err := semver.NewVersion(mw.MaxVersion)
	if err != nil {
		panic(err)
	}

	return func(writer rest.ResponseWriter, request *rest.Request) {
		version, err := semver.NewVersion(request.PathParam("version"))
		if err != nil {
			rest.Error(
				writer,
				"Invalid version: "+err.Error(),
				http.StatusBadRequest,
			)
			return
		}

		if version.LessThan(*minVersion) {
			rest.Error(
				writer,
				"Min supported version is "+minVersion.String(),
				http.StatusBadRequest,
			)
			return
		}

		if maxVersion.LessThan(*version) {
			rest.Error(
				writer,
				"Max supported version is "+maxVersion.String(),
				http.StatusBadRequest,
			)
			return
		}

		request.Env["VERSION"] = version
		handler(writer, request)
	}
}

// CORS
// ------------------------------------------------------------------

func corsMiddleware() rest.Middleware {
	return &rest.CorsMiddleware{
		RejectNonCorsRequests: false,
		OriginValidator: func(origin string, request *rest.Request) bool {
				// return origin == "https://builds.copperhead.co"
				return origin == "http://localhost"
		},
		AllowedMethods: []string{"GET", "POST", "PUT"},
		AllowedHeaders: []string{
			"Accept", "Content-Type", "X-Custom-Header", "Origin"},
		AccessControlAllowCredentials: true,
		AccessControlMaxAge:           3600,
	}
}
