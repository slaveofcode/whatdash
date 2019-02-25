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

func InitRoutes(s *wa.BucketSession) Routes {
	var waSessHandler = api.SessionHandler{Bucket: s}

	var DashboardCtrl = &api.Dashboard{SessionHandler: waSessHandler}
	var WhatsAppCtrl = &api.WhatsApp{SessionHandler: waSessHandler}

	return Routes{
		Route{
			Name:    "LIST_CONNECTED_ACCOUNTS",
			Method:  "GET",
			Path:    "/account/list-connected",
			Handler: DashboardCtrl.ListConnectedAccounts,
		},
		Route{
			Name:    "DETAIL_ACCOUNT",
			Method:  "GET",
			Path:    "/account/detail/{id}",
			Handler: DashboardCtrl.DetailAccount,
		},
		Route{
			Name:    "LOAD_CHAT_HISTORY",
			Method:  "POST",
			Path:    "/chat/history",
			Handler: DashboardCtrl.LoadChatHistory,
		},
		Route{
			Name:    "LOAD_CHATS",
			Method:  "POST",
			Path:    "/chat/list",
			Handler: DashboardCtrl.LoadChats,
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
			Name:    "WA_DESTROY_SESSION",
			Method:  "POST",
			Path:    "/wa/session/destroy",
			Handler: WhatsAppCtrl.Destroy,
		},
		Route{
			Name:    "WA_LOAD_CONTACTS",
			Method:  "POST",
			Path:    "/wa/contact/load",
			Handler: WhatsAppCtrl.LoadContacts,
		},
		Route{
			Name:    "WA_GET_CONTACTS",
			Method:  "POST",
			Path:    "/wa/contact/list",
			Handler: WhatsAppCtrl.GetContacts,
		},
		Route{
			Name:    "WA_SEND_TEXT",
			Method:  "POST",
			Path:    "/wa/send/text",
			Handler: WhatsAppCtrl.SendText,
		},
		Route{
			Name:    "WA_LOAD_MESSAGES",
			Method:  "POST",
			Path:    "/wa/messages/load",
			Handler: WhatsAppCtrl.TriggerLoadMessage,
		},
		Route{
			Name:    "WA_LOAD_NEW_MESSAGES",
			Method:  "POST",
			Path:    "/wa/messages/load-next",
			Handler: WhatsAppCtrl.TriggerLoadNewMessage,
		},
		Route{
			Name:    "WA_LOAD_OLD_MESSAGES",
			Method:  "POST",
			Path:    "/wa/messages/load-prev",
			Handler: WhatsAppCtrl.TriggerLoadOldMessage,
		},
		Route{
			Name:    "WA_EXP",
			Method:  "POST",
			Path:    "/wa/messages/longpool",
			Handler: WhatsAppCtrl.LongPoolExp,
		},
	}
}
