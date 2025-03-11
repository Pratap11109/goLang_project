package main

import (
	"encoding/json"
	"strconv"
	"net/http"
	"github.com/julienschmidt/httprouter"
)

type Movie struct{
	ID string `json:"id"`
	Name string `json:"name"`
	Rating float32 `json:"rating"`
	Director *Director "json:director"
}

type Director struct{
	FirstName string `json:"firstname"`
	LastName string `json:"lastname"`
}

var movieSlice []Movie

func __init__(){
	movieSlice = make([]Movie, 0)
}

func getMovies(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	w.Header().Set("content-Type", "application/json")
	json.NewEncoder(w).Encode(movieSlice)
}

func getMovie(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	id := ps.ByName("id")
	for _, movie := range movieSlice{
		if movie.ID == id{
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
	http.Error(w, "204 Entry not Found", http.StatusNoContent)
}

func createMovie(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	w.Header().Set("content-Type", "application/json")
	var movie Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil{
		http.Error(w, "406 not acceptable", http.StatusNotAcceptable)
	}
	movie.ID = strconv.Itoa(len(movieSlice)+1) // updating movie so we can update accordingly our database
	movieSlice = append(movieSlice, movie)
	json.NewEncoder(w).Encode(movieSlice)
}

func updateMovie(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	id := ps.ByName("id")
	for idx, movie := range movieSlice{
		if movie.ID == id{
			movieSlice = append(movieSlice[:idx], movieSlice[idx+1:]...)
			break
		}
	}
	var movie Movie
	json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = id
	movieSlice = append(movieSlice, movie)
	json.NewEncoder(w).Encode(movieSlice)
}

func deleteMovie(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	if r.Method != "DELETE"{
		http.Error(w, "404 Not supported method", http.StatusNotFound)
	}
	id := ps.ByName("id")
	w.Header().Set("content-Type", "application/json")
	for idx, movie := range movieSlice{
		if movie.ID == id{
			movieSlice = append(movieSlice[:idx], movieSlice[idx+1:]...)
			json.NewEncoder(w).Encode(movie)		
			return 	
		}
	}
	http.Error(w, "404 Not found", http.StatusNotFound)
}

func main() {
	router := httprouter.New()
	movieSlice = append(movieSlice, Movie{ID: "1", Name: "Movie one", Rating: 4.5, Director: &Director{ FirstName: "Pratap", LastName: "Thakur"} })
	movieSlice = append(movieSlice, Movie{ID: "2", Name: "Movie Two", Rating: 4.0, Director: &Director{ FirstName: "Steve", LastName: "Job"} })
	router.GET("/movies", getMovies) // reading all movies
	router.GET("/movie/:id", getMovie) // fetch single movie based on id
	router.POST("/movie", createMovie) // create new movie
	router.PUT("/movie/:id", updateMovie) // update specific movie based on id
	router.DELETE("/movie/:id", deleteMovie) // delete movie based on id
	http.ListenAndServe(":8081", router)
}