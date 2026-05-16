package main

import (
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/v3lichko/student-distribution/docs"
	"github.com/v3lichko/student-distribution/internal/db"
	"github.com/v3lichko/student-distribution/internal/handler"
	"github.com/v3lichko/student-distribution/internal/storage"
)

// @title Student Distribution API
// @version 1.0
// @host localhost:8080
// @BasePath /
func main() {
	database := db.Connect()
	defer database.Close()

	serve := http.NewServeMux()

	studentStorage := storage.NewStudentStorage(database)
	studentHandler := handler.NewStudentHandler(studentStorage)

	groupStorage := storage.NewGroupStorage(database)
	groupHandler := handler.NewGroupHandler(groupStorage)

	distributionStorage := storage.NewDistributionStorage(database)
	distributionHandler := handler.NewDistributionHandler(distributionStorage)

	serve.HandleFunc("/students", studentHandler.Students)
	serve.HandleFunc("/students/import", studentHandler.ImportStudentsCSV)

	serve.HandleFunc("/health", handler.HealthHandler)

	serve.HandleFunc("/groups", groupHandler.Groups)

	serve.HandleFunc("/distribution/run", distributionHandler.StartDistribution)
	serve.HandleFunc("/distribution", distributionHandler.Distribution)
	serve.HandleFunc("/distribution/export", distributionHandler.ExportDistributionCSV)

	serve.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	log.Println("server is started")

	err := http.ListenAndServe(":8080", serve)
	if err != nil {
		panic(err)
	}
}
