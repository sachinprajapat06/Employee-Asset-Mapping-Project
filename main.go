package main

import (
	"employee-asset-system/db"
	"employee-asset-system/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	routes.RegisterRoutes(r)

	// MongoDB URI and database name
	mongoURI := "mongodb://localhost:27017"
	databaseName := "db"

	// Initialize the database connection
	err := db.InitializeDatabase(mongoURI, databaseName)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
