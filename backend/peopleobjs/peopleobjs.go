package peopleobjs

import "time"

// Structs
// Define a Member struct to hold members in the list
type Member struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Address
	County          string    `json:"county"`
	Phone           string    `json:"phone"`
	Email           string    `json:"email"`
	Last_one_on_one time.Time `json:"last_one_on_one"`
	Issues          []string  `json:"issues"`
	Due_date_pay    time.Time `json:"due_date_pay"`
	Active          bool      `json:"active"`
}

// create the embedded Address struct
type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
	State  string `json:"state"`
	Zip    string `json:"zip"`
}

// create the user struct for logging in
type User struct {
	Username string
	Password string
}
