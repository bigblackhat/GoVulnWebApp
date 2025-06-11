package pkg

import (
	"fmt"
	"net/http"
)

/*
http.Redirect函数
http://127.0.0.1:8080/Redirect/Redirect1?url=/SSRF/SSRF1
*/
func Redirect1(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	http.Redirect(w, r, url, 301)
}

/*
手动设置Header
*/
func Redirect2(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusMovedPermanently)
}

/*
HTML元刷新（客户端跳转）
*/
func Redirect3(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<meta http-equiv="refresh" content="0; URL='`+url+`'" />`)
}
