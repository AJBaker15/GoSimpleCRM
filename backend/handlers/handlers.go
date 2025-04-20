package handlers

import (
	"GOSIMPLECRM/backend/members"
	"GOSIMPLECRM/backend/peopleobjs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//GIN Handlers

// connect gin to the listMembers method by adding a handler
func HandleListMembers(c *gin.Context) {
	members, err := members.ListMembers()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve members"})
		return
	}
	c.IndentedJSON(http.StatusOK, members)
}

// Gin handler for making member struct from json. gin.H is an alias for a map of strings used to make jsons.
func HandleLogin(c *gin.Context) {
	//create a new user struct
	var u peopleobjs.User
	//try to use gin to map the json fields to the member structs using the BindJSON method
	if err := c.BindJSON(&u); err != nil {
		//if binding fails, send an error json
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
	}
	//call the go method userLogin to check the database for the un/pw
	if members.UserLogin(u.Username, u.Password) {
		//if found and correct, send a success message
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Login Successful"})
	} else {
		//if found and incorrect or not found send a failure message
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Invalid Credentials"})
	}
}

// Gin method to handle the upload and send messages back to the frontend
func HandleUpload(c *gin.Context) {
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
	err = members.UploadCSVList(f)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "CSV parsing failed"})
		return
	}
	//send a success message to the front end.
	c.IndentedJSON(http.StatusOK, gin.H{"message": "CSV upload successful"})
}

// handler to filter one on one list
func HandleListNeedOneOnOnes(c *gin.Context) {
	members, err := members.ListMembersNeedOneOnOnes()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Could not get filtered list"})
		return
	}
	c.IndentedJSON(http.StatusOK, members)
}

// handler to filter inactive list
func HandleListInactive(c *gin.Context) {
	members, err := members.ListInactive()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Could not get inactive list."})
	}
	c.IndentedJSON(http.StatusOK, members)
}

// handler to delete a member
func HandleDeleteMember(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	err = members.DeleteMember(id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Could not delete member."})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Member deleted successfully."})
}

// handler for searching for a keyword in issues list
func HandleSearchVolunteers(c *gin.Context) {
	keyword := c.Query("issue")

	if keyword == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Missing keyword."})
		return
	}

	members, err := members.SearchVolunteers(keyword)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Search Failed."})
		return
	}
	c.IndentedJSON(http.StatusOK, members)
}

// handler for updating a member
func HandleUpdateMember(c *gin.Context) {
	var m peopleobjs.Member
	if err := c.BindJSON(&m); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id != m.Id {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Mismatch or invalid ID"})
		return
	}

	if err := members.UpdateMember(m); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Could not update member"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Member updated successfully"})
}
