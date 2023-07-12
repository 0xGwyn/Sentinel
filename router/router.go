package router

import (
	"strings"
	"time"

	"github.com/dchest/validator"
	"github.com/gofiber/fiber/v2"
	sliceutil "github.com/projectdiscovery/utils/slice"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/0xgwyn/sentinel/common"
	"github.com/0xgwyn/sentinel/models"
)

func AddRouterGroup(app *fiber.App) {
	routerGroup := app.Group("/api/v1/domains")

	// domain routes
	routerGroup.Get("/", getDomains)
	routerGroup.Delete("/:domainName", deleteDomain)
	routerGroup.Get("/:domainName", getDomain)
	routerGroup.Post("/", createDomain)
	routerGroup.Patch("/:domainName", updateDomain)

	// subdomain routes
	routerGroup.Get("/:domainName/:subdomainName", getSubdomain)
	routerGroup.Post("/:domainName", addSubdomains)
	routerGroup.Delete("/:domainName/:subdomainName", deleteSubdomain)
}

func deleteSubdomain(c *fiber.Ctx) error {
	domainName := c.Params("domainName")
	subdomainName := c.Params("subdomainName")
	coll := common.GetDBCollection("subdomains")

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

func addSubdomains(c *fiber.Ctx) error {
	domainName := c.Params("domainName")
	coll := common.GetDBCollection("subdomains")

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
			CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
			UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
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

func getSubdomain(c *fiber.Ctx) error {
	domainName := c.Params("domainName")
	subdomainName := c.Params("subdomainName")
	coll := common.GetDBCollection("subdomains")

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

func getDomain(c *fiber.Ctx) error {
	domainName := c.Params("domainName")
	domainsColl := common.GetDBCollection("domains")
	subdomainsColl := common.GetDBCollection("subdomains")

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

func deleteDomain(c *fiber.Ctx) error {
	domainName := c.Params("domainName")
	domainsColl := common.GetDBCollection("domains")
	subdomainsColl := common.GetDBCollection("subdomains")

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
