package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/0xgwyn/sentinel/handler"
)

func AddRouterGroup(app *fiber.App) {
	routerGroup := app.Group("/api/domains")

	// domain routes
	routerGroup.Get("/", handler.GetDomains)
	routerGroup.Delete("/:domainName", handler.DeleteDomain)
	routerGroup.Get("/:domainName", handler.GetDomain)
	routerGroup.Post("/", handler.CreateDomain)
	routerGroup.Patch("/:domainName", handler.UpdateDomain)

	// subdomain routes
	routerGroup.Get("/:domainName/:subdomainName", handler.GetSubdomain)
	routerGroup.Post("/:domainName", handler.AddSubdomains)
	routerGroup.Delete("/:domainName/:subdomainName", handler.DeleteSubdomain)
}
