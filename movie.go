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
	Rating       float32   `json:"rating"`
	Genres       string    `json:"genres"`
}

type MovieList struct {
	Movies []Movie `json:"movies"`
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
		if err := rows.Scan(
			&newMovie.CreationDate,
			&newMovie.AddedBy,
			&newMovie.Name,
			&newMovie.Seen,
			&newMovie.Rating,
			&newMovie.Genres); err != nil {
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
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM movies WHERE name='%s'", movieName))
	defer rows.Close()
	if rows.Next() {
		_, err = db.Exec(
			fmt.Sprintf("UPDATE movies SET added_by = '%s', seen = %t, rating = %f, genres = %s where name = '%s';",
				movie.AddedBy, movie.Seen, movie.Rating, movie.Genres, movieName))
		if err != nil {
			log.Fatal(err)
			return false
		}
		return true
	}

	_, err = db.Exec(fmt.Sprintf("INSERT INTO movies VALUES (current_timestamp, '%s', '%s', false, -1, '');", movie.AddedBy, movieName))

	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func deleteMovie(db *DB, movie Movie) bool {
	movieName := strings.Replace(movie.Name, "'", "\\'", -1)
	_, err := db.Exec(
		fmt.Sprintf("DELETE FROM movies WHERE name = '%s';",
			movieName))
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}
