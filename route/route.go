package route

import (
	"net/http"
	"whatdash/api"
)

type Route struct {
	Name    string
	Method  string
	Path    string
	Handler http.HandlerFunc
}

type Routes []Route

var DashboardCtrl = &api.Dashboard{}

var ApiRoutes = Routes{
	Route{
		Name:    "LIST_CONNECTED_ACCOUNTS",
		Method:  "GET",
		Path:    "/list-connected-accounts",
		Handler: DashboardCtrl.ListConnectedAccounts,
	},
}
