package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:id`
	Isbn     string    `json:isbn`
	Title    string    `json:title`
	Director *Director `json:director`
}

type Director struct {
	Firstname string `json:firstname`
	Lasttname string `json:lastname`
}

var movies []Movie

// ----------- all movies -------------

// w = response
// r = request
func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") //?
	json.NewEncoder(w).Encode(movies)
}

// ----------- delete -------------
func deletemovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") //?
	params := mux.Vars(r)                              // bech t3adi des parametre ( id ) w t7ot l id bil minisxule
	for index, item := range movies {                  // comme si " for each "

		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	// raja3 l movies lkol
	json.NewEncoder(w).Encode(movies)
}

// ----------- movie -------------
func getmovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") //?
	params := mux.Vars(r)
	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

// ----------- create movie -------------

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") //?
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(1000000))
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)
}

// ----------- update movie -------------

func updatemovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") //?
	params := mux.Vars(r)

	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = strconv.Itoa(rand.Intn(1000000))
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}

func main() {

	r := mux.NewRouter()

	movies = append(movies, Movie{ID: "1", Isbn: "123456", Title: "Movie One", Director: &Director{Firstname: "Jk", Lasttname: "roling"}})
	movies = append(movies, Movie{ID: "2", Isbn: "456789", Title: "Movie Two", Director: &Director{Firstname: "Agatha", Lasttname: "christy"}})

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getmovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updatemovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deletemovie).Methods("DELETE")

	fmt.Printf("Starting server at port 3000 \n")
	log.Fatal(http.ListenAndServe(":3000", r))
}
