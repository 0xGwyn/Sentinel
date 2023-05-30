package router

import (
	"context"

	"github.com/0xgwyn/sentinel/common"
	"github.com/0xgwyn/sentinel/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func AddCompanyGroup(app *fiber.App) {
	companyGroup := app.Group("/api/v1/companies")

	companyGroup.Get("/", getCompanies)
	companyGroup.Get("/:companyName", getCompany)
	companyGroup.Get("/:companyName/:domainName", getDomain)
	/*companyGroup.Post("/", createCompany)
	companyGroup.Patch("/:name", updateCompany)
	companyGroup.Delete("/:name", deleteCompany)*/
}

func getDomain(c *fiber.Ctx) error {
	domainName := c.Params("domainName")
	companyName := c.Params("companyName")
	coll := common.GetDBCollection("companies")

	// find the requested domain
	var domain models.Domain
	filter := bson.M{"name": bson.M{"$regex": companyName, "$options": "i"}, "domains.name": bson.M{"$regex": domainName, "$options": "i"}}

	if err := coll.FindOne(context.Background(), filter).Decode(&domain); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(domain)
}

func getCompanies(c *fiber.Ctx) error {
	coll := common.GetDBCollection("companies")

	// find all companies
	companies := make([]models.Company, 0)
	cursor, err := coll.Find(c.Context(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// iterate over the cursor
	for cursor.Next(c.Context()) {
		company := models.Company{}
		err := cursor.Decode(&company)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		companies = append(companies, company)
	}

	return c.Status(200).JSON(companies)
}

func getCompany(c *fiber.Ctx) error {
	companyName := c.Params("companyName")
	coll := common.GetDBCollection("companies")

	// find the requested company
	var company models.Company
	filter := bson.M{"name": bson.M{"$regex": companyName, "$options": "i"}}

	if err := coll.FindOne(c.Context(), filter).Decode(&company); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(company)
}

func createCompany(c *fiber.Ctx) error {
	coll := common.GetDBCollection("companies")

	company := new(models.Company)
	if err := c.BodyParser(company); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Check if the user already exists in the collection
	filter := bson.M{"name": company.Name}
	var existingCompany models.Company
	if err := coll.FindOne(c.Context(), filter).Decode(&existingCompany); err != nil {
		return c.Status(409).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// create a company

	return nil
}

func updateCompany(c *fiber.Ctx) error {
	return nil
}

func deleteCompany(c *fiber.Ctx) error {
	return nil
}
