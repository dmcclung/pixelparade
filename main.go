package main

import (
	"fmt"
	"net/http"
)

func pathHandler(w http.ResponseWriter, r *http.Request) {
  switch r.URL.Path {
  case "/contact": 
    contactHandler(w, r)
  case "/":
    homeHandler(w, r)
  default: 
    notFoundHandler(w, r)
  }
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusNotFound)
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  fmt.Fprint(w, "<h1>Not found</h1><p>Sorry what you requested was not found</p>")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Welcome to my awesome site</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  fmt.Fprint(w, "<h1>Contact Page</h1><p>To contact, <a href=\"mailto:webdevwithgo@gmail.com\">email me</a></p>")
}

func main() {
	http.HandleFunc("/", pathHandler)
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", nil)
}
