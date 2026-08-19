package config

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var Client *mongo.Client
var DB *mongo.Database

func ConnectDB() error {

	uri := "mongodb://localhost:27017"

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	client, err := mongo.Connect(
		options.Client().ApplyURI(uri),
	)
	if err != nil {
		return err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	// Store the client globally
	Client = client

	// Select database
	DB = client.Database("saathi")

	fmt.Println("Successfully connected to MongoDB!")

	return nil
}

func DisconnectDB() error {

	if Client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	err := Client.Disconnect(ctx)

	if err != nil {
		return err
	}

	fmt.Println("MongoDB connection closed.")

	return nil
}
