package controllers

import (
	"html/template"
	"net/http"
)

func StaticHandler(tpl Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, nil)
	}
}

func FAQ(tpl Template) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   template.HTML
	}{
		{
			Question: "What is LensLocked?",
			Answer:   "LensLocked is a web application for photographers to showcase and manage their portfolios.",
		},
		{
			Question: "How do I create an account?",
			Answer:   "You can create an account by clicking on the 'Sign Up' button and filling out the registration form.",
		},
		{
			Question: "Is LensLocked free to use?",
			Answer:   "Yes, LensLocked offers a free basic plan. We also have premium plans with additional features.",
		},
		{
			Question: "Can I organize my photos into albums?",
			Answer:   "Absolutely! LensLocked allows you to create and manage multiple albums to organize your photos.",
		},
		{
			Question: "How can I contact support?",
			Answer:   `You can reach our support team by emailing <a href="mailto:support@lenslocked.com">support@lenslocked.com</a> or through the <a href="/contact">Contact Us</a> form on our website.`,
		},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, questions)
	}
}
