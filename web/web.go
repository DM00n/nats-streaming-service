package web

import (
	"html/template"
	"log"
	"net/http"
	"wbService/cache"
)

func StartHTTP(cch *cache.Cache) {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		Handler(writer, request, cch)
	})
	fileServer := http.FileServer(http.Dir("./web/static/"))
	http.Handle("/static/", http.StripPrefix("/static", fileServer))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func Handler(w http.ResponseWriter, r *http.Request, cache *cache.Cache) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method == http.MethodGet {
		ts, err := template.ParseFiles("./web/index.page.tmpl")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}
		err = ts.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
		}
	} else if r.Method == http.MethodPost {
		id := r.FormValue("orderID")
		order := cache.Get(id)
		ts, err := template.ParseFiles("./web/index.page.tmpl")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}
		err = ts.Execute(w, order)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
		}
	}
}
