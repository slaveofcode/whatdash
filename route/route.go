package route

import (
	"net/http"
	"whatdash/api"
	"whatdash/wa"
)

type Route struct {
	Name    string
	Method  string
	Path    string
	Handler http.HandlerFunc
}

type Routes []Route

func InitRoutes(wa *wa.ActiveConnections) Routes {
	var DashboardCtrl = &api.Dashboard{WA: wa}
	var WhatsAppCtrl = &api.WhatsApp{ACS: wa}

	return Routes{
		Route{
			Name:    "LIST_CONNECTED_ACCOUNTS",
			Method:  "GET",
			Path:    "/list-connected-accounts",
			Handler: DashboardCtrl.ListConnectedAccounts,
		},
		Route{
			Name:    "WA_LOGIN",
			Method:  "POST",
			Path:    "/wa/login",
			Handler: WhatsAppCtrl.Login,
		},
		Route{
			Name:    "WA_SEND_MSG",
			Method:  "POST",
			Path:    "/wa/send-msg",
			Handler: WhatsAppCtrl.SendMsg,
		},
	}
}
