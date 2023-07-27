package common

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
	// logutil "github.com/projectdiscovery/utils/log"
)

func RunSubfinder(domain string) ([]string, error) {

	subfinderOpts := &runner.Options{
		Silent:             true,
		All:                true,
		Threads:            10,
		Timeout:            30,
		MaxEnumerationTime: 10,
		Config:             filepath.Join(userHomeDir(), ".config/subfinder/config.yaml"),
		ProviderConfig:     filepath.Join(userHomeDir(), ".config/subfinder/provider-config.yaml"),
	}

	// disable timestamps in logs / configure logger
	// logutil.DisableDefaultLogger()

	// making gologger silent
	gologger.DefaultLogger.SetMaxLevel(levels.LevelSilent)
	subfinder, err := runner.NewRunner(subfinderOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create subfinder runner: %v", err)
	}

	output := &bytes.Buffer{}
	// To run subdomain enumeration on a single domain
	if err = subfinder.EnumerateSingleDomainWithCtx(context.Background(), domain, []io.Writer{output}); err != nil {
		return nil, fmt.Errorf("failed to enumerate single domain(%v): %v", domain, err)
	}

	// Convert string to []string
	subdomains := strings.Split(output.String(), "\n")

	// Remove the last element if it's an empty string
	if subdomains[len(subdomains)-1] == "" {
		subdomains = subdomains[:len(subdomains)-1]
	}

	// enable timestamps in logs / configure logger
	// logutil.EnableDefaultLogger()

	return subdomains, nil
}
