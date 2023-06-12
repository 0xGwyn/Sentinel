package router

import (
	"strings"

	"github.com/0xgwyn/sentinel/common"
	"github.com/0xgwyn/sentinel/models"
	"github.com/dchest/validator"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddDomainGroup(app *fiber.App) {
	domainGroup := app.Group("/api/v1/domains")

	domainGroup.Get("/", getDomains)
	domainGroup.Delete("/:domainName", deleteDomain)
	domainGroup.Get("/:domainName", getDomain)
	domainGroup.Post("/", createDomain)
	domainGroup.Patch("/:domainName", updateDomain)
	domainGroup.Get("/:domainName/:subdomainName", getSubdomain)

}

func getSubdomain(c *fiber.Ctx) error {
	

	return nil
}

func updateDomain(c *fiber.Ctx) error {
	domainName := c.Params("domainName")
	coll := common.GetDBCollection("domains")

	// Parse the body
	domain := new(models.Domain)
	if err := c.BodyParser(domain); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Update the scope of the domain if it exists
	filter := bson.M{"name": strings.ToLower(domainName)}
	var update primitive.M
	if domain.InScope != nil && domain.OutOfScope != nil {
		update = bson.M{"$set": bson.M{"in_scope": domain.InScope, "out_of_scope": domain.OutOfScope}}
	} else if domain.InScope == nil && domain.OutOfScope != nil {
		update = bson.M{"$set": bson.M{"out_of_scope": domain.OutOfScope}}
	} else if domain.InScope != nil && domain.OutOfScope == nil {
		update = bson.M{"$set": bson.M{"in_scope": domain.InScope}}
	} else {
		return c.Status(500).JSON(fiber.Map{
			"error": "either in_scope or out_of_scope is needed",
		})
	}

	result, err := coll.UpdateOne(c.Context(), filter, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(result)
}

func getDomain(c *fiber.Ctx) error {
	domainName := c.Params("domainName")
	coll_domains := common.GetDBCollection("domains")
	coll_subdomains := common.GetDBCollection("subdomains")

	// find the requested domain
	domain := new(models.Domain)
	domain_filter := bson.M{"name": domainName}
	domain_projection := bson.M{"name": 1, "in_scope": 1, "out_of_scope": 1}
	domain_opts := options.FindOne().SetProjection(domain_projection)
	if err := coll_domains.FindOne(c.Context(), domain_filter, domain_opts).Decode(&domain); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// find the subdomains related to the domain
	subdomain_filter := bson.M{"domain": domain.Name}
	subdomain_projection := bson.M{"name": 1}
	subdomain_opts := options.Find().SetProjection(subdomain_projection)
	cursor, err := coll_subdomains.Find(c.Context(), subdomain_filter, subdomain_opts)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// iterate over the cursor
	subdomains := make([]string, 0)
	for cursor.Next(c.Context()) {
		subdomain := models.Subdomain{}
		err := cursor.Decode(&subdomain)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		// only get the Name field
		subdomains = append(subdomains, subdomain.Name)
	}

	return c.Status(200).JSON(fiber.Map{
		"domain":     domain,
		"subdomains": subdomains,
	})
}

func getDomains(c *fiber.Ctx) error {
	coll := common.GetDBCollection("domains")

	// find all domains
	projection := bson.M{"name": 1}
	opts := options.Find().SetProjection(projection)
	cursor, err := coll.Find(c.Context(), bson.M{}, opts)
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
	coll := common.GetDBCollection("domains")

	// Parse the body
	domain := new(models.Domain)
	if err := c.BodyParser(domain); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Check if the domain is valid
	if !validator.IsValidDomain(domain.Name) {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid domain",
		})
	}

	// Check if either out of scope or in scope is not set (at least one should be set)
	if len(domain.InScope) == 0 && len(domain.OutOfScope) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "either in_scope or out_of_scope is required",
		})
	}

	// Check if the domain already exists in the collection
	filter := bson.M{"name": strings.ToLower(domain.Name)}
	existingDomain := new(models.Domain)
	if err := coll.FindOne(c.Context(), filter).Decode(existingDomain); err == nil {
		return c.Status(409).JSON(fiber.Map{
			"error": "domain already exists",
		})
	}

	// create the domain
	_, err := coll.InsertOne(c.Context(), domain)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(domain)
}

func deleteDomain(c *fiber.Ctx) error {
	domainName := c.Params("domainName")
	coll := common.GetDBCollection("domains")

	// delete the requested domain
	result, err := coll.DeleteOne(c.Context(), bson.M{"name": domainName})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(result)
}
