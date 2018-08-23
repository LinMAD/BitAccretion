package api

import (
	"github.com/LinMAD/BitAccretion/core/cache"
	"github.com/gorilla/mux"
	"net/http"
)

// API container
type API struct {
	webPath    string
	router     *mux.Router
	controller *controller
	storage    *cache.MemoryCache
}

// NewAPI initialize all registered routes and controllers
func NewAPI(r *mux.Router, st *cache.MemoryCache, wp string) *API {
	a := &API{
		webPath: wp,
		router:  r,
		storage: st,
	}

	a.controller = &controller{
		api: a,
	}

	return a
}

// ServeAllRoutes listen and serve all registered routes
func (api *API) ServeAllRoutes(isServeStatic bool) {
	api.router.HandleFunc("/ws", api.controller.getTrafficDataViaWebSocket).Name("web_socket")
	api.router.HandleFunc("/api/traffic/data", api.controller.getTrafficData).Name("traffic_data")
	if isServeStatic {
		api.router.PathPrefix("/").Handler(http.FileServer(http.Dir(api.webPath + "/resources")))
	}
}
