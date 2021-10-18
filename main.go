package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
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

func createDB(name string) (db *sql.DB) {

	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "192.168.0.42:3306",
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + name)
	if err != nil {
		log.Fatal(err)
	}
	db.Close()

	cfg = mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "192.168.0.42:3306",
		DBName: name,
	}
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS people(" +
		"id INT NOT NULL AUTO_INCREMENT PRIMARY KEY," +
		"name VARCHAR(30) NOT NULL," +
		"age INT," +
		"sex INT" +
		");")
	if err != nil {
		log.Fatal(err)
	}

	return
}

func main() {

	db := createDB("people")

	router := gin.Default()

	if os.Getenv("ENV") == "dev" {
		router.Use(cors.Default())
	}

	router.GET("/", defaultResponse)
	router.GET("/people", getPeople)
	router.GET("/people/:id", getPersonByID)
	router.POST("/people", postPeople)

	router.Run("0.0.0.0:8080")
}
