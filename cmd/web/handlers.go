package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to quotebox"))
}

func (app *application) createQuoteForm(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/html/quotes_form_page.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

}

func (app *application) createQuote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/quote", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}

	author := r.PostForm.Get("author_name")
	category := r.PostForm.Get("category")
	quote := r.PostForm.Get("quote")
	// check the web form fields to validity
	errors := make(map[string]string)
	// check each field
	if strings.TrimSpace(author) == "" {
		errors["author"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(author) > 50 {
		errors["author"] = "This field is too long(maximum is 50)"
	}
	if strings.TrimSpace(category) == "" {
		errors["category"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(category) > 25 {
		errors["category"] = "This field is too long(maximum is 50)"
	}
	if strings.TrimSpace(quote) == "" {
		errors["quote"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(quote) > 255 {
		errors["quote"] = "This field is too long(maximum is 50)"
	}
	// Check if there are any errors in the map
	if len(errors) > 0 {
		fmt.Fprint(w, errors)
		return
	}

	//Insert a quote
	id, err := app.quotes.Insert(author, category, quote)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "row with id: %v has been inserted.", id)
}

func (app *application) displayQuotation(w http.ResponseWriter, r *http.Request) {

	q, err := app.quotes.Read()
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	// Dsplay the quotes using a template
	ts, err := template.ParseFiles("./ui/html/show_page.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, q)

	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

}
