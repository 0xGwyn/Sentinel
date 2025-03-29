package common

import (
	"log"
	"testing"
)

func TestRunSubfinder(t *testing.T) {
	subdomains, err := RunSubfinder("projectdiscovery.io")
	if err != nil {
		panic(err)
	}
	for subdomain, source := range subdomains {
		log.Println(subdomain, source)
	}
}
