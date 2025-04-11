package modules

import (
	"log"
	"testing"
)

func TestRunSubfinder(t *testing.T) {
	subdomains, err := RunSubfinder("projectdiscovery.io")
	if err != nil {
		panic(err)
	}
	for index, output := range subdomains {
		log.Println(index, output.Subdomain, output.Provider)
	}
}
