package conf
import (
	"net/http"
	"simpleBT/src/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func MyNewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

var routes = Routes{
	Route{
		"addTorrent",
		"POST",
		"/addTorrent",
		addTorrent,
	},
	Route{
		"showTorrents",
		"GET",
		"/addTorrent",
		showTorrents,
	},
	Route{
		"getIPs",
		"GET",
		"/getIPs",
		getIPs,
	},
}