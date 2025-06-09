package pkg

import (
	"fmt"
	ht "html/template"
	"net/http"
	"text/template"
)

func ReflectXSS1(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	fmt.Fprintf(w, name)
}

func ReflectXSS2(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(name))
}

func ReflectXSS3(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	t, _ := template.New("ReflectXSS3").Parse("<h1>Hello, I'am {{.Name}}</h1>")
	t.Execute(w, map[string]string{"Name": name})
}

func ReflectXSS4(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	t, _ := template.ParseFiles("./tpl/XSS/ReflectXSS4.gtpl")
	t.Execute(w, map[string]string{"Name": name})
}

func NoXSS1(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	t, _ := ht.New("NoXSS1").Parse("<h1>Hello, I'am {{.Name}}</h1>")
	t.Execute(w, map[string]string{"Name": name})
}
