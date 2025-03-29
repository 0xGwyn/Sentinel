package modules

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
	folderutil "github.com/projectdiscovery/utils/folder"
	// logutil "github.com/projectdiscovery/utils/log"
)

func RunSubfinder(domain string) (map[string][]string, error) {

	subfinderOpts := &runner.Options{
		Silent: true,
		// Verbose:            true,
		All:                true,
		Version:            true,
		Threads:            10,
		Timeout:            30,
		MaxEnumerationTime: 10,
		// somehow the config or provider variables don't apply to the actual program
		Config:         filepath.Join(folderutil.HomeDirOrDefault(""), ".config/subfinder/config.yaml"),
		ProviderConfig: filepath.Join(folderutil.HomeDirOrDefault(""), ".config/subfinder/provider-config.yaml"),
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
	var sourceMap map[string]map[string]struct{}
	// To run subdomain enumeration on a single domain
	if sourceMap, err = subfinder.EnumerateSingleDomainWithCtx(context.Background(), domain, []io.Writer{output}); err != nil {
		return nil, fmt.Errorf("failed to enumerate single domain(%v): %v", domain, err)
	}

	// use sourceMap to access the results in the application
	subdomains := make(map[string][]string, len(sourceMap))
	for subdomain, sources := range sourceMap {
		sourcesList := make([]string, 0, len(sources))
		for source := range sources {
			sourcesList = append(sourcesList, source)
		}
		subdomains[subdomain] = sourcesList
	}

	// enable timestamps in logs / configure logger
	// logutil.EnableDefaultLogger()

	return subdomains, nil
}
