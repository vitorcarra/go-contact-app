package types

type Contact struct {
	ID        int64  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

func NewContact(id int64, firstName, lastName, email, phone string) *Contact {
	return &Contact{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
	}
}
