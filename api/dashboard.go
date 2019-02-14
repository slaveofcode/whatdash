package api

import (
	"net/http"
	"whatdash/wa"
)

type Dashboard struct {
	Storage *wa.Storage
}

func (c *Dashboard) NewAccount(w http.ResponseWriter, r *http.Request) {

}

func (c *Dashboard) ListConnectedAccounts(w http.ResponseWriter, r *http.Request) {
	ResponseJSON(w, 200, []byte(`{"list": [1,2,3]}`))
}

func (c *Dashboard) ReconnectAccount(w http.ResponseWriter, r *http.Request) {

}
