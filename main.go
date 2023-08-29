package main

import (
	"fmt"
	"net/http"
)

func pathHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/faq":
		faqHandler(w, r)
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

func faqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `
  <h1>Frequently Asked Questions</h1>

  <div style="margin-bottom: 20px;">
      <h2>1. Why did you develop this web app?</h2>
      <p>We developed this web app to provide a convenient and efficient solution for managing an image gallery. Our aim is to make image galleries easier and more accessible for everyone.</p>
  </div>

  <div style="margin-bottom: 20px;">
      <h2>2. Will it always be free?</h2>
      <p>As of now, the basic features of this web app are free to use. However, we may introduce premium features in the future that could be available for a fee. We are committed to always offering a free version with essential functionalities.</p>
  </div>

  <div style="margin-bottom: 20px;">
      <h2>3. Who do I contact for help?</h2>
      <p>If you encounter any issues or have questions, you can reach out to our support team at <a href="mailto:webdevwithgo@gmail.com">webdevwithgo@gmail.com</a>. We're here to help!</p>
  </div>`)
}

func main() {
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", http.HandlerFunc(pathHandler))
}
