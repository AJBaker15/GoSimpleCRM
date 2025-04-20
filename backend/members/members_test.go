package members

//Assumptions: We are using the members_test.csv file provided
// Run the entire test suite together, not individual test functions. (Dependencies)

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

// We need this test case so we can "reset" the database after each test run. Updates and refactoring with DB
// need to happen here for scalability.
func TestMain(m *testing.M) {
	_ = os.Remove("../coalition.db")
	if err := InitializeDatabase(); err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}
	os.Exit(m.Run())
}

// This also tests the addMember and ListMembers functions in members.go in it's execution
func TestUploadCSVList(t *testing.T) {
	path := filepath.Join("testdata", "members_test.csv")
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("Failed to open test csv file: %v", err)
	}
	defer file.Close()

	err = InitializeDatabase()
	if err != nil {
		t.Fatalf("Failed to initialize the database: %v", err)
	}

	err = UploadCSVList(file)
	if err != nil {
		t.Fatalf("Failed to upload csv file: %v", err)
	}

	members, err := ListMembers()
	if err != nil {
		t.Fatalf("Failed to list members after upload: %v", err)
	}
	if len(members) == 0 {
		t.Fatalf("Expected at least one member after upload. Got: %v", err)
	}
}

func TestRowToStruct(t *testing.T) {
	row := []string{
		"John Doe",
		"123 Main St.", "Graham", "Ohio", "62704",
		"Caswell",
		"555-4597",
		"john@example.com",
		"01/01/01",
		"housing, education",
		"12/01/04",
		"true",
	}

	member, err := rowToStruct(row)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if member.Name != "John Doe" {
		t.Errorf("Expected 'John Doe' and got: %v", member.Name)
	}
	if len(member.Issues) != 2 || member.Issues[0] != "housing" {
		t.Errorf("Expected issue of two, including housing and got: %v", member.Issues)
	}
}

func TestSearchForVolunteers(t *testing.T) {
	err := InitializeDatabase()
	if err != nil {
		t.Fatalf("Could not initialize database: %v", err)
	}

	results, err := SearchVolunteers("education")
	if err != nil {
		t.Fatalf("Search Failed: %v", err)
	}

	if len(results) == 0 {
		t.Error("Expected to find a volunteer with education issue, found none.")
	}
}

func TestListInactive(t *testing.T) {
	err := InitializeDatabase()
	if err != nil {
		t.Fatalf("Could not initialize database: %v", err)
	}
	list, err := ListInactive()
	if err != nil {
		t.Fatalf("Search Failed: %v", err)
	}
	if len(list) == 0 {
		t.Error("Expected at least one inactive member, found none.")
	}
}

func TestDeleteMember(t *testing.T) {
	err := InitializeDatabase()
	if err != nil {
		t.Fatalf("Could not initialize database: %v", err)
	}
	//get all the members in the db and pick the first one to delete.
	members, _ := ListMembers()
	targetID := members[0].Id
	err = DeleteMember(targetID)
	if err != nil {
		t.Fatalf("Failed to delete member: %v", err)
	}
	//loop through the list looking for the targetID
	members, _ = ListMembers()
	for _, m := range members {
		if m.Id == targetID {
			t.Errorf("Member with id: %d found still in the list.", targetID)
		}
	}
}

func TestUpdateMember(t *testing.T) {
	err := InitializeDatabase()
	if err != nil {
		t.Fatalf("Could not initialize database: %v", err)
	}

	members, _ := ListMembers()
	if len(members) == 0 {
		t.Fatal("No members in database")
	}
	target := members[1]

	//Modify the user data to fields we know are not in the .csv test file.
	target.City = "Newtown"
	target.Phone = "555-2222"
	target.Active = false

	err = UpdateMember(target)
	if err != nil {
		t.Errorf("Failed to update member: %v", err)
	}

	//Fetch member to verify updated fields
	updatedMembers, _ := ListMembers()
	found := false
	for _, m := range updatedMembers {
		if m.Id == target.Id {
			found = true
			if m.City != "Newtown" || m.Phone != "555-2222" || m.Active != false {
				t.Errorf("Member update not applied correctly: %v", err)
			}
		}
	}
	if !found {
		t.Errorf("Updated member with ID %d not found", target.Id)
	}
}
