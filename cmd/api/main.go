package main

import (
	"log"
	"net/http"

	"github.com/v3lichko/student-distribution/internal/db"
	"github.com/v3lichko/student-distribution/internal/handler"
)

func main() {
	database := db.Connect()
	defer database.Close()

	serve := http.NewServeMux()
	studentHandler := handler.NewStudentHandler(database)
	serve.HandleFunc("/students", studentHandler.CreateStudent)
	serve.HandleFunc("/health", handler.HealthHandler)
	log.Println("server is started")
	err := http.ListenAndServe(":8080", serve)
	if err != nil {
		panic(err)
	}

}
