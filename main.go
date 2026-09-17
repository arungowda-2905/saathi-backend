// package main

// import (
// 	"log"
// 	"github.com/joho/godotenv"

// 	"saathi-backend/config"
// 	"saathi-backend/handlers"
// 	"saathi-backend/routes"
// 	tutorialservice "saathi-backend/tutorial-service"

// 	"github.com/gofiber/fiber/v2"
// )

// // func main() {

// // 	err := godotenv.Load("dev.env")
// // 	if err != nil {
// // 		log.Fatal("Error loading dev.env")
// // 	}

// // 	if err := config.ConnectDB(); err != nil {
// // 		log.Fatal("MongoDB connection failed:", err)
// // 	}
// // 	defer config.DisconnectDB()

// // 	ctx := context.Background()

// // 	// Create GCS service
// // 	gcsService, err := gcs.NewService(ctx)
// // 	if err != nil {
// // 		log.Fatalf("Failed to initialize GCS service: %v", err)
// // 	}
// // 	defer func() {
// // 		if err := gcsService.Close(); err != nil {
// // 			log.Printf("Failed to close GCS service: %v", err)
// // 		}
// // 	}()

// // 	app := fiber.New(fiber.Config{
// // 		BodyLimit: 500 * 1024 * 1024,
// // 	})

// // 	app.Use(cors.New(cors.Config{
// // 		AllowOrigins: "http://localhost:3000,http://127.0.0.1:3000",
// // 		AllowHeaders: "Origin, Content-Type, Accept, X-Role, Authorization",
// // 		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
// // 	}))

// // 	tutorialHandler := handlers.NewTutorialHandler(tutorialService)
// // 	routes.SetupRoutes(app, tutorialHandler)

// // 	log.Println("Server running on http://localhost:8080")

// // 	if err := app.Listen(":8080"); err != nil {
// // 		log.Fatal(err)
// // 	}
// // }

// func main() {
// 	err := godotenv.Load("dev.env")
// 	if err != nil {
// 		log.Println("Warning: .env file not found")
// 	}

// 	if err := config.ConnectDB(); err != nil {
// 		log.Fatal("MongoDB connection failed:", err)
// 	}
// 	defer func() {
// 		if err := config.DisconnectDB(); err != nil {
// 			log.Printf("MongoDB disconnect failed: %v", err)
// 		}
// 	}()

// 	tutorialService := tutorialservice.NewTutorialService("")
// 	tutorialHandler := handlers.NewTutorialHandler(tutorialService)

// 	app := fiber.New()
// 	routes.SetupRoutes(app, tutorialHandler)

// 	log.Println("Server running on http://localhost:8080")
// 	if err := app.Listen(":8080"); err != nil {
// 		log.Fatal(err)
// 	}
// }

package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"saathi-backend/config"
	"saathi-backend/handlers"
	"saathi-backend/routes"
	tutorialservice "saathi-backend/tutorial-service"
)

func main() {

	// Load environment variables from dev.env
	err := godotenv.Load("dev.env")
	if err != nil {
		log.Println("Warning: dev.env file not found:", err)
	}

	// Connect to MongoDB
	if err := config.ConnectDB(); err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}

	defer func() {
		if err := config.DisconnectDB(); err != nil {
			log.Printf("MongoDB disconnect failed: %v", err)
		}
	}()

	// Get GCS video prefix from environment
	prefix := os.Getenv("GCS_VIDEO_PREFIX")

	if prefix == "" {
		log.Fatal("GCS_VIDEO_PREFIX is not configured")
	}

	log.Printf("GCS video prefix: %s", prefix)

	// Create tutorial service
	tutorialService := tutorialservice.NewTutorialService(prefix)

	// Create tutorial handler
	tutorialHandler := handlers.NewTutorialHandler(tutorialService)

	// Create Fiber app
	app := fiber.New()

	// Setup routes
	routes.SetupRoutes(app, tutorialHandler)

	log.Println("Server running on http://localhost:8080")

	// Start server
	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
