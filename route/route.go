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

func InitRoutes(s *wa.Storage) Routes {
	var DashboardCtrl = &api.Dashboard{Storage: s}
	var WhatsAppCtrl = &api.WhatsApp{Storage: s}

	return Routes{
		Route{
			Name:    "LIST_CONNECTED_ACCOUNTS",
			Method:  "GET",
			Path:    "/list-connected-accounts",
			Handler: DashboardCtrl.ListConnectedAccounts,
		},
		Route{
			Name:    "WA_CREATE_SESSION",
			Method:  "POST",
			Path:    "/wa/session/create",
			Handler: WhatsAppCtrl.CreateSession,
		},
		Route{
			Name:    "WA_CHECK_REGISTER",
			Method:  "POST",
			Path:    "/wa/session/check",
			Handler: WhatsAppCtrl.CheckSession,
		},
		Route{
			Name:    "WA_SEND_TEXT",
			Method:  "POST",
			Path:    "/wa/send/text",
			Handler: WhatsAppCtrl.SendText,
		},
	}
}
