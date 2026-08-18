package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func ConnectDB() (*mongo.Client, error) {
	// 1. Define your connection string (replace with your local or MongoDB Atlas URI)
	// For Atlas, get the URI from the MongoDB Atlas Dashboard (Connect -> Connect your application)
	uri := "mongodb://localhost:27017"

	// 2. Set up connection options
	clientOptions := options.Client().ApplyURI(uri)

	// 3. Create a context with a timeout for the connection phase
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 4. Connect to MongoDB
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	// 5. Ensure connection is safely closed when main() exits
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalf("Failed to disconnect cleanly: %v", err)
		}
		fmt.Println("Connection to MongoDB closed.")
	}()

	// 6. Ping the database to verify the connection is alive
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("Could not ping database: %v", err)
	}

	fmt.Println("Successfully connected to MongoDB!")

	// 7. Access a database and collection
	database := client.Database("my_database")
	collection := database.Collection("my_collection")

	// Your CRUD operations go here using the 'collection' instance
	_ = collection
	return client, nil
}
