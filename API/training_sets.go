package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type TrainingSet struct {
	Id   int `json:"id"`
	P1Id int `json:"p1_id"`
	P2Id int `json:"p2_id"`
	P3Id int `json:"p3_id"`
	P4Id int `json:"p4_id"`
	P5Id int `json:"p5_id"`
}

type TrainingSetController struct {
	db *sql.DB
}

type TrainingSetCriteria struct {
	levels, formats, subjects            []string
	minScore, maxScore, minYear, maxYear int
}

func (controller *TrainingSetController) GetTrainingSets(c *gin.Context) {
	trainingSets := []TrainingSet{}
	rows, err := controller.db.Query(`
		SELECT *
		FROM Training_Sets;
	`)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		trainingSet := TrainingSet{}
		err = rows.Scan(
			&trainingSet.Id,
			&trainingSet.P1Id,
			&trainingSet.P2Id,
			&trainingSet.P3Id,
			&trainingSet.P4Id,
			&trainingSet.P5Id,
		)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		trainingSets = append(trainingSets, trainingSet)
	}
	c.JSON(http.StatusOK, trainingSets)
}

func (controller *TrainingSetController) GetTrainingSet(c *gin.Context) {
	trainingSet := TrainingSet{}
	trainingSetId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	stmt, err := controller.db.Prepare(`
		SELECT *
		FROM Training_Sets
		WHERE id = ?;
	`)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(trainingSetId).Scan(
		&trainingSet.Id,
		&trainingSet.P1Id,
		&trainingSet.P2Id,
		&trainingSet.P3Id,
		&trainingSet.P4Id,
		&trainingSet.P5Id,
	)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, trainingSet)
}

func (controller *TrainingSetController) CreateTrainingSet(c *gin.Context) {
	trainingSet := TrainingSet{Id: -1}
	trainingSetCriteria := TrainingSetCriteria{}

	trainingSetCriteria.getFromQuery(c)

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
		RETURNING id;
	`)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(
		trainingSet.P1Id,
		trainingSet.P2Id,
		trainingSet.P3Id,
		trainingSet.P4Id,
		trainingSet.P5Id,
	).Scan(&trainingSet.Id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	tx.Commit()

	fmt.Println("Training Set Created")
	fmt.Println(trainingSet)

	c.JSON(http.StatusOK, trainingSet)
}

func (trainingSetCriteria *TrainingSetCriteria) getFromQuery(c *gin.Context) {
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

	trainingSetCriteria.maxScore, err = strconv.Atoi(c.DefaultQuery("max_score", "100"))
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
}

func (controller *TrainingSetController) getPuzzlesForTrainingSet(
	trainingSet *TrainingSet, criteria TrainingSetCriteria,
) error {
	queryString, queryArgs := createPuzzleQuery(criteria)

	stmt, err := controller.db.Prepare(queryString)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(queryArgs...)
	if err != nil {
		return err
	}
	defer rows.Close()
	rows.Next()

	puzzleIds := [5]*int{
		&trainingSet.P1Id,
		&trainingSet.P2Id,
		&trainingSet.P3Id,
		&trainingSet.P4Id,
		&trainingSet.P5Id,
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

func createPuzzleQuery(criteria TrainingSetCriteria) (string, []interface{}) {
	builder := strings.Builder{}
	queryArgs := []interface{}{}

	for i := range criteria.levels {
		queryArgs = append(queryArgs, criteria.levels[i])
	}

	for i := range criteria.formats {
		queryArgs = append(queryArgs, criteria.formats[i])
	}

	for i := range criteria.subjects {
		builder.WriteString("subjects LIKE '%' || ? || '%'")
		if i < len(criteria.subjects)-1 {
			builder.WriteString(" OR ")
		}
		queryArgs = append(queryArgs, criteria.subjects[i])
	}

	subjectsCondition := builder.String()

	levelsQueryPlaceholders := strings.Repeat("?, ", len(criteria.levels)-1) + "?"
	formatsQueryPlaceholders := strings.Repeat("?, ", len(criteria.formats)-1) + "?"

	queryString := `
		SELECT id FROM Puzzles
		WHERE level IN (` + levelsQueryPlaceholders + `)
		AND format IN (` + formatsQueryPlaceholders + `) AND
	(` + subjectsCondition + `)
		AND score >= ?
		AND score <= ?
		AND year >= ?
		AND year <= ?
		AND id NOT IN (SELECT p1_id FROM Training_Sets)
		AND id NOT IN (SELECT p2_id FROM Training_Sets)
		AND id NOT IN (SELECT p3_id FROM Training_Sets)
		AND id NOT IN (SELECT p4_id FROM Training_Sets)
		AND id NOT IN (SELECT p5_id FROM Training_Sets)
		ORDER BY RANDOM()
		LIMIT 5;
	`
	queryArgs = append(
		queryArgs,
		criteria.minScore,
		criteria.maxScore,
		criteria.minYear,
		criteria.maxYear,
	)

	return queryString, queryArgs
}
