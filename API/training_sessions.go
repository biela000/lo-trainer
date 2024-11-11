package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TrainingSession struct {
	Id            int `json:"id"`
	TrainingSetId int `json:"training_set_id"`
	P1Score       int `json:"p1_score"`
	P2Score       int `json:"p2_score"`
	P3Score       int `json:"p3_score"`
	P4Score       int `json:"p4_score"`
	P5Score       int `json:"p5_score"`
	FullScore     int `json:"full_score"`
	P1Time        int `json:"p1_time"`
	P2Time        int `json:"p2_time"`
	P3Time        int `json:"p3_time"`
	P4Time        int `json:"p4_time"`
	P5Time        int `json:"p5_time"`
	FullTime      int `json:"full_time"`
	DateTaken     int `json:"date_taken"`
	Finished      int `json:"finished"`
}

type TrainingSessionController struct {
	db *sql.DB
}

func (controller *TrainingSessionController) GetTrainingSessions(c *gin.Context) {
	trainingSessions := []TrainingSession{}
	rows, err := controller.db.Query(`
		SELECT *
		FROM Training_Sessions;
	`)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		trainingSession := TrainingSession{}
		err = rows.Scan(
			&trainingSession.Id,
			&trainingSession.TrainingSetId,
			&trainingSession.P1Score,
			&trainingSession.P2Score,
			&trainingSession.P3Score,
			&trainingSession.P4Score,
			&trainingSession.P5Score,
			&trainingSession.FullScore,
			&trainingSession.P1Time,
			&trainingSession.P2Time,
			&trainingSession.P3Time,
			&trainingSession.P4Time,
			&trainingSession.P5Time,
			&trainingSession.FullTime,
			&trainingSession.DateTaken,
			&trainingSession.Finished,
		)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		trainingSessions = append(trainingSessions, trainingSession)
	}
	c.JSON(http.StatusOK, trainingSessions)
}

func (controller *TrainingSessionController) GetTrainingSession(c *gin.Context) {
	trainingSession := TrainingSession{}
	trainingSessionId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	stmt, err := controller.db.Prepare(`
		SELECT *
		FROM Training_Sessions
		WHERE id = ?;
	`)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(trainingSessionId).Scan(
		&trainingSession.Id,
		&trainingSession.TrainingSetId,
		&trainingSession.P1Score,
		&trainingSession.P2Score,
		&trainingSession.P3Score,
		&trainingSession.P4Score,
		&trainingSession.P5Score,
		&trainingSession.FullScore,
		&trainingSession.P1Time,
		&trainingSession.P2Time,
		&trainingSession.P3Time,
		&trainingSession.P4Time,
		&trainingSession.P5Time,
		&trainingSession.FullTime,
		&trainingSession.DateTaken,
		&trainingSession.Finished,
	)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, trainingSession)
}

func (controller *TrainingSessionController) CreateTrainingSession(c *gin.Context) {
	trainingSession := TrainingSession{}
	err := c.BindJSON(&trainingSession)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	tx, err := controller.db.Begin()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	stmt, err := tx.Prepare(`
		INSERT INTO Training_Sessions (
			training_set_id,
			p1_score,
			p2_score,
			p3_score,
			p4_score,
			p5_score,
			full_score,
			p1_time,
			p2_time,
			p3_time,
			p4_time,
			p5_time,
			full_time,
			finished
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id, date_taken;
	`)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(
		trainingSession.TrainingSetId,
		trainingSession.P1Score,
		trainingSession.P2Score,
		trainingSession.P3Score,
		trainingSession.P4Score,
		trainingSession.P5Score,
		trainingSession.FullScore,
		trainingSession.P1Time,
		trainingSession.P2Time,
		trainingSession.P3Time,
		trainingSession.P4Time,
		trainingSession.P5Time,
		trainingSession.FullTime,
		trainingSession.Finished,
	).Scan(&trainingSession.Id, &trainingSession.DateTaken)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tx.Commit()

	fmt.Println("Created training session with id:", trainingSession.Id)

	c.JSON(http.StatusOK, trainingSession)
}

func (controller *TrainingSessionController) UpdateTrainingSession(c *gin.Context) {
	trainingSession := TrainingSession{}
	trainingSessionId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	trainingSession.Id = trainingSessionId
	err = c.BindJSON(&trainingSession)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	tx, err := controller.db.Begin()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	stmt, err := tx.Prepare(`
		UPDATE Training_Sessions
		SET
			training_set_id = ?,
			p1_score = ?,
			p2_score = ?,
			p3_score = ?,
			p4_score = ?,
			p5_score = ?,
			full_score = ?,
			p1_time = ?,
			p2_time = ?,
			p3_time = ?,
			p4_time = ?,
			p5_time = ?,
			full_time = ?,
			finished = ?
		WHERE id = ?
		RETURNING date_taken;
	`)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		trainingSession.TrainingSetId,
		trainingSession.P1Score,
		trainingSession.P2Score,
		trainingSession.P3Score,
		trainingSession.P4Score,
		trainingSession.P5Score,
		trainingSession.FullScore,
		trainingSession.P1Time,
		trainingSession.P2Time,
		trainingSession.P3Time,
		trainingSession.P4Time,
		trainingSession.P5Time,
		trainingSession.FullTime,
		trainingSession.Finished,
		trainingSession.Id,
	).Scan(&trainingSession.DateTaken)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	tx.Commit()
	fmt.Println("Updated training session with id:", trainingSession.Id)
	c.JSON(http.StatusOK, trainingSession)
}

func (controller *TrainingSessionController) DeleteTrainingSession(c *gin.Context) {
	trainingSessionId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	tx, err := controller.db.Begin()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	stmt, err := tx.Prepare(`
		DELETE FROM Training_Sessions
		WHERE id = ?;
	`)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(trainingSessionId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	tx.Commit()
	fmt.Println("Deleted training session with id:", trainingSessionId)

	c.JSON(http.StatusNoContent, nil)
}
