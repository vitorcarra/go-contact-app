package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/vitorcarra/go-contact-app/api"
	"github.com/vitorcarra/go-contact-app/db"
	"github.com/vitorcarra/go-contact-app/types"
)

func main() {
	var dir string

	flag.StringVar(&dir, "dir", "./ui/static/", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	// fs := http.FileServer(http.Dir("./ui/static/"))
	// http.Handle("/static/", http.StripPrefix("/static/", fs))
	contactStore := db.NewInMemoryContactStore()
	contactStore.CreateContact(types.NewContact(-1, "John", "Doe", "john.doe@hollywood.com", "123456789"))
	contactStore.CreateContact(types.NewContact(-1, "Samuel", "L. Jackson", "samuel.ljackson@hollywood.com", "987123778"))
	contactStore.CreateContact(types.NewContact(-1, "Robert", "De Niro", "robert.niro@hollywood.com", "9871234987"))

	contactsHandler := api.NewContactHandler(contactStore)
	contactRouter := r.PathPrefix("/contacts").Subrouter()
	contactRouter.StrictSlash(true).HandleFunc("/new", contactsHandler.HandleGetNewContact).Methods("GET")
	contactRouter.StrictSlash(true).HandleFunc("/new", contactsHandler.HandlePostNewContact).Methods("POST")
	contactRouter.StrictSlash(true).HandleFunc("/{id}/edit", contactsHandler.HandleGetUpdateContact).Methods("GET")
	contactRouter.StrictSlash(true).HandleFunc("/{id}/edit", contactsHandler.HandleUpdateContact).Methods("POST")
	contactRouter.StrictSlash(true).HandleFunc("/{id}/delete", contactsHandler.HandleDeleteContact).Methods("POST")
	contactRouter.HandleFunc("/", contactsHandler.HandleGetContacts)

	rootHandler := api.NewRootHandler()
	r.HandleFunc("/", rootHandler.HandleGetRoot)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	srv.ListenAndServe()
}
