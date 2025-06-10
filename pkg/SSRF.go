package pkg

import (
	"io"
	"net/http"
)

func SSRF1(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	if url == "" {
		http.Error(w, "url is empty", http.StatusBadRequest)
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	w.Header().Set("Content-Type", "text/html")
	io.Copy(w, resp.Body)
}

func NoSSRF(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path")
	if path == "" {
		path = "go/go-tutorial.html"
	}
	url := "https://www.runoob.com/" + path

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	w.Header().Set("Content-Type", "text/html")
	io.Copy(w, resp.Body)
}
