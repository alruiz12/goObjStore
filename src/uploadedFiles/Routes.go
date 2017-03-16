package conf
import (
	"net/http"
	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}
type Routes []Route

/*
Router using gorilla/mux
*/
func MyNewRouter() *mux.Router {

	/*
	go func() {
		http.ListenAndServe(":8080", nil)
	}()*/

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/upLoadFile", upLoadFile)
	for _, route := range routes {
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
		"addPeer",
		"POST",
		"/addPeer",
		addPeer,
	},
	Route{
		"getTorrentsList",
		"GET",
		"/getTorrentsList",
		getTorrentsList,
	},
	Route{
		"getIPs",
		"POST",
		"/getIPs",
		getIPs,
	},
	/*Route{
		"upLoadFile",
		"POST",
		"/upLoadFile",
		upLoadFile,
	},*/

}