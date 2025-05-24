package database

import (
	"context"
	"time"

	"github.com/0xgwyn/sentinel/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func InsertMockData() error {
	now := time.Now()

	// Create domain documents
	domain1 := models.Domain{
		Name:       "microsoft.com",
		InScope:    []string{"*.microsoft.com", "*.azure.com"},
		OutOfScope: []string{"*.test.microsoft.com", "internal.microsoft.com"},
	}

	domain2 := models.Domain{
		Name:       "meta.com",
		InScope:    []string{"*.meta.com", "*.facebook.com", "*.instagram.com"},
		OutOfScope: []string{"*.internal.meta.com", "*.dev.meta.com"},
	}

	// Create subdomains
	sub1_1 := models.Subdomain{
		Domain:     "microsoft.com",
		Name:       "www.microsoft.com",
		CreatedAt:  bson.NewDateTimeFromTime(now.Add(-72 * time.Hour)),
		UpdatedAt:  bson.NewDateTimeFromTime(now),
		Providers:  []string{"subfinder", "crtsh"},
		WatchDNS:   true,
		WatchHTTP:  true,
		DNSStatus:  models.FreshSubdomain,
		HTTPStatus: models.FreshService,
	}

	sub1_2 := models.Subdomain{
		Domain:     "microsoft.com",
		Name:       "docs.microsoft.com",
		CreatedAt:  bson.NewDateTimeFromTime(now.Add(-48 * time.Hour)),
		UpdatedAt:  bson.NewDateTimeFromTime(now),
		Providers:  []string{"subfinder", "abuseipdb"},
		WatchDNS:   true,
		WatchHTTP:  true,
		DNSStatus:  models.FreshResolved,
		HTTPStatus: models.FreshService,
	}

	sub1_3 := models.Subdomain{
		Domain:     "microsoft.com",
		Name:       "login.microsoft.com",
		CreatedAt:  bson.NewDateTimeFromTime(now.Add(-24 * time.Hour)),
		UpdatedAt:  bson.NewDateTimeFromTime(now),
		Providers:  []string{"subfinder"},
		WatchDNS:   true,
		WatchHTTP:  true,
		DNSStatus:  models.FreshSubdomain,
		HTTPStatus: models.FreshService,
	}

	sub2_1 := models.Subdomain{
		Domain:     "meta.com",
		Name:       "www.meta.com",
		CreatedAt:  bson.NewDateTimeFromTime(now.Add(-96 * time.Hour)),
		UpdatedAt:  bson.NewDateTimeFromTime(now),
		Providers:  []string{"subfinder", "crtsh", "censys"},
		WatchDNS:   true,
		WatchHTTP:  true,
		DNSStatus:  models.FreshSubdomain,
		HTTPStatus: models.FreshService,
	}

	sub2_2 := models.Subdomain{
		Domain:     "meta.com",
		Name:       "developers.meta.com",
		CreatedAt:  bson.NewDateTimeFromTime(now.Add(-72 * time.Hour)),
		UpdatedAt:  bson.NewDateTimeFromTime(now),
		Providers:  []string{"subfinder", "censys"},
		WatchDNS:   true,
		WatchHTTP:  true,
		DNSStatus:  models.LastResolved,
		HTTPStatus: models.FreshService,
	}

	// Insert domains and subdomains
	domainColl := GetDBCollection("domains")
	subdomainColl := GetDBCollection("subdomains")
	dnsColl := GetDBCollection("dns")
	httpColl := GetDBCollection("http")

	_, err := domainColl.InsertMany(context.Background(), []models.Domain{domain1, domain2})
	if err != nil {
		return err
	}

	subdomains := []models.Subdomain{sub1_1, sub1_2, sub1_3, sub2_1, sub2_2}
	_, err = subdomainColl.InsertMany(context.Background(), subdomains)
	if err != nil {
		return err
	}

	// DNS records
	dnsRecords := []models.DNS{
		{Domain: "microsoft.com", Subdomain: "www.microsoft.com", ResolutionDate: bson.NewDateTimeFromTime(now), CnameRecords: []string{"www-microsoft-com.akadns.net"}, ARecords: []string{"23.45.229.117", "23.45.229.118"}, NSRecords: []string{"ns1.msft.net", "ns2.msft.net"}, MXRecords: []string{"microsoft-com.mail.protection.outlook.com"}},
		{Domain: "microsoft.com", Subdomain: "www.microsoft.com", ResolutionDate: bson.NewDateTimeFromTime(now.Add(-24 * time.Hour)), CnameRecords: []string{"www-microsoft-com.akadns.net"}, ARecords: []string{"23.45.229.120", "23.45.229.121"}, NSRecords: []string{"ns1.msft.net", "ns2.msft.net"}},
		{Domain: "microsoft.com", Subdomain: "www.microsoft.com", ResolutionDate: bson.NewDateTimeFromTime(now.Add(-48 * time.Hour)), ARecords: []string{"23.45.229.125", "23.45.229.126"}, NSRecords: []string{"ns1.msft.net", "ns2.msft.net"}},

		{Domain: "microsoft.com", Subdomain: "docs.microsoft.com", ResolutionDate: bson.NewDateTimeFromTime(now), CnameRecords: []string{"docs.microsoft.com.edgekey.net"}, ARecords: []string{"104.72.155.182", "104.72.155.183"}, NSRecords: []string{"ns1-204.azure-dns.com", "ns2-204.azure-dns.net"}},
		{Domain: "microsoft.com", Subdomain: "docs.microsoft.com", ResolutionDate: bson.NewDateTimeFromTime(now.Add(-24 * time.Hour)), ARecords: []string{"104.72.155.184", "104.72.155.185"}, NSRecords: []string{"ns1-204.azure-dns.com", "ns2-204.azure-dns.net"}},
		{Domain: "microsoft.com", Subdomain: "docs.microsoft.com", ResolutionDate: bson.NewDateTimeFromTime(now.Add(-48 * time.Hour)), ARecords: []string{"104.72.155.186", "104.72.155.187"}, NSRecords: []string{"ns1-204.azure-dns.com", "ns2-204.azure-dns.net"}},

		{Domain: "microsoft.com", Subdomain: "login.microsoft.com", ResolutionDate: bson.NewDateTimeFromTime(now), CnameRecords: []string{"login.msa.msidentity.com"}, ARecords: []string{"40.126.31.145", "40.126.31.146"}, NSRecords: []string{"ns1.msft.net", "ns2.msft.net"}},
		{Domain: "microsoft.com", Subdomain: "login.microsoft.com", ResolutionDate: bson.NewDateTimeFromTime(now.Add(-12 * time.Hour)), ARecords: []string{"40.126.31.147", "40.126.31.148"}, NSRecords: []string{"ns1.msft.net", "ns2.msft.net"}},
		{Domain: "microsoft.com", Subdomain: "login.microsoft.com", ResolutionDate: bson.NewDateTimeFromTime(now.Add(-24 * time.Hour)), ARecords: []string{"40.126.31.149", "40.126.31.150"}, NSRecords: []string{"ns1.msft.net", "ns2.msft.net"}},

		{Domain: "meta.com", Subdomain: "www.meta.com", ResolutionDate: bson.NewDateTimeFromTime(now), CnameRecords: []string{"www.meta.com.edgekey.net"}, ARecords: []string{"157.240.214.35", "157.240.214.36"}, NSRecords: []string{"ns1.meta.com", "ns2.meta.com"}, MXRecords: []string{"aspmx.l.google.com"}},
		{Domain: "meta.com", Subdomain: "www.meta.com", ResolutionDate: bson.NewDateTimeFromTime(now.Add(-24 * time.Hour)), ARecords: []string{"157.240.214.37", "157.240.214.38"}, NSRecords: []string{"ns1.meta.com", "ns2.meta.com"}},
		{Domain: "meta.com", Subdomain: "www.meta.com", ResolutionDate: bson.NewDateTimeFromTime(now.Add(-48 * time.Hour)), ARecords: []string{"157.240.214.39", "157.240.214.40"}, NSRecords: []string{"ns1.meta.com", "ns2.meta.com"}},

		{Domain: "meta.com", Subdomain: "developers.meta.com", ResolutionDate: bson.NewDateTimeFromTime(now), CnameRecords: []string{"developers.meta.com.edgekey.net"}, ARecords: []string{"157.240.195.35", "157.240.195.36"}, NSRecords: []string{"ns3.meta.com", "ns4.meta.com"}},
		{Domain: "meta.com", Subdomain: "developers.meta.com", ResolutionDate: bson.NewDateTimeFromTime(now.Add(-24 * time.Hour)), ARecords: []string{"157.240.195.37", "157.240.195.38"}, NSRecords: []string{"ns3.meta.com", "ns4.meta.com"}},
		{Domain: "meta.com", Subdomain: "developers.meta.com", ResolutionDate: bson.NewDateTimeFromTime(now.Add(-48 * time.Hour)), ARecords: []string{"157.240.195.39", "157.240.195.40"}, NSRecords: []string{"ns3.meta.com", "ns4.meta.com"}},
	}
	_, err = dnsColl.InsertMany(context.Background(), dnsRecords)
	if err != nil {
		return err
	}

	// HTTP records
	httpRecords := []models.HTTP{
		{Domain: "microsoft.com", Subdomain: "www.microsoft.com", ScanningDate: bson.NewDateTimeFromTime(now), Location: "https://www.microsoft.com/en-us", StatusCode: 200, Title: "Microsoft - Cloud, Computers, Apps & Gaming", CDNName: "Akamai", CDNType: "Enterprise", Technologies: []string{"Azure", "IIS", "Windows Server"}, Words: 1500, Lines: 450, Port: "443", ContentLength: 245678, ResponseHeaders: map[string]any{"Server": "Microsoft-IIS/10.0", "X-Powered-By": "ASP.NET"}, Hashes: map[string]any{"sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"}},
		{Domain: "microsoft.com", Subdomain: "www.microsoft.com", ScanningDate: bson.NewDateTimeFromTime(now.Add(-24 * time.Hour)), StatusCode: 200, Title: "Microsoft - Cloud, Computers, Apps & Gaming", CDNName: "Akamai", Technologies: []string{"Azure", "IIS", "Windows Server"}, Words: 1450, Lines: 440, Port: "443", ContentLength: 245000},

		{Domain: "microsoft.com", Subdomain: "docs.microsoft.com", ScanningDate: bson.NewDateTimeFromTime(now), Location: "https://docs.microsoft.com/en-us", StatusCode: 200, Title: "Microsoft Docs", CDNName: "Akamai", CDNType: "Enterprise", Technologies: []string{"Azure", "Next.js", "React"}, Words: 2000, Lines: 600, Port: "443", ContentLength: 356789, ResponseHeaders: map[string]any{"Server": "AkamaiGHost", "X-Azure-Ref": "0123456789"}, Hashes: map[string]any{"sha256": "f3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"}},
		{Domain: "microsoft.com", Subdomain: "docs.microsoft.com", ScanningDate: bson.NewDateTimeFromTime(now.Add(-24 * time.Hour)), StatusCode: 200, Title: "Microsoft Docs", CDNName: "Akamai", Technologies: []string{"Azure", "Next.js", "React"}, Words: 1950, Lines: 580, Port: "443", ContentLength: 350000},

		{Domain: "microsoft.com", Subdomain: "login.microsoft.com", ScanningDate: bson.NewDateTimeFromTime(now), Location: "https://login.live.com", StatusCode: 302, Title: "Microsoft Account Login", CDNName: "Azure", CDNType: "Enterprise", Technologies: []string{"Azure AD", "OAuth"}, Words: 500, Lines: 150, Port: "443", ContentLength: 125678, ResponseHeaders: map[string]any{"Strict-Transport-Security": "max-age=31536000"}, Hashes: map[string]any{"sha256": "a3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"}},
		{Domain: "microsoft.com", Subdomain: "login.microsoft.com", ScanningDate: bson.NewDateTimeFromTime(now.Add(-12 * time.Hour)), StatusCode: 302, Title: "Microsoft Account Login", CDNName: "Azure", Technologies: []string{"Azure AD", "OAuth"}, Words: 480, Lines: 145, Port: "443", ContentLength: 124000},

		{Domain: "meta.com", Subdomain: "www.meta.com", ScanningDate: bson.NewDateTimeFromTime(now), Location: "https://www.meta.com", StatusCode: 200, Title: "Meta - Home", CDNName: "Akamai", CDNType: "Enterprise", Technologies: []string{"React", "GraphQL", "Node.js"}, Words: 2500, Lines: 800, Port: "443", ContentLength: 567890, ResponseHeaders: map[string]any{"x-fb-debug": "xyz123", "x-frame-options": "DENY"}, Hashes: map[string]any{"sha256": "b3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"}},
		{Domain: "meta.com", Subdomain: "www.meta.com", ScanningDate: bson.NewDateTimeFromTime(now.Add(-24 * time.Hour)), StatusCode: 200, Title: "Meta - Home", CDNName: "Akamai", Technologies: []string{"React", "GraphQL", "Node.js"}, Words: 2400, Lines: 780, Port: "443", ContentLength: 565000},

		{Domain: "meta.com", Subdomain: "developers.meta.com", ScanningDate: bson.NewDateTimeFromTime(now), Location: "https://developers.meta.com", StatusCode: 200, Title: "Meta for Developers", CDNName: "Akamai", CDNType: "Enterprise", Technologies: []string{"React", "GraphQL", "Express"}, Words: 3000, Lines: 1000, Port: "443", ContentLength: 789012, ResponseHeaders: map[string]any{"x-fb-debug": "abc789", "x-content-type-options": "nosniff"}, Hashes: map[string]any{"sha256": "c3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"}},
		{Domain: "meta.com", Subdomain: "developers.meta.com", ScanningDate: bson.NewDateTimeFromTime(now.Add(-24 * time.Hour)), StatusCode: 200, Title: "Meta for Developers", CDNName: "Akamai", Technologies: []string{"React", "GraphQL", "Express"}, Words: 2900, Lines: 980, Port: "443", ContentLength: 785000},
	}
	_, err = httpColl.InsertMany(context.Background(), httpRecords)
	if err != nil {
		return err
	}

	return nil
}
