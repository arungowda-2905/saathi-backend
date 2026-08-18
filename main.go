package main

import (
	"fmt"
	"saathi-backend/config"
)

func main() {
	// Call the function to connect to MongoDB
	config.ConnectDB()
	fmt.Println("MongoDB connection established successfully.")
}
