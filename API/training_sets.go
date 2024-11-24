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
	P1Id int `json:"p1Id"`
	P2Id int `json:"p2Id"`
	P3Id int `json:"p3Id"`
	P4Id int `json:"p4Id"`
	P5Id int `json:"p5Id"`
}

type PopulatedTrainingSet struct {
	Id     int    `json:"id"`
	P1Link string `json:"p1Link"`
	P2Link string `json:"p2Link"`
	P3Link string `json:"p3Link"`
	P4Link string `json:"p4Link"`
	P5Link string `json:"p5Link"`
}

type NullablePopulatedTrainingSet struct {
	Id     int            `json:"id"`
	P1Link sql.NullString `json:"p1Link"`
	P2Link sql.NullString `json:"p2Link"`
	P3Link sql.NullString `json:"p3Link"`
	P4Link sql.NullString `json:"p4Link"`
	P5Link sql.NullString `json:"p5Link"`
}

type TrainingSetController struct {
	db *sql.DB
}

type TrainingSetCriteria struct {
	levels, formats, subjects            []string
	minScore, maxScore, minYear, maxYear int
}

func (controller *TrainingSetController) GetTrainingSets(c *gin.Context) {
	trainingSets := []PopulatedTrainingSet{}
	rows, err := controller.db.Query(`
		SELECT 
			t.id,
			p1.link AS p1_link,
			p2.link AS p2_link,
			p3.link AS p3_link,
			p4.link AS p4_link,
			p5.link AS p5_link
		FROM 
			Training_Sets t
		LEFT JOIN 
			puzzles p1 ON t.p1_id = p1.id
		LEFT JOIN 
			puzzles p2 ON t.p2_id = p2.id
		LEFT JOIN 
			puzzles p3 ON t.p3_id = p3.id
		LEFT JOIN 
			puzzles p4 ON t.p4_id = p4.id
		LEFT JOIN 
			puzzles p5 ON t.p5_id = p5.id;
	`)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		trainingSet := NullablePopulatedTrainingSet{}
		err = rows.Scan(
			&trainingSet.Id,
			&trainingSet.P1Link,
			&trainingSet.P2Link,
			&trainingSet.P3Link,
			&trainingSet.P4Link,
			&trainingSet.P5Link,
		)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		trainingSetResult := createPopulatedTrainingSet(trainingSet)

		trainingSets = append(trainingSets, trainingSetResult)
	}
	c.JSON(http.StatusOK, trainingSets)
}

func (controller *TrainingSetController) GetTrainingSet(c *gin.Context) {
	trainingSet := NullablePopulatedTrainingSet{}
	trainingSetId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	stmt, err := controller.db.Prepare(`
		SELECT 
			t.id,
			p1.link AS p1_link,
			p2.link AS p2_link,
			p3.link AS p3_link,
			p4.link AS p4_link,
			p5.link AS p5_link
		FROM 
			Training_Sets t
		LEFT JOIN 
			puzzles p1 ON t.p1_id = p1.id
		LEFT JOIN 
			puzzles p2 ON t.p2_id = p2.id
		LEFT JOIN 
			puzzles p3 ON t.p3_id = p3.id
		LEFT JOIN 
			puzzles p4 ON t.p4_id = p4.id
		LEFT JOIN 
			puzzles p5 ON t.p5_id = p5.id
		WHERE
			t.id = ?;
	`)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(trainingSetId).Scan(
		&trainingSet.Id,
		&trainingSet.P1Link,
		&trainingSet.P2Link,
		&trainingSet.P3Link,
		&trainingSet.P4Link,
		&trainingSet.P5Link,
	)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	trainingSetResult := createPopulatedTrainingSet(trainingSet)
	c.JSON(http.StatusOK, trainingSetResult)
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

	populateStmt, err := controller.db.Prepare(`
		SELECT 
			t.id,
			p1.link AS p1_link,
			p2.link AS p2_link,
			p3.link AS p3_link,
			p4.link AS p4_link,
			p5.link AS p5_link
		FROM 
			Training_Sets t
		LEFT JOIN 
			puzzles p1 ON t.p1_id = p1.id
		LEFT JOIN 
			puzzles p2 ON t.p2_id = p2.id
		LEFT JOIN 
			puzzles p3 ON t.p3_id = p3.id
		LEFT JOIN 
			puzzles p4 ON t.p4_id = p4.id
		LEFT JOIN 
			puzzles p5 ON t.p5_id = p5.id
		WHERE
			t.id = ?;
	`)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer populateStmt.Close()

	populatedTrainingSet := NullablePopulatedTrainingSet{}

	err = populateStmt.QueryRow(trainingSet.Id).Scan(
		&populatedTrainingSet.Id,
		&populatedTrainingSet.P1Link,
		&populatedTrainingSet.P2Link,
		&populatedTrainingSet.P3Link,
		&populatedTrainingSet.P4Link,
		&populatedTrainingSet.P5Link,
	)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	populatedTrainingSetResult := createPopulatedTrainingSet(populatedTrainingSet)

	c.JSON(http.StatusOK, populatedTrainingSetResult)
}

func (controller *TrainingSetController) DeleteTrainingSet(c *gin.Context) {
	trainingSetId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	stmt, err := controller.db.Prepare(`
		DELETE FROM Training_Sets
		WHERE id = ?;
	`)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(trainingSetId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
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

func createPopulatedTrainingSet(trainingSet NullablePopulatedTrainingSet) PopulatedTrainingSet {
	return PopulatedTrainingSet{
		Id:     trainingSet.Id,
		P1Link: trainingSet.P1Link.String,
		P2Link: trainingSet.P2Link.String,
		P3Link: trainingSet.P3Link.String,
		P4Link: trainingSet.P4Link.String,
		P5Link: trainingSet.P5Link.String,
	}
}
