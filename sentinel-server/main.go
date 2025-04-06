package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/0xgwyn/sentinel/config"
	"github.com/0xgwyn/sentinel/database"
	"github.com/0xgwyn/sentinel/router"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	// init db
	err := database.InitDB()
	if err != nil {
		return err
	}

	// defer closing db
	defer database.CloseDB()

	// Initialize indexes unless explicitly skipped
	if skip, _ := config.LoadEnv("SKIP_INDEXES"); skip != "true" {
		if err := database.InitIndexes(); err != nil {
			return err
		}
		log.Println("Indexes initialized successfully")
	}

	// Insert mock data if enabled
	if insert, _ := config.LoadEnv("INSERT_MOCK_DATA"); insert == "true" {
		if err := database.InsertMockData(); err != nil {
			return err
		}
		log.Println("Mock data inserted successfully")
	}

	// create app
	app := fiber.New()

	// add basic middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// add routes
	router.AddRouterGroup(app)

	// start server
	var port string
	if port, err = config.LoadEnv("PORT"); port == "" {
		port = "9000"
	}
	app.Listen(":" + port)

	return nil
}
