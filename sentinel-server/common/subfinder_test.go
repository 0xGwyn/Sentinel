package common

import (
	"fmt"
	"testing"
)

func TestRunSubfinder(t *testing.T) {
	subdomains, err := RunSubfinder("projectdiscovery.io")
	if err != nil {
		panic(err)
	}
	fmt.Println(subdomains)
}
