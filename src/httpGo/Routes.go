package httpGo
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

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/StorageNodeListen", StorageNodeListen)
	router.HandleFunc("/SNodeListenNoP2P", SNodeListenNoP2P)

	router.HandleFunc("/p2pRequest", p2pRequest)
	router.HandleFunc("/ReturnData",ReturnData)
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
		"GetNodes",
		"GET",
		"/GetNodes",
		GetNodes,
	},
	Route{
		"/GetNodesForKey",
		"GET",
		"/GetNodesForKey",
		GetNodesForKey,
	},
	Route{
		"/GetChunks",
		"POST",
		"/GetChunks",
		GetChunks,
	},



}
