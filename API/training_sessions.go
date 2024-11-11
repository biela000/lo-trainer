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
		&trainingSession.Finished,
	)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, trainingSession)
}
