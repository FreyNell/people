package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type person struct {
	ID   int64       `json:"id"`
	Name string      `json:"name"`
	Age  json.Number `json:"age"`
	Sex  json.Number `json:"sex"`
}

var DBNAME string = "people"

func createDB() (db *sql.DB) {

	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   os.Getenv("DBIP") + ":" + os.Getenv("DBPORT"),
	}

	log.Print(os.Getenv("DBIP") + ":" + os.Getenv("DBPORT"))
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + DBNAME)
	if err != nil {
		log.Fatal(err)
	}
	db.Close()

	cfg = mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   os.Getenv("DBIP") + ":" + os.Getenv("DBPORT"),
		DBName: DBNAME,
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

func getPeople(c *gin.Context) {
	db := createDB()

	var people []person
	rows, err := db.Query("SELECT id,name,age,sex FROM people;")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer rows.Close()
	defer db.Close()

	for rows.Next() {
		var per person
		if err := rows.Scan(&per.ID, &per.Name, &per.Age, &per.Sex); err != nil {
			log.Fatal(err)
			return
		}
		people = append(people, per)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
		return
	}

	c.IndentedJSON(http.StatusOK, people)
}

func postPeople(c *gin.Context) {
	var newPerson person

	if err := c.BindJSON(&newPerson); err != nil {
		log.Fatal(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An error has occur"})
	}

	db := createDB()
	defer db.Close()

	result, err := db.Exec("INSERT INTO people(name,age,sex) VALUES (?,?,?);", newPerson.Name, newPerson.Age, newPerson.Sex)

	if err != nil {
		log.Fatal(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An error has occur"})
	}

	newPerson.ID = id
	c.IndentedJSON(http.StatusOK, newPerson)
}

func getPersonByID(c *gin.Context) {
	id := c.Param("id")

	var per person

	db := createDB()
	row := db.QueryRow("SELECT id,name,age,sex FROM people WHERE id = ?", id)
	defer db.Close()
	if err := row.Scan(&per.ID, &per.Name, &per.Age, &per.Sex); err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Person not found."})
		}
		log.Fatal(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "An error has ocurr."})
	}
	c.IndentedJSON(http.StatusOK, per)

}

func defaultResponse(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "People Service")
}

func main() {

	router := gin.Default()

	log.Print(os.Getenv("ENV"))
	if os.Getenv("ENV") == "dev" {

		router.Use(cors.Default())
	}

	router.GET("/", defaultResponse)
	router.GET("/people", getPeople)
	router.GET("/people/:id", getPersonByID)
	router.POST("/people", postPeople)

	router.Run("0.0.0.0:8080")
}
