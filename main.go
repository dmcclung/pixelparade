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
    http.Error(w, "Page not found", http.StatusNotFound)
  }
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
