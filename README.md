# GoSimpleCRM
SQLite/Go/Gin/JavaScript Simple CRM

# Description: 
Simple CRM that creates a sqlite3 database file to store structs of Members. The React UI interface will display the table of members, allow for an upload of a csv file (not functional through testing & this assignment it is uploaded for us), update, searc, add, or delete members from the csv. 
The GO backend offers methods to interact with the SQL database
The Gin creates a RESTful API and acts as a middleware creating and sending JSON requests from the GO backend and the React UI
The React UI will be where the user interacts with the database and sends JSON requests through Gin.

# Dependencies -> Tech Stack
 
sqlite3 -> database, storage
gin -> Rest API, middleware
JavaScript React -> UI, frontend

# Packages
handlers/ --- holds all gin handler functions
members/ --- holds all go functions for sql database interaction
peopleobjs/ --- holds all structs (Members, Users, Address)

# Unit Testing
When testing run all tests together-- dependencies in the DB CSV upload. Easy fix if we want to test single functions.
Helpful to look at the members_test.csv file provided -- some test functions are working with this known data. If new testdata is added, test functions need to be updated
Automated testing files are located in packages: 
members/
Test with: go test ./members -v 
Integration Testing between Gin/React still needed. 

# Scaling Warnings: 
In test cases before scaling, refactor DB usage. Use :memory: usage instead of persistent shared .db file.
Add a parser to the .csv so any format can be uploaded -> will map headers to struct fields, this is hardcoded currently.
Ideally we would use different storage options than a local db file for larger datasets. (BigQuery, Cloud)

# Running the Program
**Open up the backend/ as root
under backend/ run go run main/main.go to start the server listening on port: 8080
under frontend/ run npm run dev to start the frontend react UI. Follow the url to reach the login page
Default testing login: 
Username = "admin"
Password = "admin"

UI Functionality is not yet complete. 
Use the corresponding unit cases for testing backend go functionality. 



