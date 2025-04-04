package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Domain struct {
	Name       string   `json:"name,omitempty" bson:"name"`
	InScope    []string `json:"in_scope,omitempty" bson:"in_scope"`
	OutOfScope []string `json:"out_of_scope,omitempty" bson:"out_of_scope"`
}

type StatusType string

const (
	// freshly added subdomains
	FreshSubdomain StatusType = "fresh_subdomain"
	// freshly resolved subdomains
	FreshResolved StatusType = "fresh_resolved"
	// subdomains that were resolved to ip address but now are not resolved to any ip
	LastResolved StatusType = "last_resolved"
	// subdomains that are not resolved to any ip after being added to the database for the first time
	UnresolvedSubdomain StatusType = "unresolved_subdomain"
	// subdomains that have an http service behind them
	FreshService StatusType = "fresh_service"
	// subdomains that have an http service behind them but the service status code has changed
	ChangedService StatusType = "changed_service"
	// not any of the above (after a status type change, the status type should change to this after a while.)
	Normal StatusType = "normal"
)

type HTTP struct {
	ScanningDate    bson.DateTime  `json:"scanning_date,omitempty" bson:"scanning_date"`
	Location        string         `json:"location,omitempty" bson:"location"`
	StatusCode      int            `json:"status_code,omitempty" bson:"status_code"`
	Title           string         `json:"title,omitempty" bson:"title"`
	CDNName         string         `json:"cdn_name,omitempty" bson:"cdn_name"`
	CDNType         string         `json:"cdn_type,omitempty" bson:"cdn_type"`
	Technologies    []string       `json:"technologies,omitempty" bson:"technologies"`
	Hashes          map[string]any `json:"hashes,omitempty" bson:"hashes"`
	Words           int            `json:"words,omitempty" bson:"words"`
	Lines           int            `json:"lines,omitempty" bson:"lines"`
	Failed          bool           `json:"failed,omitempty" bson:"failed"`
	Port            string         `json:"port,omitempty" bson:"port"`
	Latest          bool           `json:"latest,omitempty" bson:"latest"`
	ResponseHeaders map[string]any `json:"response_headers,omitempty" bson:"response_headers"`
	ContentLength   int            `json:"content_length,omitempty" bson:"content_length"`
}

type DNS struct {
	ResolutionDate bson.DateTime `json:"resolution_date,omitempty" bson:"resolution_date"`
	CnameRecords   []string      `json:"cname_records,omitempty" bson:"cname_records"`
	ARecords       []string      `json:"a_records,omitempty" bson:"a_records"`
	NSRecords      []string      `json:"ns_records,omitempty" bson:"ns_records"`
	PTRRecords     []string      `json:"ptr_records,omitempty" bson:"ptr_records"`
	MXRecords      []string      `json:"mx_records,omitempty" bson:"mx_records"`
}

type Subdomain struct {
	Domain    string        `json:"domain,omitempty" bson:"domain"`
	Name      string        `json:"name,omitempty" bson:"name"`
	CreatedAt bson.DateTime `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt bson.DateTime `json:"updated_at,omitempty" bson:"updated_at"`
	Providers []string      `json:"providers,omitempty" bson:"providers"`
	Watch     bool          `json:"watch,omitempty" bson:"watch"`

	// DNS data
	DNS []DNS `json:"dns,omitempty" bson:"dns"`

	// HTTP service data
	HTTP []HTTP `json:"http,omitempty" bson:"http"`

	// Status type of the subdomain
	Status StatusType `json:"status,omitempty" bson:"status"`
}
