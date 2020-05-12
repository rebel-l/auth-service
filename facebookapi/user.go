package facebookapi

type User struct {
	ID        string `json:"id"`
	EMail     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Name      string `json:"name"`
}
