package api

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vitorcarra/go-contact-app/db"
	"github.com/vitorcarra/go-contact-app/types"
)

type ContactHandler struct {
	contactStore db.ContactStore
}

type Params struct {
	Q        string
	Contacts []*types.Contact
}

func (h *ContactHandler) HandleGetContacts(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("q")
	contacts, err := h.contactStore.GetContacts()
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	output := []*types.Contact{}
	if search != "" {
		log.Printf("Search query: %q", search)
		for _, item := range contacts {
			if item.FirstName == search || item.LastName == search || item.Email == search || item.Phone == search {
				output = append(output, item)
				break
			}
		}
	} else {
		output = contacts
	}

	template, err := template.ParseFiles(
		"ui/html/base.tmpl.html",
		"ui/components/base/navbar.tmpl.html",
		"ui/html/contacts.tmpl.html",
	)

	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	p := Params{Q: search, Contacts: output}
	err = template.Execute(w, p)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *ContactHandler) HandleGetUpdateContact(w http.ResponseWriter, r *http.Request) {
	muxVars := mux.Vars(r)
	muxVarsIdStr := muxVars["id"]
	muxVarsId, err := strconv.ParseInt(muxVarsIdStr, 10, 64)
	if err != nil {
		log.Printf("Invalid ID: %v", muxVarsIdStr)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	contact, err := h.contactStore.GetContact(muxVarsId)

	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	template, err := template.ParseFiles(
		"ui/html/base.tmpl.html",
		"ui/components/base/navbar.tmpl.html",
		"ui/html/contact_edit.tmpl.html",
	)

	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = template.Execute(w, contact)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *ContactHandler) HandleGetNewContact(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles(
		"ui/html/base.tmpl.html",
		"ui/components/base/navbar.tmpl.html",
		"ui/html/contact_new.tmpl.html",
	)

	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = template.Execute(w, nil)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *ContactHandler) HandlePostNewContact(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	c := types.NewContact(
		-1,
		r.FormValue("first_name"),
		r.FormValue("last_name"),
		r.FormValue("email"),
		r.FormValue("phone"),
	)

	err = h.contactStore.CreateContact(c)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("New contact: %v", c)
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

func (h *ContactHandler) HandleUpdateContact(w http.ResponseWriter, r *http.Request) {
	muxVars := mux.Vars(r)
	muxVarsIdStr := muxVars["id"]
	muxVarsId, err := strconv.ParseInt(muxVarsIdStr, 10, 64)
	if err != nil {
		log.Printf("Invalid ID: %v", muxVarsIdStr)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	currentContact, err := h.contactStore.GetContact(muxVarsId)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	err = r.ParseForm()
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	c := types.NewContact(
		currentContact.ID,
		r.FormValue("first_name"),
		r.FormValue("last_name"),
		r.FormValue("email"),
		r.FormValue("phone"),
	)

	err = h.contactStore.UpdateContact(c)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("Updated contact: %v", c)
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

func (h *ContactHandler) HandleDeleteContact(w http.ResponseWriter, r *http.Request) {
	muxVars := mux.Vars(r)
	muxVarsIdStr := muxVars["id"]
	id, err := strconv.ParseInt(muxVarsIdStr, 10, 64)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = h.contactStore.DeleteContact(&types.Contact{ID: id})
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

func NewContactHandler(c db.ContactStore) *ContactHandler {
	return &ContactHandler{
		contactStore: c,
	}
}
