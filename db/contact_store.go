package db

import (
	"errors"

	"github.com/vitorcarra/go-contact-app/types"
)

type ContactStore interface {
	GetContacts() ([]*types.Contact, error)
	GetContact(id int64) (*types.Contact, error)
	CreateContact(contact *types.Contact) error
	DeleteContact(contact *types.Contact) error
	UpdateContact(contact *types.Contact) error
}

type InMemoryContactStore struct {
	contacts []*types.Contact
	id       int64
}

func NewInMemoryContactStore() *InMemoryContactStore {
	return &InMemoryContactStore{
		contacts: []*types.Contact{},
		id:       0,
	}
}

func (s *InMemoryContactStore) GetContacts() ([]*types.Contact, error) {
	return s.contacts, nil
}

func (s *InMemoryContactStore) GetContact(id int64) (*types.Contact, error) {
	for _, contact := range s.contacts {
		if contact.ID == id {
			return contact, nil
		}
	}
	return nil, errors.New("contact not found")
}

func (s *InMemoryContactStore) CreateContact(contact *types.Contact) error {
	s.id += 1
	contact.ID = s.id
	s.contacts = append(s.contacts, contact)
	return nil
}

func (s *InMemoryContactStore) UpdateContact(contact *types.Contact) error {
	for i, item := range s.contacts {
		if item.ID == contact.ID {
			s.contacts[i] = contact
			return nil
		}
	}
	return errors.New("contact not found")
}

func (s *InMemoryContactStore) DeleteContact(contact *types.Contact) error {
	for i, item := range s.contacts {
		if item.ID == contact.ID {
			s.contacts = append(s.contacts[:i], s.contacts[i+1:]...)
			return nil
		}
	}
	return errors.New("contact not found")
}
