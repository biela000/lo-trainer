package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	_ "github.com/mattn/go-sqlite3"
)

type Puzzle struct {
	name, number, link, level, subjects, format string
	year, score                                 int
}

func main() {
	collector := colly.NewCollector(
		colly.AllowedDomains("www.uklo.org"),
	)

	defer collector.Visit("https://www.uklo.org/past-exam-papers")

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL)
	})

	collector.OnError(func(response *colly.Response, err error) {
		fmt.Println("Error:", err)
	})

	var puzzles []Puzzle

	collector.OnHTML("#tablepress-14 > tbody", func(table *colly.HTMLElement) {
		table.ForEach("tr", func(_ int, row *colly.HTMLElement) {
			puzzle := Puzzle{}

			err := puzzle.getNameNumberYearFromColumn(
				row.ChildText("td.column-2 > a"),
			)
			if err != nil {
				log.Fatal(err)
			}

			puzzle.link = row.ChildAttr("td.column-2 > a", "href")

			err = puzzle.getLevelScoreFromColumn(
				row.ChildText("td.column-1"),
			)
			if err != nil {
				log.Fatal(err)
			}

			puzzle.format = row.ChildText("td.column-4")

			puzzle.getSubjectsFromColumn(row.ChildText("td.column-3"))

			puzzles = append(puzzles, puzzle)
		})
	})

	collector.OnScraped(func(response *colly.Response) {
		fmt.Println("Name | Number | Link | Year | Level | Score | Subjects | Format")
		fmt.Println("-----------------------------------------------------------------")
		for _, puzzle := range puzzles {
			fmt.Printf(
				"%s | %s | %s | %d | %s | %d | %s | %s\n",
				puzzle.name,
				puzzle.number,
				puzzle.link,
				puzzle.year,
				puzzle.level,
				puzzle.score,
				puzzle.subjects,
				puzzle.format,
			)
		}
		fmt.Println("Total puzzles:", len(puzzles))

		os.Remove("./puzzles.db")

		db, err := sql.Open("sqlite3", "./puzzles.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		sqlStmt := `
		CREATE TABLE puzzles (id integer not null primary key, name text, year integer, level text, number text, subjects text, format text, score integer, link text);
		DELETE FROM puzzles;
	`

		_, err = db.Exec(sqlStmt)
		if err != nil {
			log.Printf("%q: %s\n", err, sqlStmt)
			return
		}

		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}

		stmt, err := tx.Prepare("INSERT INTO puzzles(name, year, level, number, subjects, format, score, link) VALUES(?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		for _, puzzle := range puzzles {
			_, err = stmt.Exec(
				puzzle.name,
				puzzle.year,
				puzzle.level,
				puzzle.number,
				puzzle.subjects,
				puzzle.format,
				puzzle.score,
				puzzle.link,
			)
			if err != nil {
				log.Fatal(err)
			}
		}

		err = tx.Commit()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Database created and populated successfully")
	})
}

func (puzzle *Puzzle) getNameNumberYearFromColumn(columnText string) error {
	puzzleYearAndNumber := columnText[:strings.IndexByte(columnText, ' ')]

	puzzle.name = columnText[strings.IndexByte(columnText, ' ')+1:]

	puzzle.number = puzzleYearAndNumber[strings.LastIndexByte(puzzleYearAndNumber, '_')+1:]

	var err error
	puzzle.year, err = strconv.Atoi(puzzleYearAndNumber[:4])
	if err != nil {
		return err
	}

	return nil
}

func (puzzle *Puzzle) getLevelScoreFromColumn(columnText string) error {
	re := regexp.MustCompile(`^[a-zA-Z\s\/]*`)
	puzzle.level = re.FindStringSubmatch(columnText)[0]

	if strings.Contains(columnText, "No Data") {
		puzzle.score = -1
		return nil
	}

	score1, err := strconv.Atoi(
		columnText[strings.IndexByte(columnText, '[')+1 : strings.IndexByte(columnText, '%')],
	)
	if err != nil {
		return err
	}

	lastSlashIndex := strings.LastIndexByte(columnText, '/')

	if lastSlashIndex == -1 {
		puzzle.score = score1
	} else {
		score2, err := strconv.Atoi(columnText[lastSlashIndex+2 : len(columnText)-2])
		if err != nil {
			return err
		}
		puzzle.score = (score1 + score2) / 2
	}

	return nil
}

func (puzzle *Puzzle) getSubjectsFromColumn(columnText string) {
	subjects := map[string]string{
		"Co": "Compounding",
		"Mo": "Morphology",
		"Nu": "Numbers",
		"Ph": "Phonology and Phonetics",
		"Se": "Semantics",
		"Sy": "Syntax",
		"Wr": "Writing System",
	}

	var builder strings.Builder

	for i := 2; i < len(columnText)-1; i += 3 {
		builder.WriteString(subjects[columnText[i:i+2]])

		if i+3 < len(columnText) {
			builder.WriteString(", ")
		}
	}

	puzzle.subjects = builder.String()
}
