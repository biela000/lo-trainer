package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TrainingSet struct {
	id, p1Id, p2Id, p3Id, p4Id, p5Id int
}

type TrainingSetController struct {
	db *sql.DB
}

type TrainingSetCriteria struct {
	levels, formats, subjects            []string
	minScore, maxScore, minYear, maxYear int
}

func (controller *TrainingSetController) CreateTrainingSet(c *gin.Context) {
	trainingSet := TrainingSet{}
	trainingSetCriteria := TrainingSetCriteria{}

	trainingSetCriteria.levels = c.QueryArray("levels")
	if len(trainingSetCriteria.levels) == 0 {
		trainingSetCriteria.levels = DefaultLevels[:]
	}

	trainingSetCriteria.formats = c.QueryArray("formats")
	if len(trainingSetCriteria.formats) == 0 {
		trainingSetCriteria.formats = DefaultFormats[:]
	}

	trainingSetCriteria.subjects = c.QueryArray("subjects")
	if len(trainingSetCriteria.subjects) == 0 {
		trainingSetCriteria.subjects = DefaultSubjects[:]
	}

	var err error
	trainingSetCriteria.minScore, err = strconv.Atoi(c.DefaultQuery("min_score", "0"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	trainingSetCriteria.maxScore, err = strconv.Atoi(c.DefaultQuery("min_score", "100"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	trainingSetCriteria.minYear, err = strconv.Atoi(c.DefaultQuery("min_year", "0"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	trainingSetCriteria.maxYear, err = strconv.Atoi(c.DefaultQuery("max_year", "2049"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	tx, err := controller.db.Begin()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	controller.getPuzzlesForTrainingSet(&trainingSet, trainingSetCriteria)

	stmt, err := tx.Prepare(`
		INSERT INTO Training_Sets (
			p1_id,
			p2_id,
			p3_id,
			p4_id,
			p5_id
		) VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer stmt.Close()

	if _, err = stmt.Exec(
		trainingSet.p1Id,
		trainingSet.p2Id,
		trainingSet.p3Id,
		trainingSet.p4Id,
		trainingSet.p5Id,
	); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	tx.Commit()

	c.JSON(http.StatusOK, trainingSet)
}

func (controller *TrainingSetController) getPuzzlesForTrainingSet(
	trainingSet *TrainingSet, criteria TrainingSetCriteria,
) error {
	stmt, err := controller.db.Prepare(`
		SELECT id FROM Puzzles
		WHERE level IN (?)
		AND format IN (?)
		AND subject IN (?)
		AND score >= ?
		AND score <= ?
		AND year >= ?
		AND year <= ?
		ORDER BY RANDOM()
		LIMIT 5;
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(
		criteria.levels,
		criteria.formats,
		criteria.subjects,
		criteria.minScore,
		criteria.maxScore,
		criteria.minYear,
		criteria.maxYear,
	)

	puzzleIds := [5]*int{
		&trainingSet.p1Id,
		&trainingSet.p2Id,
		&trainingSet.p3Id,
		&trainingSet.p4Id,
		&trainingSet.p5Id,
	}

	for _, puzzleId := range puzzleIds {
		if err = rows.Scan(puzzleId); err != nil {
			return err
		}
		if exists := rows.Next(); !exists {
			break
		}
	}
	return nil
}
