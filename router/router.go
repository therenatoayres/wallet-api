package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/therenatoayres/wallet-api/controller"
	"github.com/therenatoayres/wallet-api/logging"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

const (
	basePath         string = "/wallet"
	exchangeRatePath string = "/rate"
)

type listOfRoutes []route

//Router creates and returns the list of routes for the server
func Router(router *mux.Router) *mux.Router {
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = logging.HttpLogger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

var routes = listOfRoutes{

	// rate: used to get the currency rate between one currency to one or more diffent currencies
	route{"GetCurrencyRate", "GET", basePath + exchangeRatePath, controller.GetCurrencyRate},
}
