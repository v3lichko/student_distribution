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
	serve.HandleFunc("/students", studentHandler.Students)
	groupHandler := handler.NewGroupHandler(database)
	serve.HandleFunc("/health", handler.HealthHandler)
	serve.HandleFunc("/groups", groupHandler.Groups)
	distributionHandler := handler.NewDistributionHandler(database)
	serve.HandleFunc("/distribution/run", distributionHandler.StartDistribution)
	serve.HandleFunc("/distribution", distributionHandler.Distribution)
	serve.HandleFunc("/distribution/export", distributionHandler.ExportDistributionCSV)
	serve.HandleFunc("/students/import", studentHandler.ImportStudentsCSV)
	log.Println("server is started")
	err := http.ListenAndServe(":8080", serve)
	if err != nil {
		panic(err)
	}

}
