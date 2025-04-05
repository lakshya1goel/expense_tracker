package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/lakshya1goel/expense_tracker/database"
)


func main() {

	database.ConnectDb()

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Hello, World!")
	})

	log.Fatal(http.ListenAndServe(":8000", nil))
}
