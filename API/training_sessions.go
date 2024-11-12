package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TrainingSession struct {
	Id            int       `json:"id"`
	TrainingSetId int       `json:"trainingSetId"`
	P1Score       int       `json:"p1Score"`
	P2Score       int       `json:"p2Score"`
	P3Score       int       `json:"p3Score"`
	P4Score       int       `json:"p4Score"`
	P5Score       int       `json:"p5Score"`
	FullScore     int       `json:"fullScore"`
	P1Time        int       `json:"p1Time"`
	P2Time        int       `json:"p2Time"`
	P3Time        int       `json:"p3Time"`
	P4Time        int       `json:"p4Time"`
	P5Time        int       `json:"p5Time"`
	FullTime      int       `json:"fullTime"`
	DateTaken     time.Time `json:"dateTaken"`
	Finished      bool      `json:"finished"`
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
		) VALUES (?, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
		RETURNING *;
	`)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(
		trainingSession.TrainingSetId,
	).Scan(
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

func (controller *TrainingSessionController) GetTrainingSessionStreak(c *gin.Context) {
	query := `
		WITH ordered_dates AS (
			SELECT 
				id,
				date_taken,
				finished,
				ROW_NUMBER() OVER (ORDER BY date_taken DESC) AS row_num
			FROM 
				Training_Sessions
			WHERE 
				finished = 1
		),
		acceptable_latest_date AS (
			SELECT 
				date_taken AS latest_date
			FROM 
				ordered_dates
			WHERE 
				date_taken IN (DATE('now'), DATE('now', '-1 day'))
			ORDER BY 
				date_taken DESC
           		LIMIT 1
		),
		consecutive_days AS (
			SELECT 
				id,
				date_taken,
				finished,
				DATE(date_taken, '-' || (ROW_NUMBER() OVER (ORDER BY date_taken DESC)) || ' day') AS consecutive_group
            		FROM 
				ordered_dates
		)
		SELECT 
			COUNT(*)
		FROM 
			consecutive_days
		WHERE 
			consecutive_group = (
				SELECT 
					DATE(latest_date, '-' || (ROW_NUMBER() OVER (ORDER BY latest_date DESC)) || ' day')
				FROM 
					acceptable_latest_date
			);
	`

	var count int
	err := controller.db.QueryRow(query).Scan(&count)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	fmt.Println("Streak:", count)
	c.JSON(http.StatusOK, map[string]int{"streak": count})
}
