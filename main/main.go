package main

//imports
import (
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

//Main backend logic for Coalition App

// I could encapsulate these structs into a separate package and write some statements for them like getters
// and setters if I need to, but I think working with Gin it might be redundant. If I am just translating json
// strings and interacting with the database, should I even bother with encapsulation?
// Define a Member struct to hold members in the list
type Member struct {
	id          int    //autoincremented unique id for the sql database
	name        string //First and last
	Address            //embedded struct
	county      string
	phone       int
	email       string    //this is a unique id as well
	last1_1Date time.Time //this is holding a date in the format "mm/dd/yy"
	issues      []string  //stored as a comma separated string
	dueDatePay  time.Time
	active      bool
}

// create the embedded Address struct
type Address struct {
	street string
	city   string
	state  string
	zip    int
}

// create the user struct for logging in
type User struct {
	Username string
	Password string
}

//create the sql table to hold the member structs
/*
CREATE TABLE IF NOT EXISTS members (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	street TEXT,
	city TEXT,
	state TEXT,
	zip INTEGER,
	phone INTEGER,
	email TEXT UNIQUE
	last_one_on_one TEXT,
	issues TEXT,
	due_date_pay TEXT,
	active BOOLEAN
);

//create a users table in sql to hold the users that can log in and access the data
//cyber security specialist everywhere are losing their MINDS right now.
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_name TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL
);
*/

//Main features

//login authentication. The user name and password will be hardcoded for testing
//or displayed on the page for testing th feature. User: admin; pw: admin

func userLogin(userName string, password string) bool {
	//check the username and password against a table of usernames and passwords
	// SELECT * FROM users WHERE user_name = userName AND password = password
	//return true if user/pw is correct/found, false if not

}

// Gin handler for making member struct from json
func handleLogin(c *gin.Context) {
	//Bind json to user struct
	//call userLogin()
	//if found/correct send json: {"message" : "Login Successful"}
	//else send json: {"error" : "Invalid credentials"}  401 status
}

// let the user upload a new .csv file and use it to populate the database.
// not sure about this parameter yet. Also -- should it return a bool/error?
func uploadCSVList(filereader io.Reader) error {
	//get a .csv file from gin as a json (probably).
	//Parse the json and call rowToStruct() on each row
	//call addNewMember() to add the structs to the db
}

// take a row from the json and map it to a member struct
func rowToStruct(row []string) (Member, error) {
	//map the string array to the fields in the member struct
}

// return a list of all of the members (after updating, adding or deleting). Hand off to
// gin to display on the user interface. May need a parameter here, not sure yet.
func listMembers() []string {
	//SELECT * FROM members
	//call rowToStruct() to send to gin to be put in a json.
	//do we really need to create a member struct here, since the db has the form the json needs in
	//the first place?
}

// return a list of all of the members that need to have one on one conversations.
// This means any inactive or member whose last convo was more than six months ago
// not sure about parameters/returns yet. How does this work with gin, are we building the json here?
func listMembersNeedOneOnOnes() ([]string, error) {
	//today := time.Now()
	//threshold := today minus six months
	//SELECT * FROM members WHERE last_one_on_one < threshold
	//convert to member structs? rowToStruct()
}

// add a new member to the database. Here do we want to parse the json with the member request and put it
// in a member struct before we pass it here, or should we do it all in this method?
// not sure if I want to use a string[] or a member struct as the argument.
func addNewMember(member []string) error {
	//add the new struct to the members table in the database
}

// delete a member from the sql member table. Given a member struct to delete. How are we identifying
// the member to delete? We could use the UI to pass the entire member struct via the id of the struct
// identified by UI elements like a mouse click. Could also return a bool?
func deleteMember(m Member) error {
	//use a unique identifier like the email or the id to locate the memeber struct
	//delete the member struct from the database
	//DELETE FROM members WHERE id = m.id
	//DELETE FROM members WHERE email = m.email
}

// updates a member feild from the UI selections. Should we update just the changes or swap out the whole
// member struct no matter how much it has changed? For performance I think the former. This could also
// return a bool
func updateMember(m Member) error {
	//may just get the whole member struct we need to change and then have to look for the
	//fields that actually changed in this method.
	//update the fields that changed in the database.
	//UPDATE members SET ... WHERE id = m.id
}

// this needs to return a list of member structs where each member has the keyword given in the issues
// column of the database (comma separated string -> issues). User needs to choose from a drop down to avoid
// complexity for now. This can be expanded later.
func searchVolunteers(keyword string) ([]Member, error) {
	//user gives a keyword (using a drop down on the UI)
	//search the database for a list of members that have that keyword in
	//SELECT * FROM members WHERE issues LIKE '%keyword%'
}

// returns a list of members that has an active status set to false.
func listInactive() ([]Member, error) {
	//search the database for a list of members that have active = false
	//SELECT * FROM members WHERE active = false
}
