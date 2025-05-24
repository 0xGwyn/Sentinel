package models

import "go.mongodb.org/mongo-driver/v2/bson"

type StatusType string

const (
	// DNS STATUS

	// freshly added subdomains
	FreshSubdomain StatusType = "fresh_subdomain"
	// freshly resolved subdomains
	FreshResolved StatusType = "fresh_resolved"
	// subdomains that were resolved to ip address but now are not resolved to any ip
	LastResolved StatusType = "last_resolved"
	// subdomains that are resolved but not fresh anymore
	ResolvedSubdomain StatusType = "resolved_subdomain"
	// subdomains that are not resolved to any ip after being added to the database for the first time
	UnresolvedSubdomain StatusType = "unresolved_subdomain"

	// HTTP STATUS

	// subdomains that have an http service behind and are fresh (first time checking for the service)
	FreshService StatusType = "fresh_service"
	// subdomains that have an http service behind them and are not fresh
	NormalService StatusType = "normal_service"
	// subdomains that have an http service behind them but the service status code has changed
	ChangedService StatusType = "changed_service"
	// subdomains that had an http service behind them but the service is no longer available
	LastService StatusType = "last_service"
)

type Domain struct {
	Name       string   `json:"name,omitempty" bson:"name"`
	InScope    []string `json:"in_scope,omitempty" bson:"in_scope"`
	OutOfScope []string `json:"out_of_scope,omitempty" bson:"out_of_scope"`
}

type Subdomain struct {
	Domain    string        `json:"domain,omitempty" bson:"domain"`
	Name      string        `json:"name,omitempty" bson:"name"`
	CreatedAt bson.DateTime `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt bson.DateTime `json:"updated_at,omitempty" bson:"updated_at"`
	Providers []string      `json:"providers,omitempty" bson:"providers"`
	WatchHTTP bool          `json:"watch_http,omitempty" bson:"watch_http"`
	WatchDNS  bool          `json:"watch_dns,omitempty" bson:"watch_dns"`

	// Status types of a subdomain
	DNSStatus  StatusType `json:"dns_status,omitempty" bson:"dns_status"`
	HTTPStatus StatusType `json:"http_status,omitempty" bson:"http_status"`
}

type HTTP struct {
	ScanningDate    bson.DateTime  `json:"scanning_date,omitempty" bson:"scanning_date"`
	Domain          string         `json:"domain,omitempty" bson:"domain"`
	Subdomain       string         `json:"subdomain,omitempty" bson:"subdomain"`
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
	ResponseHeaders map[string]any `json:"response_headers,omitempty" bson:"response_headers"`
	ContentLength   int            `json:"content_length,omitempty" bson:"content_length"`
}

type DNS struct {
	ResolutionDate bson.DateTime `json:"resolution_date,omitempty" bson:"resolution_date"`
	Domain         string        `json:"domain,omitempty" bson:"domain"`
	Subdomain      string        `json:"subdomain,omitempty" bson:"subdomain"`
	CnameRecords   []string      `json:"cname_records,omitempty" bson:"cname_records"`
	ARecords       []string      `json:"a_records,omitempty" bson:"a_records"`
	AAAARecords    []string      `json:"aaaa_records,omitempty" bson:"aaaa_records"`
	NSRecords      []string      `json:"ns_records,omitempty" bson:"ns_records"`
	PTRRecords     []string      `json:"ptr_records,omitempty" bson:"ptr_records"`
	MXRecords      []string      `json:"mx_records,omitempty" bson:"mx_records"`
	TXTRecords     []string      `json:"txt_records,omitempty" bson:"txt_records"`
}
