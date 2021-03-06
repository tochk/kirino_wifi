package memorandums

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/tochk/kirino/auth"
	"github.com/tochk/kirino/templates/html"
)

func FormsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Loaded %s page from %s", r.URL.Path, r.Header.Get("X-Real-IP"))
	vars := mux.Vars(r)
	pageType := vars["type"]
	if pageType == "" {
		pageType = "wifi"
	}
	if auth.IsAdmin(r) {
		pageType = "admin"
	}
	switch vars["type"] {
	case "wifi", "":
		fmt.Fprint(w, html.WifiPage(pageType))
	case "phone":
		fmt.Fprint(w, html.PhonePage(pageType))
	case "mail":
		fmt.Fprint(w, html.MailPage(pageType))
	case "domain":
		fmt.Fprint(w, html.DomainPage(pageType))
	case "ethernet":
		fmt.Fprint(w, html.EthernetPage(pageType))
	case "admin":
		if auth.IsAdmin(r) {
			http.Redirect(w, r, "/wifi/memorandums/view/1", 302)
			return
		}
		fmt.Fprint(w, html.AdminPage())
	}
}
