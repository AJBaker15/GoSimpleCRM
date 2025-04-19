package main

//imports
import (
	"GOSIMPLECRM/handlers"
	"GOSIMPLECRM/members"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// testing main
func main() {
	err := members.InitializeDatabase()
	if err != nil {
		log.Println(err)
	}

	router := gin.Default()
	router.POST("/login", handlers.HandleLogin)
	router.POST("/upload", handlers.HandleUpload)
	router.GET("/members", handlers.HandleListMembers)
	router.GET("/members/need-one-on-ones", handlers.HandleListNeedOneOnOnes)
	router.GET("/members/inactive", handlers.HandleListInactive)
	router.GET("/members/search", handlers.HandleSearchVolunteers)
	router.DELETE("/member/:id", handlers.HandleDeleteMember)
	router.PUT("/member/:id", handlers.HandleUpdateMember)

	router.Run(":8080")
}
