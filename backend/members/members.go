package members

//imports
import (
	"GOSIMPLECRM/backend/peopleobjs"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

//Main backend logic for Coalition App

//Database

// create the sql table to hold the member structs
func InitializeDatabase() error {
	db, err := sql.Open("sqlite3", "../coalition.db")
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

	// Insert default admin user if not exists
	var exists bool
	err = db.QueryRow(`SELECT EXISTS (SELECT 1 FROM users WHERE user_name = ?)`, "admin").Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking for default admin user: %v", err)
	}

	if !exists {
		_, err = db.Exec(`INSERT INTO users (user_name, password) VALUES (?, ?)`, "admin", "admin")
		if err != nil {
			return fmt.Errorf("error inserting default admin user: %v", err)
		}
		log.Println("Default admin user created: admin / admin")
	}

	//seed the database with sample data so that we have something to show on the UI during demos
	row := db.QueryRow("SELECT COUNT(*) FROM members")
	var count int
	err = row.Scan(&count)
	if err != nil {
		log.Println("Error checking Member count")
		return err
	} else {
		fmt.Printf("Current Member Count: %d", count)
	}

	//If the database is empty, load the sample .csv that is provided.

	if count == 0 {
		log.Println("Database is empty, seeding with sample data.")
		file, err := os.Open("members/testdata/members_test.csv")
		if err != nil {
			log.Println("Failed to open members_test.csv", err)
			return err
		}
		defer file.Close()

		err = UploadCSVList(file)
		if err != nil {
			log.Println("Failed to seed sample data.", err)
			return err
		}
	}
	log.Println("Sample data added to members table.")

	return nil
}

//Main features ---- Go Methods

//login authentication. The user name and password will be hardcoded for testing
//or displayed on the page for testing th feature. User: admin; pw: admin

func UserLogin(userName string, password string) bool {
	//check the username and password are not empty
	if userName == "" || password == "" {
		log.Println("Username or password cannot be empty.")
	}
	//open the database
	db, err := sql.Open("sqlite3", "../coalition.db")
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	var exists bool
	//use the database package to query the sqlite3 database
	err = db.QueryRow(`SELECT EXISTS (SELECT 1 FROM users WHERE user_name = ? AND password = ?)`, userName, password).Scan(&exists)
	if err != nil {
		log.Println(err)
		return false
	}
	//return true if user/pw is correct/found, false if not
	log.Printf("Login attempt for '%s': exists = %v\n", userName, exists)

	return exists
}

// let the user upload a new .csv file and use it to populate the database.
// not sure about this parameter yet. Also -- should it return a bool/error?
func UploadCSVList(filereader io.Reader) error {
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
			log.Println(err)
			failures++
			continue
		}
		//try to convert the row to a struct
		member, err := rowToStruct(row)
		if err != nil {
			log.Println(err)
			failures++
			continue
		}
		//add the new struct to the database
		err = AddNewMember(member)
		if err != nil {
			log.Println("Could not add new member", err)
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
func rowToStruct(row []string) (peopleobjs.Member, error) {
	//rows are just a slice of strings from a []string
	//create a new member struct and error to return.
	var m peopleobjs.Member
	var err error
	//access each row in the []string slice and map it to the correct field in the member struct
	m.Name = row[0]
	m.Address = peopleobjs.Address{
		Street: row[1],
		City:   row[2],
		State:  row[3],
		Zip:    row[4],
	}
	m.County = row[5]
	m.Phone = row[6]
	m.Email = row[7]

	//formating the date correctly
	m.Last_one_on_one, err = time.Parse("01/02/06", row[8])
	if err != nil {
		return m, err
	}
	//issues were saved as a comma separated string, so here we are splitting the string at each comma
	//to fill in the issues column.
	m.Issues = strings.Split(row[9], ",")
	m.Due_date_pay, err = time.Parse("01/02/06", row[10])
	if err != nil {
		return m, err
	}

	m.Active = strings.ToLower(row[11]) == "true" || row[11] == "1"
	return m, nil
}

// add a new member to the database. Here do we want to parse the json with the member request and put it
// in a member struct before we pass it here, or should we do it all in this method?
// not sure if I want to use a string[] or a member struct as the argument.
func AddNewMember(m peopleobjs.Member) error {
	db, err := sql.Open("sqlite3", "../coalition.db")
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
		m.Last_one_on_one.Format("01/02/06"),
		strings.Join(m.Issues, ","),
		m.Due_date_pay.Format("01/02/06"),
		m.Active,
	)
	if err != nil {
		return err
	}
	return nil
}

// return a list of all of the members (after updating, adding or deleting). Hand off to
// gin to display on the user interface.
func ListMembers() ([]peopleobjs.Member, error) {
	//create an array of Members to give back to gin
	var members []peopleobjs.Member
	db, err := sql.Open("sqlite3", "../coalition.db")
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
		var m peopleobjs.Member
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
		m.Last_one_on_one, _ = time.Parse("01/02/06", lastOneOnOne)
		m.Issues = strings.Split(issuesCSV, ",")
		m.Due_date_pay, _ = time.Parse("01/02/06", dueDatePay)
		//add the new member to the list of members
		members = append(members, m)
	}
	return members, nil
}

// return a list of all of the members that need to have one on one conversations.
// This means any inactive or member whose last convo was more than six months ago
func ListMembersNeedOneOnOnes() ([]peopleobjs.Member, error) {
	var members []peopleobjs.Member
	db, err := sql.Open("sqlite3", "../coalition.db")
	if err != nil {
		log.Println(err)
		return members, err
	}
	defer db.Close()

	//build the query
	query := `SELECT * FROM members WHERE active = false OR date(last_one_on_one) < date('now', '-6 months');`

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return members, err
	}
	defer rows.Close()

	//build the member
	for rows.Next() {
		var m peopleobjs.Member
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
		m.Last_one_on_one, _ = time.Parse("01/02/06", lastOneOnOne)
		m.Issues = strings.Split(issuesCSV, ",")
		m.Due_date_pay, _ = time.Parse("01/02/06", dueDatePay)

		members = append(members, m)
	}
	return members, nil
}

// delete a member from the sql member table. Given a member struct to delete.
func DeleteMember(id int) error {
	db, err := sql.Open("sqlite3", "../coalition.db")
	if err != nil {
		log.Println(err)
		return err
	}
	defer db.Close()

	query := `DELETE FROM members WHERE id = ?`
	_, err = db.Exec(query, id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// updates a member feild from the UI selections. Whole form is given for simplicity and we just replace all fields
func UpdateMember(m peopleobjs.Member) error {
	db, err := sql.Open("sqlite3", "../coalition.db")
	if err != nil {
		log.Println(err)
		return err
	}
	defer db.Close()

	query := `
	UPDATE members
	SET name = ?, street = ?, city = ?, state = ?, zip = ?, county = ?, 
		phone = ?, email = ?, last_one_on_one = ?, issues = ?, due_date_pay = ?, active = ?
	WHERE id = ?;
	`
	_, err = db.Exec(query,
		m.Name,
		m.Street,
		m.City,
		m.State,
		m.Zip,
		m.County,
		m.Phone,
		m.Email,
		m.Last_one_on_one.Format("01/02/06"),
		strings.Join(m.Issues, ","),
		m.Due_date_pay.Format("01/02/06"),
		m.Active,
		m.Id,
	)
	return err
}

// this needs to return a list of member structs where each member has the keyword given in the issues
// column of the database (comma separated string -> issues). User needs to choose from a drop down to avoid
// complexity for now. This can be expanded later.
func SearchVolunteers(keyword string) ([]peopleobjs.Member, error) {
	var members []peopleobjs.Member
	db, err := sql.Open("sqlite3", "../coalition.db")
	if err != nil {
		log.Println(err)
		return members, err
	}
	defer db.Close()

	//build the query
	query := `SELECT * FROM members WHERE issues LIKE ?`
	searchTerm := "%" + keyword + "%"

	rows, err := db.Query(query, searchTerm)
	if err != nil {
		log.Println(err)
		return members, err
	}
	defer rows.Close()

	for rows.Next() {
		var m peopleobjs.Member
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
		m.Last_one_on_one, _ = time.Parse("01/02/06", lastOneOnOne)
		m.Issues = strings.Split(issuesCSV, ",")
		m.Due_date_pay, _ = time.Parse("01/02/06", dueDatePay)

		members = append(members, m)
	}
	return members, nil
}

// returns a list of members that has an active status set to false.
func ListInactive() ([]peopleobjs.Member, error) {
	var members []peopleobjs.Member
	db, err := sql.Open("sqlite3", "../coalition.db")
	if err != nil {
		log.Println(err)
		return members, err
	}
	defer db.Close()

	//build the query
	query := `SELECT * FROM members WHERE active = false`

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return members, err
	}
	defer rows.Close()

	for rows.Next() {
		var m peopleobjs.Member
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
		m.Last_one_on_one, _ = time.Parse("01/02/06", lastOneOnOne)
		m.Issues = strings.Split(issuesCSV, ",")
		m.Due_date_pay, _ = time.Parse("01/02/06", dueDatePay)

		members = append(members, m)
	}
	return members, nil
}

// this helps with debugging the log in functionality. It checks and lists the entries in the users table
// so that we can add a test login
func DebugListUsers() {
	db, _ := sql.Open("sqlite3", "./coalition.db")
	defer db.Close()

	rows, err := db.Query("SELECT id, user_name, password FROM users")
	if err != nil {
		log.Println("Query failed:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var username, password string
		rows.Scan(&id, &username, &password)
		log.Printf("User: id=%d username=%s password=%s", id, username, password)
	}
}
