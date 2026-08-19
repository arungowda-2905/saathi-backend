package config

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var DB *mongo.Database

func ConnectDB() error {

	// Local MongoDB
	uri := "mongodb://localhost:27017"

	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return err
	}

	// Check whether MongoDB is reachable
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	fmt.Println("Successfully connected to MongoDB!")

	// Select database
	DB = client.Database("saathi_db")

	return nil
}