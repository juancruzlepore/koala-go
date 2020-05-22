package main

import (
	. "database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

type Movie struct {
	CreationDate time.Time `json:"creationDate"`
	AddedBy      string    `json:"addedBy"`
	Name         string    `json:"name"`
	Seen         bool      `json:"seen"`
}

type MovieList struct {
	Movies []Movie `json:"movies"`
}

func connect() *DB {
	const (
		//host             = "ec2-174-129-206-173.compute-1.amazonaws.com"
		//dbPort           = 5432
		//user             = "bzmsjvfflczcup"
		//password         = "51e4f5708d04bfd7a3513573b53018e5ce45e7664605ef800e6f8cc2098388b8"
		//dbname           = "dd9vnmb96ok7ie"
		connectionString = "dbname=dd9vnmb96ok7ie host=ec2-174-129-206-173.compute-1.amazonaws.com port=5432 user=bzmsjvfflczcup password=51e4f5708d04bfd7a3513573b53018e5ce45e7664605ef800e6f8cc2098388b8 sslmode=require"
	)
	fmt.Println("Connecting...")

	db, err := Open("postgres", connectionString)
	dbInstance := db
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return dbInstance
}

func getNextDate(db *DB) string {
	rows, err := db.Query("SELECT start_date FROM next_dates WHERE end_date > current_timestamp ORDER BY end_date ASC LIMIT 1")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	date := "None"
	for rows.Next() {
		if err := rows.Scan(&date); err != nil {
			log.Fatal(err)
		}
		fmt.Println(date)
	}
	return date
}

func getMovies(db *DB) string {
	rows, err := db.Query("SELECT * FROM movies")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var movies []Movie
	for rows.Next() {
		newMovie := Movie{}
		if err := rows.Scan(&newMovie.CreationDate, &newMovie.AddedBy, &newMovie.Name, &newMovie.Seen); err != nil {
			log.Fatal(err)
		} else {
			movies = append(movies, newMovie)
		}
	}
	var movieList MovieList
	movieList.Movies = movies
	moviesJson, err := json.Marshal(movieList)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	fmt.Println(len(movies))
	return string(moviesJson[:])
}

func addMovie(db *DB, movie Movie) bool {
	movieName := strings.Replace(movie.Name, "'", "\\'", -1)
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM movies where name='%s'", movieName))
	defer rows.Close()
	if rows.Next() {
		return false
	}

	_, err = db.Exec(fmt.Sprintf("insert into movies values (current_timestamp, '%s', '%s', false);", movie.AddedBy, movieName))

	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func addDate(db *DB, dateStart string, dateEnd string) bool {
	dateStart = strings.Replace(dateStart, "'", "\\'", -1)
	dateEnd = strings.Replace(dateEnd, "'", "\\'", -1)

	rows, err := db.Query(fmt.Sprintf("SELECT * FROM next_dates where start_date='%s'", dateStart))
	//defer rows.Close()

	if err == nil && rows.Next() {
		log.Fatal("start date already existing")
		return false
	}

	_, err = db.Exec("insert into next_dates values ($1, $2);", dateStart, dateEnd)
	if err != nil {
		log.Fatal(err)
		return false
	}

	return true
}
