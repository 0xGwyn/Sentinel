package handler

import (
	"strings"
	"time"

	"github.com/0xgwyn/sentinel/database"
	"github.com/0xgwyn/sentinel/models"
	"github.com/gofiber/fiber/v2"
	sliceutil "github.com/projectdiscovery/utils/slice"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func DeleteSubdomain(c *fiber.Ctx) error {
	domainName := c.Params("domainName")
	subdomainName := c.Params("subdomainName")
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
	domainName := c.Params("domainName")
	coll := database.GetDBCollection("subdomains")

	// Parse the body
	newSubdomains := new([]string)
	if err := c.BodyParser(newSubdomains); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get all subdomains related to the given domain from the collection
	filter := bson.M{"domain": domainName}
	projection := bson.M{"name": 1, "_id": 0}
	opts := options.Find().SetProjection(projection)
	cursor, err := coll.Find(c.Context(), filter, opts)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	defer cursor.Close(c.Context())

	// Get the name of all subdomains inside the collection
	oldSubdomains := make([]string, 0)
	for cursor.Next(c.Context()) {
		subdomain := new(models.Subdomain)
		err := cursor.Decode(subdomain)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		// only get the Name field
		oldSubdomains = append(oldSubdomains, subdomain.Name)
	}

	// Trim spaces from all subdomains
	for i := range *newSubdomains {
		(*newSubdomains)[i] = strings.TrimSpace((*newSubdomains)[i])
	}

	// Make the new subdomains unique then compare old ones(in DB) with
	// the new ones(in the request) and if no subdomain is new, let the
	// user know
	newUniqueSubdomains := sliceutil.Dedupe(*newSubdomains)
	_, subsToBeAdded := sliceutil.Diff(oldSubdomains, newUniqueSubdomains)
	if subsToBeAdded == nil {
		return c.Status(200).JSON(fiber.Map{
			"info": "no new subdomain was found",
		})
	}

	// Create new subdomains then insert them to mongodb
	for _, name := range subsToBeAdded {
		subdomain := models.Subdomain{
			Domain:    strings.ToLower(domainName),
			Name:      strings.ToLower(name),
			CreatedAt: bson.NewDateTimeFromTime(time.Now()),
			UpdatedAt: bson.NewDateTimeFromTime(time.Now()),
		}
		_, err := coll.InsertOne(c.Context(), subdomain)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	return c.Status(200).JSON(subsToBeAdded)
}

func GetSubdomain(c *fiber.Ctx) error {
	domainName := c.Params("domainName")
	subdomainName := c.Params("subdomainName")
	coll := database.GetDBCollection("subdomains")

	// find the requested subdomain
	subdomain := new(models.Subdomain)
	filter := bson.M{"domain": domainName, "name": subdomainName}
	projection := bson.M{"_id": 0}
	opts := options.FindOne().SetProjection(projection)
	if err := coll.FindOne(c.Context(), filter, opts).Decode(subdomain); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(subdomain)
}
