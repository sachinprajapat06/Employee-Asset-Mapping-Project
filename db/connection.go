package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Database *mongo.Database // Global variable to hold the database instance

// InitializeDatabase connects to MongoDB and assigns the Database instance.
func InitializeDatabase(uri string, dbName string) error {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	err = client.Connect(ctx)
	if err != nil {
		return err
	}

	// Ping the MongoDB server
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	// Assign the connected database to the global variable
	Database = client.Database(dbName)
	log.Println("Successfully connected to MongoDB!")
	return nil
}
