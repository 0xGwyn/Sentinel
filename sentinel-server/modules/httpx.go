package modules

import (
	"fmt"
	"log"
	"sync"

	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
	"github.com/projectdiscovery/httpx/runner"
)

type httpxOutput struct {
	Domain       string
	StatusCode   int
	Title        string
	CDN          string
	Technology   string
	BodySha256   string
	Failed       bool
	HeaderSha256 string
	Words        int
	Lines        int
}

func RunHttpx(domains []string, threads int) error {
	// increase the verbosity (optional)
	gologger.DefaultLogger.SetMaxLevel(levels.LevelVerbose)

	// output := []httpxOutput{}

	var mu sync.Mutex
	options := runner.Options{
		RandomAgent:      true,
		OutputCDN:        "true",
		Threads:          threads,
		Methods:          "GET",
		InputTargetHost:  domains,
		JSONOutput:       true,
		OutputWordsCount: true,
		OutputLinesCount: true,
		OnResult: func(r runner.Result) {
			// handle error
			if r.Err != nil {
				log.Printf("[Err] %s: %s\n", r.Input, r.Err)
				return
			}
			mu.Lock()
			fmt.Printf("URL: %#v\n", r.URL)
			fmt.Printf("status-code: %#v\n", r.StatusCode)
			fmt.Printf("title: %#v\n", r.Title)
			fmt.Printf("cdn: %#v\n", r.CDNName)
			fmt.Printf("tech: %#v\n", r.Technologies)
			fmt.Printf("hashes: %#v\n", r.Hashes)
			fmt.Printf("failed: %#v\n", r.Failed)
			fmt.Printf("respones-headers: %#v\n", r.ResponseHeaders)
			fmt.Printf("words: %#v\n", r.Words)
			fmt.Printf("lines: %#v\n", r.Lines)
			mu.Unlock()

		},
	}

	if err := options.ValidateOptions(); err != nil {
		return fmt.Errorf("option validation failed: %v", err)
	}

	httpxRunner, err := runner.New(&options)
	if err != nil {
		return fmt.Errorf("creating httpx runner failed: %v", err)
	}
	defer httpxRunner.Close()

	httpxRunner.RunEnumeration()

	return nil
}
