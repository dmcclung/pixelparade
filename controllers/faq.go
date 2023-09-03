package controllers

import (
	"html/template"
	"net/http"
)

func Faq(tmplt Template) http.HandlerFunc {
	faqItems := []struct {
		Question string
		Answer   template.HTML
	}{
		{
			Question: "Why did you develop this web app?",
			Answer:   "We developed this web app to provide a convenient and efficient solution for managing an image gallery. Our aim is to make image galleries easier and more accessible for everyone.",
		},
		{
			Question: "Will it always be free?",
			Answer:   "As of now, the basic features of this web app are free to use. However, we may introduce premium features in the future that could be available for a fee. We are committed to always offering a free version with essential functionalities.",
		},
		{
			Question: "Who do I contact for help?",
			Answer:   "If you encounter any issues or have questions, you can reach out to our support team at <a href=\"mailto:pixelparade@gmail.com\">pixelparade@gmail.com</a>. We're here to help!",
		},
		{
			Question: "Will you translate the app?",
			Answer:   "Yes! Translations of the app will be up soon",
		},
		{
			Question: "Where is your team located?",
			Answer:   "Team is remote, but based in the US and Canada",
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		err := tmplt.Execute(w, faqItems)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
