package main

import (
	"fmt"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	pageText := `<h1>Hello, World!</h1>
				  <p><a href='/contact'>Contact</a></p>
				  <p><a href='/faq'>FAQ</a></p>`
	fmt.Fprint(w, pageText)
}

func contactFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	pageText := `<h1>Contact Page</h1>
				 <p><a href='/'>Home</a></p>`
	fmt.Fprint(w, pageText)
}

func faqFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	pageText := `<h1>FAQ Page</h1>
	<ul>
		<li><b>Q: Is there a free version?</b></li>
		<p>A: Yes! We offer a free trial for 30 days on any paid plans.</p>
		<li><b>Q: What are your support hours?</b></li>
		<p>A: We have support staff answering emails 24/7, through response times may be a bit slower on weekends.</p>
		<li><b>Q: How do i contact support?</b></li>
		<p>A: Email us - <a href='mailto:support@lenslocked.com'> support@lenslocked.com</p>
	</ul>
	<p><a href='/'>Home</a></p>`
	fmt.Fprint(w, pageText)
}

func pathHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/contact":
		contactFunc(w, r)
	case "/faq":
		faqFunc(w, r)
	default:
		http.Error(w, "Page not found", http.StatusNotFound)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", pathHandler)
	fmt.Println("Server is running on port 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
