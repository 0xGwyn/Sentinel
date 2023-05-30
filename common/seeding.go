package common

import (
	"context"
	"log"
	"time"

	"github.com/0xgwyn/sentinel/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Seeding() error {
	coll := GetDBCollection("companies")

	// Insert some seed data
	company1 := models.Company{
		ID:   primitive.NewObjectID(),
		Name: "Google",
		Domains: []models.Domain{
			{
				ID:   primitive.NewObjectID(),
				Name: "google.com",
				Subdomains: []models.Subdomain{
					{
						ID:         primitive.NewObjectID(),
						Name:       "www.google.com",
						StatusCode: "200",
						Title:      "Google",
						CDN:        "Google Cloud",
						Technology: "AngularJS",
						CreatedAt:  primitive.NewDateTimeFromTime(time.Now()),
						UpdatedAt:  primitive.NewDateTimeFromTime(time.Now()),
						Paths:      []string{"/search", "/images", "/maps"},
					},
					{
						ID:    primitive.NewObjectID(),
						Name:  "api.google.com",
						Paths: []string{"/inbox", "/sent", "/drafts"},
					},
				},
			},
		},
	}

	company2 := models.Company{
		ID:   primitive.NewObjectID(),
		Name: "Apple",
		Domains: []models.Domain{
			{
				ID:   primitive.NewObjectID(),
				Name: "apple.com",
				Subdomains: []models.Subdomain{
					{
						ID:         primitive.NewObjectID(),
						Name:       "www.apple.com",
						StatusCode: "200",
						Title:      "Apple",
						CDN:        "Akamai",
						Technology: "React",
						CreatedAt:  primitive.NewDateTimeFromTime(time.Now()),
						UpdatedAt:  primitive.NewDateTimeFromTime(time.Now()),
						Paths:      []string{"/mac", "/iphone", "/ipad"},
					},
				},
			},
		},
	}

	company3 := models.Company{
		ID:   primitive.NewObjectID(),
		Name: "Microsoft",
		Domains: []models.Domain{
			{
				ID:   primitive.NewObjectID(),
				Name: "microsoft.com",
				Subdomains: []models.Subdomain{
					{
						ID:         primitive.NewObjectID(),
						StatusCode: "200",
						Title:      "Microsoft",
						CDN:        "Akamai",
						Technology: "React",
						CreatedAt:  primitive.NewDateTimeFromTime(time.Now()),
						UpdatedAt:  primitive.NewDateTimeFromTime(time.Now()),
						Paths:      []string{"/windows", "/office", "/surface"},
					},
					{
						ID:    primitive.NewObjectID(),
						Name:  "azure.microsoft.com",
						Paths: []string{"/pricing", "/docs", "/support"},
					},
				},
			},
			{
				ID:   primitive.NewObjectID(),
				Name: "github.com/microsoft",
				Subdomains: []models.Subdomain{
					{
						ID:         primitive.NewObjectID(),
						Name:       "github.com/microsoft/vscode",
						StatusCode: "200",
						Title:      "Visual Studio Code",
						Technology: "Electron",
						CreatedAt:  primitive.NewDateTimeFromTime(time.Now()),
						UpdatedAt:  primitive.NewDateTimeFromTime(time.Now()),
						Paths:      []string{"/docs", "/extensions", "/settings"},
					},
					{
						ID:    primitive.NewObjectID(),
						Paths: []string{"/docs", "/samples", "/community"},
					},
				},
			},
		},
	}

	_, err := coll.InsertMany(context.Background(), []interface{}{company1, company2, company3})
	if err != nil {
		return err
	}

	log.Println("Seed data inserted successfully!")

	return nil
}
