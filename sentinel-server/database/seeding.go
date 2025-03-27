package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/0xgwyn/sentinel/models"
)

func Seeding() error {

	// Insert some seed data
	domain1 := models.Domain{
		Name:    "apple.com",
		InScope: []string{"*.apple.com"},
	}
	sub1_1 := models.Subdomain{
		// ID:         primitive.NewObjectID(),
		Domain:     "apple.com",
		Name:       "www.apple.com",
		StatusCode: 200,
		Title:      "Apple",
		CDN:        "Akamai",
		Technology: "React",
		CreatedAt:  bson.NewDateTimeFromTime(time.Now().Add(-24 * time.Hour)),
		UpdatedAt:  bson.NewDateTimeFromTime(time.Now()),
	}

	domain2 := models.Domain{
		Name:    "example.com",
		InScope: []string{"*.example.com"},
	}
	sub2_1 := models.Subdomain{
		// ID:         primitive.NewObjectID(),
		Domain:     "example.com",
		Name:       "www.example.com",
		StatusCode: 200,
		Title:      "Example Domain",
		CDN:        "Cloudflare",
		Technology: "Nginx",
		CreatedAt:  bson.NewDateTimeFromTime(time.Now().Add(-24 * time.Hour)),
		UpdatedAt:  bson.NewDateTimeFromTime(time.Now()),
		ARecords:   []string{"192.0.2.1", "192.0.2.2"},
		Providers:  []string{"crtsh", "subfinder"},
	}

	domain3 := models.Domain{
		Name:       "google.com",
		OutOfScope: []string{"*.test.google.com"},
	}
	sub3_1 := models.Subdomain{
		// ID:         primitive.NewObjectID(),
		Domain:     "google.com",
		Name:       "www",
		StatusCode: 200,
		Title:      "Google",
		CDN:        "Google Cloud CDN",
		Technology: "Google Web Server",
		CreatedAt:  bson.NewDateTimeFromTime(time.Now().Add(-48 * time.Hour)),
		UpdatedAt:  bson.NewDateTimeFromTime(time.Now()),

		ARecords:  []string{"172.217.6.68", "172.217.6.67"},
		Providers: []string{"abuseipdb", "subfinder"},
	}
	sub3_2 := models.Subdomain{
		// ID:         primitive.NewObjectID(),
		Domain:     "google.com",
		Name:       "mail",
		StatusCode: 200,
		Title:      "Gmail",
		CDN:        "",
		Technology: "Google Web Server",
		CreatedAt:  bson.NewDateTimeFromTime(time.Now().Add(-48 * time.Hour)),
		UpdatedAt:  bson.NewDateTimeFromTime(time.Now()),
		Providers:  []string{"subfinder"},
	}

	// insert domains into DB
	coll_domains := GetDBCollection("domains")
	_, err := coll_domains.InsertMany(context.Background(), []interface{}{domain1, domain2, domain3})
	if err != nil {
		return err
	}

	// insert subdomains into DB
	coll_subdomains := GetDBCollection("subdomains")
	_, err = coll_subdomains.InsertMany(context.Background(), []interface{}{sub1_1, sub2_1, sub3_1, sub3_2})
	if err != nil {
		return err
	}

	log.Println("Seed data inserted successfully!")

	return nil
}
