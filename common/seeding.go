package common

import (
	"context"
	"log"
	"time"

	"github.com/0xgwyn/sentinel/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Seeding() error {
	coll := GetDBCollection("domains")

	// Insert some seed data
	domain1 := models.Domain{
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
				CreatedAt:  primitive.NewDateTimeFromTime(time.Now().Add(-24 * time.Hour)),
				UpdatedAt:  primitive.NewDateTimeFromTime(time.Now()),
				Paths:      []string{"/mac", "/iphone", "/ipad"},
			},
		},
	}

	domain2 := models.Domain{
		ID:   primitive.NewObjectID(),
		Name: "example.com",
		Subdomains: []models.Subdomain{
			{
				ID:         primitive.NewObjectID(),
				Name:       "www",
				StatusCode: "200 OK",
				Title:      "Example Domain",
				CDN:        "Cloudflare",
				Technology: "Nginx",
				CreatedAt:  primitive.NewDateTimeFromTime(time.Now().Add(-24 * time.Hour)),
				UpdatedAt:  primitive.NewDateTimeFromTime(time.Now()),
				IPs:        []string{"192.0.2.1", "192.0.2.2"},
				Providers:  []string{"Google Cloud", "DigitalOcean"},
				Paths:      []string{"/", "/about", "/contact"},
			},
		},
	}

	domain3 := models.Domain{
		ID:   primitive.NewObjectID(),
		Name: "google.com",
		Subdomains: []models.Subdomain{
			{
				ID:         primitive.NewObjectID(),
				Name:       "www",
				StatusCode: "200 OK",
				Title:      "Google",
				CDN:        "Google Cloud CDN",
				Technology: "Google Web Server",
				CreatedAt:  primitive.NewDateTimeFromTime(time.Now().Add(-48 * time.Hour)),
				UpdatedAt:  primitive.NewDateTimeFromTime(time.Now()),
				IPs:        []string{"172.217.6.68", "172.217.6.67"},
				Providers:  []string{"Google Cloud", "Akamai"},
				Paths:      []string{"/", "/about", "/search"},
			},
			{
				ID:         primitive.NewObjectID(),
				Name:       "mail",
				StatusCode: "200 OK",
				Title:      "Gmail",
				CDN:        "",
				Technology: "Google Web Server",
				CreatedAt:  primitive.NewDateTimeFromTime(time.Now().Add(-48 * time.Hour)),
				UpdatedAt:  primitive.NewDateTimeFromTime(time.Now()),
				Providers:  []string{"Google Cloud", "Akamai"},
				Paths:      []string{"/", "/inbox", "/settings"},
			},
		},
	}

	_, err := coll.InsertMany(context.Background(), []interface{}{domain1, domain2, domain3})
	if err != nil {
		return err
	}

	log.Println("Seed data inserted successfully!")

	return nil
}
