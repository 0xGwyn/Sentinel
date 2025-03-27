package handler

import (
	"strings"

	"github.com/0xgwyn/sentinel/database"
	"github.com/0xgwyn/sentinel/models"
	"github.com/dchest/validator"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func UpdateDomain(c *fiber.Ctx) error {
	domainName := c.Params("domainName")
	coll := database.GetDBCollection("domains")

	// Parse the body
	domain := new(models.Domain)
	if err := c.BodyParser(domain); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Update the scope of the domain if it exists
	filter := bson.M{"name": strings.ToLower(domainName)}
	var update bson.M
	if domain.InScope != nil && domain.OutOfScope != nil {
		update = bson.M{
			"$set": bson.M{"in_scope": domain.InScope, "out_of_scope": domain.OutOfScope},
		}
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

func GetDomain(c *fiber.Ctx) error {
	domainName := c.Params("domainName")
	domainsColl := database.GetDBCollection("domains")
	subdomainsColl := database.GetDBCollection("subdomains")

	// find the requested domain
	domain := new(models.Domain)
	domainFilter := bson.M{"name": domainName}
	domainProjection := bson.M{"name": 1, "in_scope": 1, "out_of_scope": 1}
	domainOpts := options.FindOne().SetProjection(domainProjection)
	if err := domainsColl.FindOne(c.Context(), domainFilter, domainOpts).Decode(domain); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// find the subdomains related to the domain
	subdomainFilter := bson.M{"domain": domain.Name}
	subdomainProjection := bson.M{"name": 1}
	subdomainOpts := options.Find().SetProjection(subdomainProjection)
	cursor, err := subdomainsColl.Find(c.Context(), subdomainFilter, subdomainOpts)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	defer cursor.Close(c.Context())

	// iterate over the cursor
	subdomains := make([]string, 0)
	for cursor.Next(c.Context()) {
		subdomain := new(models.Subdomain)
		err := cursor.Decode(subdomain)
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

func GetDomains(c *fiber.Ctx) error {
	coll := database.GetDBCollection("domains")

	// find all domains
	projection := bson.M{"name": 1}
	opts := options.Find().SetProjection(projection)
	cursor, err := coll.Find(c.Context(), bson.M{}, opts)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	defer cursor.Close(c.Context())

	// iterate over the cursor
	domains := make([]string, 0)
	for cursor.Next(c.Context()) {
		domain := new(models.Domain)
		err := cursor.Decode(domain)
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

func CreateDomain(c *fiber.Ctx) error {
	coll := database.GetDBCollection("domains")

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

	// create the domain also save the domain in lowercase
	domain.Name = strings.ToLower(domain.Name)
	_, err := coll.InsertOne(c.Context(), domain)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(domain)
}

func DeleteDomain(c *fiber.Ctx) error {
	domainName := c.Params("domainName")
	domainsColl := database.GetDBCollection("domains")
	subdomainsColl := database.GetDBCollection("subdomains")

	// delete the requested domain if it exists
	domainsFilter := bson.M{"name": domainName}
	result, err := domainsColl.DeleteOne(c.Context(), domainsFilter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if result.DeletedCount == 0 {
		return c.Status(500).JSON(fiber.Map{
			"error": "the domain does not exist",
		})
	}

	// delete all subdomains of that domain
	subdomainsFilter := bson.M{"domain": domainName}
	_, err = subdomainsColl.DeleteMany(c.Context(), subdomainsFilter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"info": domainName + " and its subdomains are removed",
	})
}
