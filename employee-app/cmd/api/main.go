// @title Employee API
// @version 1.0
// @description Simple Employee System using Go and PostgreSQL.
// @host localhost:8080
// @BasePath /
package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "employee-app/docs"
	"employee-app/internal/config"
	"employee-app/internal/db"
	httpHandlers "employee-app/internal/http/handlers"
	"employee-app/internal/http/middleware"
	"employee-app/internal/repository"
)

func main() {
	// 1) Load .env (for local dev)
	if err := godotenv.Load(); err != nil {
		log.Println("  No .env file found, using OS environment variables")
	}

	//  Load config from environment variables
	cfg := config.LoadConfig()

	//  Connect to DB using config
	db.ConnectDB(cfg)
	defer db.DB.Close()

	//  Setup repository and handlers
	repo := repository.NewPostgresEmployeeRepository(db.DB)

	authH := &httpHandlers.AuthHandler{Repo: repo}
	empH := httpHandlers.NewEmployeeHandler(repo)

	//  Setup routes
	mux := http.NewServeMux()

	// Auth routes
	mux.HandleFunc("/register", authH.Register) 
	mux.HandleFunc("/login", authH.Login)

	// Employee routes
	mux.HandleFunc("/employees/", empH.GetByID)

	// Swagger UI
	mux.Handle("/swagger/", httpSwagger.Handler(
    httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // URL to swagger.json
))

	//  Start HTTP server with logging middleware
	addr := ":8080"
	log.Println("Server running on", addr)

	if err := http.ListenAndServe(addr, middleware.Logging(mux)); err != nil {
		log.Fatal(err)
	}
}
