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
		RandomAgent:        true,
		OutputCDN:          "true",
		OutputServerHeader: true,
		OutputContentType:  true,
		Threads:            threads,
		Methods:            "GET",
		InputTargetHost:    domains,
		FollowRedirects:    true,
		JSONOutput:         true,
		OutputWordsCount:   true,
		// Proxy:                   "socks5://127.0.0.1:1080",
		// Probe:                   true,
		OutputLinesCount:        true,
		ExtractTitle:            true,
		Hashes:                  "sha256",
		ResponseHeadersInStdout: true,
		// ExtractFqdn:        true,

		// MaxResponseBodySizeToSave: 2147483647, // Max body size to save
		MaxResponseBodySizeToRead: 2147483647, // Max body size to read
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
			fmt.Printf("cdn: %#v\n", r.CDN)
			fmt.Printf("cdn name: %#v\n", r.CDNName)
			fmt.Printf("cdn type: %#v\n", r.CDNType)
			fmt.Printf("location: %#v\n", r.Location)
			fmt.Printf("tech: %#v\n", r.Technologies)
			// fmt.Printf("tech details: %#v\n", r.TechnologyDetails)
			fmt.Printf("hashes: %#v\n", r.Hashes)
			fmt.Printf("failed: %#v\n", r.Failed)
			fmt.Printf("respones-headers: %#v\n", r.ResponseHeaders)
			fmt.Printf("words: %#v\n", r.Words)
			fmt.Printf("lines: %#v\n", r.Lines)
			fmt.Printf("content-length: %#v\n", r.ContentLength)
			// fmt.Printf("fqdn: %#v\n", r.Fqdns)
			/*fmt.Printf("json: %#v\n", r.JSON(&runner.ScanOptions{
				OutputTitle:               true,
				OutputStatusCode:          true,
				OutputLocation:            true,
				OutputContentLength:       true,
				OutputServerHeader:        true,
				OutputContentType:         true,
				OutputCDN:                 "true",
				TechDetect:                true,
				MaxResponseBodySizeToRead: 2147483647,
				OutputLinesCount:          true,
				OutputWordsCount:          true,
				Hashes:                    "sha256",
			}))*/
			fmt.Printf("\n\n")
			mu.Unlock()

		},
	}

	// if err := options.ValidateOptions(); err != nil {
	// 	return fmt.Errorf("option validation failed: %v", err)
	// }

	httpxRunner, err := runner.New(&options)
	if err != nil {
		return fmt.Errorf("creating httpx runner failed: %v", err)
	}

	httpxRunner.RunEnumeration()
	httpxRunner.Close()

	return nil
}
