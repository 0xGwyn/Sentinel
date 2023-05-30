package main

import (
	"fmt"
	"os"

	"github.com/0xgwyn/sentinel/common"
	"github.com/0xgwyn/sentinel/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	err := run()

	if err != nil {
		panic(err)
	}
}

func run() error {
	// init env
	err := common.LoadEnv()
	if err != nil {
		return err
	}

	// init db
	err = common.InitDB()
	if err != nil {
		return err
	}

	// seeding phase
	/*err = common.Seeding()
	if err != nil {
		return err
	}*/

	// defer closing db
	defer common.CloseDB()

	// create app
	app := fiber.New()

	// add basic middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// add routes
	router.AddCompanyGroup(app)

	// start server
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "9000"
	}
	app.Listen(":" + port)
	fmt.Println("testing")

	return nil
}
