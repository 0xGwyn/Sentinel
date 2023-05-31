package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Domain struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	Name       string             `json:"name,omitempty" bson:"name"`
	Subdomains []Subdomain        `json:"subdomains,omitempty" bson:"subdomains"`
}

type Subdomain struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	Name       string             `json:"name,omitempty" bson:"name"`
	StatusCode string             `json:"statuscode,omitempty" bson:"statuscode"`
	Title      string             `json:"title,omitempty" bson:"title"`
	CDN        string             `json:"CDN,omitempty" bson:"CDN"`
	Technology string             `json:"technology,omitempty" bson:"technology"`
	CreatedAt  primitive.DateTime `json:"createdAt,omitempty" bson:"createdAt"`
	UpdatedAt  primitive.DateTime `json:"updatedAt,omitempty" bson:"updatedAt"`
	IPs        []string           `json:"ips,omitempty" bson:"ips"`
	Providers  []string           `json:"providers,omitempty" bson:"providers"`
	Paths      []string           `json:"paths,omitempty" bson:"paths"`
}
