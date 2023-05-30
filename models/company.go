package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Company struct {
	ID      primitive.ObjectID `json:"_id" bson:"_id"`
	Name    string             `json:"name,omitempty" bson:"name"`
	Domains []Domain           `json:"domains,omitempty" bson:"domains"`
}

type Domain struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	Name       string             `json:"name,omitempty" bson:"name"`
	Subdomains []Subdomain        `json:"subdomains,omitempty" bson:"subdomains"`
}

type Subdomain struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	Name       string             `json:"name,omitempty" bson:"name"`
	StatusCode string             `json:"statuscode,omitempty" bson:"statuscode"`
	Title      string             `json:"title,omitempty" bson:"title"`
	CDN        string             `json:"CDN,omitempty" bson:"CDN"`
	Technology string             `json:"technology,omitempty" bson:"technology"`
	CreatedAt  primitive.DateTime `json:"createdAt,omitempty" bson:"createdAt"`
	UpdatedAt  primitive.DateTime `json:"updatedAt,omitempty" bson:"updatedAt"`
	Paths      []string           `json:"paths,omitempty" bson:"paths"`
}
