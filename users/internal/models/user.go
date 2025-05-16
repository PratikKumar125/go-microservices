package models

type User struct {
	Id		string	`json:"id"`
	Name	string	`json:"name"`
	Email	string	`json:"email"`
}

type UserSearchFilters struct {
	Email string
	Name string
}