package main

import "github.com/ant0ine/go-json-rest/rest"

type Route struct {
	Path    string
	Method  string
	Handler rest.HandlerFunc
}

var routes = []Route{
	// Releases
	Route{"/", "GET", Releases},
	Route{"/releases", "GET", Releases},
	Route{"/releases/:id", "GET", Release},
	Route{"/releases", "POST", ReleaseCreate},
	Route{"/releases/:id", "PUT", ReleaseUpdate},
	Route{"/releases/:id", "DELETE", ReleaseDelete},

	// Files
	Route{"/files", "GET", Files},
	Route{"/files/:id", "GET", File},
	Route{"/files", "POST", FileCreate},

	// Users
	Route{"/users", "GET", Users},

	// Authentication
	Route{"/login", "POST", Login},
	Route{"/logout", "DESTROY", Logout},

	// CMUpdater
	Route{"/", "POST", Releases},
	Route{"/v1/build/get_delta", "GET", Releases},
	Route{"/changelog/:incremental.txt", "GET", ReleaseChangelog},
}

// Generates go-json-rest routes from our simple Route
// Example: rest.Get("/countries", GetAllCountries),
func generateRoute(route Route, versioned bool) *rest.Route {
	fn := rest.Get
	if route.Method == "POST" {
		fn = rest.Post
	}
	if route.Method == "PUT" {
		fn = rest.Put
	}
	if route.Method == "DELETE" {
		fn = rest.Delete
	}
	if versioned {
		// Convert route to a semver-style versioned route
		// For ex: /releases -> /api/1.0.0/releases
		return fn(cat("/#version", route.Path), ApiVersions.MiddlewareFunc(route.Handler))
	} else {
		return fn(route.Path, route.Handler)
	}
}

func generateRoutes() []*rest.Route {
	// Generate both versioned and non-versioned APIs
	restRoutes := make([]*rest.Route, 0)
	for _, route := range routes {
		rt := generateRoute(route, false)
		restRoutes = append(restRoutes, rt)
		rt = generateRoute(route, true)
		restRoutes = append(restRoutes, rt)
	}
	return restRoutes
}
