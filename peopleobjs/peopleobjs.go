package peopleobjs

import "time"

// Structs
// Define a Member struct to hold members in the list
type Member struct {
	Id              int    //autoincremented unique id for the sql database
	Name            string //First and last
	Address                //embedded struct
	County          string
	Phone           string
	Email           string    //this is a unique id as well
	Last_one_on_one time.Time //this is holding a date in the format "mm/dd/yy"
	Issues          []string  //stored as a comma separated string
	Due_date_pay    time.Time
	Active          bool
}

// create the embedded Address struct
type Address struct {
	Street string
	City   string
	State  string
	Zip    string
}

// create the user struct for logging in
type User struct {
	Username string
	Password string
}
