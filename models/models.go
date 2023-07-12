package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Domain struct {
	Name       string   `json:"name,omitempty" bson:"name"`
	InScope    []string `json:"in_scope,omitempty" bson:"in_scope"`
	OutOfScope []string `json:"out_of_scope,omitempty" bson:"out_of_scope"`
}

type Subdomain struct {
	Domain       string             `json:"domain,omitempty" bson:"domain"`
	Name         string             `json:"name,omitempty" bson:"name"`
	CreatedAt    primitive.DateTime `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt    primitive.DateTime `json:"updated_at,omitempty" bson:"updated_at"`
	StatusCode   int                `json:"status_code,omitempty" bson:"status_code"`
	Title        string             `json:"title,omitempty" bson:"title"`
	CDN          string             `json:"CDN,omitempty" bson:"CDN"`
	Technology   string             `json:"technology,omitempty" bson:"technology"`
	Providers    []string           `json:"providers,omitempty" bson:"providers"`
	BodySha256   string             `json:"body_sha256,omitempty" bson:"body_sha256"`
	HeaderSha256 string             `json:"header_sha256,omitempty" bson:"header_sha256"`
	Words        int                `json:"words,omitempty" bson:"words"`
	Lines        int                `json:"lines,omitempty" bson:"lines"`
	Failed       bool               `json:"failed,omitempty" bson:"failed"`
	Watch        bool               `json:"watch,omitempty" bson:"watch"`
	CnameRecords []string           `json:"cname_records,omitempty" bson:"cname_records"`
	ARecords     []string           `json:"a_records,omitempty" bson:"a_records"`
	NSRecords    []string           `json:"ns_records,omitempty" bson:"ns_records"`
	PTRRecords   []string           `json:"ptr_records,omitempty" bson:"ptr_records"`
	MXRecords    []string           `json:"mx_records,omitempty" bson:"mx_records"`
}
