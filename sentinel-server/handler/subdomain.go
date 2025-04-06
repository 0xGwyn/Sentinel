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
	subdomainsColl := database.GetDBCollection("subdomains")
	httpColl := database.GetDBCollection("http")
	dnsColl := database.GetDBCollection("dns")

	// Delete the requested subdomain
	subFilter := bson.M{"domain": domainName, "name": subdomainName}
	subResult, err := subdomainsColl.DeleteOne(c.Context(), subFilter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if subResult.DeletedCount == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "subdomain not found",
		})
	}

	// Delete related HTTP records
	httpFilter := bson.M{"domain": domainName, "subdomain": subdomainName}
	_, err = httpColl.DeleteMany(c.Context(), httpFilter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Delete related DNS records
	dnsFilter := bson.M{"domain": domainName, "subdomain": subdomainName}
	_, err = dnsColl.DeleteMany(c.Context(), dnsFilter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": subdomainName + " subdomain and related records deleted successfully",
	})
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
	subdomainsColl := database.GetDBCollection("subdomains")
	httpColl := database.GetDBCollection("http")
	dnsColl := database.GetDBCollection("dns")

	// find the requested subdomain
	subdomain := models.Subdomain{}
	filter := bson.M{"domain": domainName, "name": subdomainName}
	projection := bson.M{"_id": 0}
	opts := options.FindOne().SetProjection(projection)
	if err := subdomainsColl.FindOne(c.Context(), filter, opts).Decode(&subdomain); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get the latest HTTP record
	httpRecord := models.HTTP{}
	httpFilter := bson.M{"domain": domainName, "subdomain": subdomainName}
	httpOpts := options.FindOne().
		SetProjection(bson.M{"_id": 0}).
		SetSort(bson.M{"scanning_date": -1})
	_ = httpColl.FindOne(c.Context(), httpFilter, httpOpts).Decode(&httpRecord)

	// Get the latest DNS record
	dnsRecord := models.DNS{}
	dnsFilter := bson.M{"domain": domainName, "subdomain": subdomainName}
	dnsOpts := options.FindOne().
		SetProjection(bson.M{"_id": 0}).
		SetSort(bson.M{"resolution_date": -1})
	_ = dnsColl.FindOne(c.Context(), dnsFilter, dnsOpts).Decode(&dnsRecord)

	// Combine all data in the desired order
	response := SubdomainResponse{
		Subdomain: subdomain,
		DNS:       dnsRecord,
		HTTP:      httpRecord,
	}

	return c.Status(200).JSON(response)
}

type SubdomainResponse struct {
	Subdomain models.Subdomain `json:"subdomain"`
	DNS       models.DNS       `json:"latest_dns"`
	HTTP      models.HTTP      `json:"latest_http"`
}
