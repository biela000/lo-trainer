package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	router := gin.Default()

	router.Use(ErrorHandler)

	db, err := sql.Open("sqlite3", "./lo_trainer.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTables(db)

	trainingSetController := TrainingSetController{db: db}

	v1 := router.Group("/v1")
	{
		v1.POST("/training_sets", trainingSetController.CreateTrainingSet)
	}

	router.Run(":8080")
}

func createTables(db *sql.DB) error {
	sqlStmt := `
		CREATE TABLE IF NOT EXISTS Training_Sets (id integer not null primary key, p1_id integer, p2_id integer, p3_id integer, p4_id integer, p5_id integer);
		DELETE FROM Training_Sets;
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	sqlStmt = `
		CREATE TABLE IF NOT EXISTS Scores (id integer not null primary key, training_set_id integer, p1_score integer, p2_score integer, p3_score integer, p4_score integer, p5_score integer, full_score integer);
		DELETE FROM Scores;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	sqlStmt = `
		CREATE TABLE IF NOT EXISTS Times (id integer not null primary key, training_set_id integer, p1_time integer, p2_time integer, p3_time integer, p4_time integer, p5_time integer, full_time integer);
		DELETE FROM Times;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return nil
}

func ErrorHandler(c *gin.Context) {
	c.Next()

	for _, err := range c.Errors {
		log.Println(err)
	}

	c.JSON(-1, gin.H{"error": "An error occurred"})
}
