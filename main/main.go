package main

//imports
import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

//Main backend logic for Coalition App

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

//Database

// create the sql table to hold the member structs
func initializeDatabase() error {
	db, err := sql.Open("sqlite3", "./coalition.db")
	if err != nil {
		log.Println("Could not open database: ", err)
		return err
	}
	//defer means that the db will close when the function is finished executing - even though it's declared here
	defer db.Close()
	createMembers := `CREATE TABLE IF NOT EXISTS members (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		street TEXT,
		city TEXT,
		state TEXT,
		zip TEXT,
		county TEXT,
		phone TEXT,
		email TEXT UNIQUE,
		last_one_on_one TEXT,
		issues TEXT,
		due_date_pay TEXT,
		active BOOLEAN
		);`
	//the db.Exec function will handle sql statements that insert or update the database.
	_, err = db.Exec(createMembers)
	if err != nil {
		return err
	}

	//create a users table in sql to hold the users that can log in and access the data
	//cyber security specialist everywhere are losing their MINDS right now.
	createUsers := `CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_name TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL
	);`

	_, err = db.Exec(createUsers)
	if err != nil {
		return err
	}
	return nil
}

//Main features ---- Go Methods

//login authentication. The user name and password will be hardcoded for testing
//or displayed on the page for testing th feature. User: admin; pw: admin

func userLogin(userName string, password string) bool {
	//check the username and password are not empty
	if userName == "" || password == "" {
		log.Println("Username or password cannot be empty.")
	}
	//open the database
	db, err := sql.Open("sqlite3", "./coalition.db")
	if err != nil {
		log.Println("Could not open database: %w", err)
	}

	defer db.Close()

	var exists bool
	//use the database package to query the sqlite3 database
	err = db.QueryRow(`SELECT EXISTS (SELECT 1 FROM users WHERE username = ? AND password = ?)`, userName, password).Scan(&exists)
	if err != nil {
		log.Println("Could not query database: %w", err)
		return false
	}
	//return true if user/pw is correct/found, false if not
	return exists
}

// let the user upload a new .csv file and use it to populate the database.
// not sure about this parameter yet. Also -- should it return a bool/error?
func uploadCSVList(filereader io.Reader) error {
	//open up a reader using the encoding/csv package
	reader := csv.NewReader(filereader)
	//skip the header row
	_, _ = reader.Read()

	//lets track how many rows fail, if any
	var failures int
	//loop through the rows and call rowToStruct() to convert each row to a member struct
	for {
		//read each row
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Could not read row: %w", err)
			failures++
			continue
		}
		//try to convert the row to a struct
		member, err := rowToStruct(row)
		if err != nil {
			log.Println("Could not convert row to struct: %w", err)
			failures++
			continue
		}
		//add the new struct to the database
		err = addNewMember(member)
		if err != nil {
			log.Println("Could not add new member")
			failures++
			continue
		}
	}
	if failures > 0 {
		return fmt.Errorf("csv processed with %d failures", failures)
	}
	return nil
}

// take a row from the json and map it to a member struct
func rowToStruct(row []string) (Member, error) {
	//rows are just a slice of strings from a []string
	//create a new member struct and error to return.
	var m Member
	var err error
	//access each row in the []string slice and map it to the correct field in the member struct
	m.Name = row[0]
	m.Address = Address{
		Street: row[1],
		City:   row[2],
		State:  row[3],
		Zip:    row[4],
	}
	m.County = row[5]
	m.Phone = row[6]
	m.Email = row[7]

	//formating the date correctly
	m.Last_one_on_one, err = time.Parse("01/02/2006", row[8])
	if err != nil {
		return m, err
	}
	//issues were saved as a comma separated string, so here we are splitting the string at each comma
	//to fill in the issues column.
	m.Issues = strings.Split(row[9], ",")
	m.Due_date_pay, err = time.Parse("01/02/2006", row[10])
	if err != nil {
		return m, err
	}

	m.Active = strings.ToLower(row[11]) == "true" || row[11] == "1"
	return m, nil
}

// add a new member to the database. Here do we want to parse the json with the member request and put it
// in a member struct before we pass it here, or should we do it all in this method?
// not sure if I want to use a string[] or a member struct as the argument.
func addNewMember(m Member) error {
	db, err := sql.Open("sqlite3", "./coalition.db")
	if err != nil {
		return err
	}
	//add the new struct to the members table in the database
	query := `INSERT INTO members (
				name, street, city, state, zip, county, 	
				phone, email, last_one_on_one, issues,
				due_date_pay, active
	) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	//add the variables that hold the actual data into the db.Exec function, they will replace the '?'
	_, err = db.Exec(query,
		m.Name,
		m.Street,
		m.City,
		m.State,
		m.Zip,
		m.County,
		m.Phone,
		m.Email,
		m.Last_one_on_one.Format("01/02/2006"),
		strings.Join(m.Issues, ","),
		m.Due_date_pay.Format("01/02/2006"),
		m.Active,
	)
	if err != nil {
		return err
	}
	return nil
}

// return a list of all of the members (after updating, adding or deleting). Hand off to
// gin to display on the user interface.
func listMembers() ([]Member, error) {
	//create an array of Members to give back to gin
	var members []Member
	db, err := sql.Open("sqlite3", "./coalition.db")
	if err != nil {
		log.Println(err)
		return members, err
	}
	defer db.Close()
	//query the db to get the rows
	rows, err := db.Query(`SELECT * FROM members`)
	if err != nil {
		log.Println(err)
		return members, err
	}
	defer rows.Close()
	//loop through each row and add it to a new member's fields.
	for rows.Next() {
		var m Member
		var lastOneOnOne string
		var issuesCSV string
		var dueDatePay string

		err := rows.Scan(
			&m.Id,
			&m.Name,
			&m.Street,
			&m.City,
			&m.State,
			&m.Zip,
			&m.County,
			&m.Phone,
			&m.Email,
			&lastOneOnOne,
			&issuesCSV,
			&dueDatePay,
			&m.Active,
		)
		if err != nil {
			log.Println(err)
			continue
		}
		//some of the fields need to be formatted
		m.Last_one_on_one, _ = time.Parse("01/02/2006", lastOneOnOne)
		m.Issues = strings.Split(issuesCSV, ",")
		m.Due_date_pay, _ = time.Parse("01/02/2006", dueDatePay)
		//add the new member to the list of members
		members = append(members, m)
	}
	return members, nil
}

// return a list of all of the members that need to have one on one conversations.
// This means any inactive or member whose last convo was more than six months ago
// not sure about parameters/returns yet. How does this work with gin, are we building the json here?
func listMembersNeedOneOnOnes() error {
	//today := time.Now()
	//threshold := today minus six months
	//SELECT * FROM members WHERE last_one_on_one < threshold
	//convert to member structs? rowToStruct()
	return nil
}

// delete a member from the sql member table. Given a member struct to delete. How are we identifying
// the member to delete? We could use the UI to pass the entire member struct via the id of the struct
// identified by UI elements like a mouse click. Could also return a bool?
func deleteMember(m Member) error {
	//use a unique identifier like the email or the id to locate the memeber struct
	//delete the member struct from the database
	//DELETE FROM members WHERE id = m.id
	//DELETE FROM members WHERE email = m.email
	return nil

}

// updates a member feild from the UI selections. Should we update just the changes or swap out the whole
// member struct no matter how much it has changed? For performance I think the former. This could also
// return a bool
func updateMember(m Member) error {
	//may just get the whole member struct we need to change and then have to look for the
	//fields that actually changed in this method.
	//update the fields that changed in the database.
	//UPDATE members SET ... WHERE id = m.id
	return nil
}

// this needs to return a list of member structs where each member has the keyword given in the issues
// column of the database (comma separated string -> issues). User needs to choose from a drop down to avoid
// complexity for now. This can be expanded later.
func searchVolunteers(keyword string) error {
	//user gives a keyword (using a drop down on the UI)
	//search the database for a list of members that have that keyword in
	//SELECT * FROM members WHERE issues LIKE '%keyword%'
	return nil
}

// returns a list of members that has an active status set to false.
func listInactive() error {
	//search the database for a list of members that have active = false
	//SELECT * FROM members WHERE active = false
	return nil
}

//GIN Handlers

// connect gin to the listMembers method by adding a handler
func handleListMembers(c *gin.Context) {
	members, err := listMembers()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve members"})
		return
	}
	c.IndentedJSON(http.StatusOK, members)
}

// Gin handler for making member struct from json. gin.H is an alias for a map of strings used to make jsons.
func handleLogin(c *gin.Context) {
	//create a new user struct
	var u User
	//try to use gin to map the json fields to the member structs using the BindJSON method
	if err := c.BindJSON(&u); err != nil {
		//if binding fails, send an error json
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
	}
	//call the go method userLogin to check the database for the un/pw
	if userLogin(u.Username, u.Password) {
		//if found and correct, send a success message
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Login Successful"})
	} else {
		//if found and incorrect or not found send a failure message
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Invalid Credentials"})
	}
}

// Gin method to handle the upload and send messages back to the frontend
func handleUpload(c *gin.Context) {
	//c.FormFile is a method that gets the file metadata from the frontend react FormData fetch request.
	file, err := c.FormFile("file")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	//file.Open creates a reader when it opens the file
	f, err := file.Open()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Could not open file"})
		return
	}
	//close the file reader when the function exits
	defer f.Close()
	//upload the csv so we can take the rows and map them into member structs. f is already a reader
	err = uploadCSVList(f)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "CSV parsing failed"})
		return
	}
	//send a success message to the front end.
	c.IndentedJSON(http.StatusOK, gin.H{"message": "CSV upload successful"})
}

// testing main
func main() {
	err := initializeDatabase()
	if err != nil {
		log.Println(err)
	}

	router := gin.Default()
	router.POST("/login", handleLogin)
	router.POST("/upload", handleUpload)
	router.GET("/members", handleListMembers)

	router.Run(":8080")
}
