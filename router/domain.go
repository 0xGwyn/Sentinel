package router

import (
	"context"

	"github.com/0xgwyn/sentinel/common"
	"github.com/0xgwyn/sentinel/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddDomainGroup(app *fiber.App) {
	domainGroup := app.Group("/api/v1/domains")

	domainGroup.Get("/", getDomains)
	domainGroup.Get("/:domainName", getDomain)

}

func getDomain(c *fiber.Ctx) error {
	domainName := c.Params("domainName")
	coll := common.GetDBCollection("domains")

	// find the requested domain
	var domain models.Domain
	filter := bson.M{"name": domainName}
	projection := bson.M{"name": 1, "subdomains.name": 1}
	opts := options.FindOne().SetProjection(projection)

	if err := coll.FindOne(context.Background(), filter, opts).Decode(&domain); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	subdomains := make([]string, 0)
	for _, subdomain := range domain.Subdomains {
		subdomains = append(subdomains, subdomain.Name)
	}

	return c.Status(200).JSON(fiber.Map{
		"domain":     domain.Name,
		"subdomains": subdomains,
	})
}

func getDomains(c *fiber.Ctx) error {
	coll := common.GetDBCollection("domains")

	// find all domains
	projection := bson.M{"name": 1}
	opts := options.Find().SetProjection(projection)
	cursor, err := coll.Find(context.Background(), bson.M{}, opts)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// iterate over the cursor
	domains := make([]string, 0)
	for cursor.Next(c.Context()) {
		domain := models.Domain{}
		err := cursor.Decode(&domain)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		// only get the Name field
		domains = append(domains, domain.Name)
	}

	return c.Status(200).JSON(fiber.Map{
		"domains": domains,
	})
}

func createDomain(c *fiber.Ctx) error {
	/*coll := common.GetDBCollection("companies")

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
	*/
	return nil
}

func updateDomain(c *fiber.Ctx) error {
	return nil
}

func deleteDomain(c *fiber.Ctx) error {
	return nil
}
