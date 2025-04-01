package handler

import (
	"strings"
	"time"

	"github.com/0xgwyn/sentinel/database"
	"github.com/0xgwyn/sentinel/models"
	"github.com/dchest/validator"
	"github.com/gofiber/fiber/v2"
	sliceutil "github.com/projectdiscovery/utils/slice"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func DeleteSubdomain(c *fiber.Ctx) error {
	domainName := strings.ToLower(c.Params("domainName"))
	subdomainName := strings.ToLower(c.Params("subdomainName"))
	coll := database.GetDBCollection("subdomains")

	// delete the requested subdomain
	filter := bson.M{"domain": domainName, "name": subdomainName}
	result, err := coll.DeleteOne(c.Context(), filter)
	if err != nil {
		return c.Status(200).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(result)
}

func AddSubdomains(c *fiber.Ctx) error {
	domainName := strings.ToLower(c.Params("domainName"))
	coll := database.GetDBCollection("subdomains")

	// Parse the body
	newSubdomains := []string{}
	if err := c.BodyParser(&newSubdomains); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Check if the subdomains are valid
	for _, subdomain := range newSubdomains {
		if !validator.IsValidDomain(subdomain) {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid subdomain: " + subdomain,
			})
		}
	}

	// Remove duplicates from the new subdomains
	newUniqueSubdomains := sliceutil.Dedupe(newSubdomains)

	// Prepare a list of subdomains to be added
	subsToBeAdded := []models.Subdomain{}
	for _, name := range newUniqueSubdomains {
		// Check if the subdomain already exists in the database
		filter := bson.M{"domain": domainName, "name": name}
		count, err := coll.CountDocuments(c.Context(), filter)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// If the subdomain does not exist, add it to the list
		if count == 0 {
			subdomain := models.Subdomain{
				Domain:    strings.ToLower(domainName),
				Name:      strings.ToLower(name),
				CreatedAt: bson.NewDateTimeFromTime(time.Now()),
				UpdatedAt: bson.NewDateTimeFromTime(time.Now()),
			}
			subsToBeAdded = append(subsToBeAdded, subdomain)
		}
	}

	// If no new subdomains are found, return a message
	if len(subsToBeAdded) == 0 {
		return c.Status(200).JSON(fiber.Map{
			"info": "no new subdomain was found",
		})
	}

	// Insert the new subdomains into the database
	_, err := coll.InsertMany(c.Context(), subsToBeAdded)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return the newly added subdomains
	return c.Status(200).JSON(subsToBeAdded)
}

func GetSubdomain(c *fiber.Ctx) error {
	domainName := strings.ToLower(c.Params("domainName"))
	subdomainName := strings.ToLower(c.Params("subdomainName"))
	coll := database.GetDBCollection("subdomains")

	// find the requested subdomain
	subdomain := models.Subdomain{}
	filter := bson.M{"domain": domainName, "name": subdomainName}
	projection := bson.M{"_id": 0}
	opts := options.FindOne().SetProjection(projection)
	if err := coll.FindOne(c.Context(), filter, opts).Decode(&subdomain); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(subdomain)
}
