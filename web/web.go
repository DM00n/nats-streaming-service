package web

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {
	r.FormValue("id")
}
