package main

import (
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/v3lichko/student-distribution/docs"
	"github.com/v3lichko/student-distribution/internal/config"
	"github.com/v3lichko/student-distribution/internal/db"
	"github.com/v3lichko/student-distribution/internal/handler"
	"github.com/v3lichko/student-distribution/internal/storage"
)

// @title Student Distribution API
// @version 1.0
// @host localhost:8080
// @BasePath /
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	database, err := db.Connect(cfg.DB)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer database.Close()

	mux := http.NewServeMux()

	studentStorage := storage.NewStudentStorage(database)
	groupStorage := storage.NewGroupStorage(database)
	distributionStorage := storage.NewDistributionStorage(database)

	studentHandler := handler.NewStudentHandler(studentStorage)
	groupHandler := handler.NewGroupHandler(groupStorage)
	distributionHandler := handler.NewDistributionHandler(distributionStorage)

	mux.HandleFunc("GET /students", studentHandler.GetStudents)
	mux.HandleFunc("POST /students", studentHandler.CreateStudent)
	mux.HandleFunc("DELETE /students", studentHandler.DeleteStudent)
	mux.HandleFunc("POST /students/import", studentHandler.ImportStudentsCSV)

	mux.HandleFunc("GET /health", handler.HealthHandler)

	mux.HandleFunc("GET /groups", groupHandler.GetGroups)
	mux.HandleFunc("POST /groups", groupHandler.CreateGroup)

	mux.HandleFunc("POST /distribution/run", distributionHandler.StartDistribution)
	mux.HandleFunc("GET /distribution", distributionHandler.GetDistribution)
	mux.HandleFunc("GET /distribution/export", distributionHandler.ExportDistributionCSV)

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	log.Printf("listening on %s", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, mux); err != nil {
		log.Fatalf("server: %v", err)
	}
}
