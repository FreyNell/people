package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type person struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int64  `json:"age"`
	Sex  string `json:"sex"`
}

// starting data.
var people = []person{
	{ID: "1", Name: "Erick", Age: 25, Sex: "male"},
	{ID: "2", Name: "Stefanny", Age: 25, Sex: "female"},
	{ID: "3", Name: "George", Age: 23, Sex: "male"},
}

func getPeople(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, people)
}

func postPeople(c *gin.Context) {
	var newPerson person

	if err := c.BindJSON(&newPerson); err != nil {
		return
	}

	people = append(people, newPerson)
	c.IndentedJSON(http.StatusCreated, newPerson)
}

func getPersonByID(c *gin.Context) {
	id := c.Param("id")

	for _, a := range people {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "person not found."})
}

func defaultResponse(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "People Service")
}

func main() {
	router := gin.Default()
	router.GET("/", defaultResponse)
	router.GET("/people", getPeople)
	router.GET("/people/:id", getPersonByID)
	router.POST("/people", postPeople)

	router.Run("0.0.0.0:8082")
}
