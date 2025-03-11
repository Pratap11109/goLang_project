package main

import (
	"fmt"
	"log"
	"net/http"
)

func handelFunc(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err!=nil{
		http.Error(w, "not able parse form", http.StatusBadRequest)
	}

	name := r.Form.Get("name")
	address := r.Form.Get("address")
	fmt.Fprintf(w, " Name: %s\t Address: %s\n", name, address)
	fmt.Fprintf(w, "POST request success full..")
}


func homePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/home"{
		http.Error(w, "404 not found", http.StatusBadRequest)
	}
	if r.Method != "GET"{
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
	}
	_, err := fmt.Fprintf(w, "Home age")
	if err!= nil{
		log.Fatal("Print error in home page")
	}
}


func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/form", handelFunc)
	http.HandleFunc("/home", homePage)
	if err:= http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}

	
}