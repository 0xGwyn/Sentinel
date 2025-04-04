package modules

import (
	"fmt"
	// "log"
	"math"

	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
	"github.com/projectdiscovery/httpx/runner"
)

type httpxOutput struct {
	Input           string
	StatusCode      int
	Title           string
	CDNName         string
	CDNType         string
	Technologies    []string
	Failed          bool
	Words           int
	Lines           int
	Port            string
	Location        string
	Hashes          map[string]any
	ResponseHeaders map[string]any
	ContentLength   int
}

func RunHttpx(domains []string, threads int) ([]httpxOutput, error) {
	// Decreasing verbosity level to disable stdout(json output)
	gologger.DefaultLogger.SetMaxLevel(levels.LevelFatal)

	output := []httpxOutput{}

	// var mu sync.Mutex
	options := runner.Options{
		RandomAgent:         true,
		OutputCDN:           "true",
		OutputServerHeader:  true,
		OutputContentType:   true,
		Threads:             threads,
		Methods:             "GET",
		InputTargetHost:     domains,
		FollowRedirects:     true,
		FollowHostRedirects: true,
		JSONOutput:          true,
		// CSVOutput:                 true,
		// DisableStdout:             true,
		TechDetect:                true,
		OutputWordsCount:          true,
		OutputLinesCount:          true,
		ExtractTitle:              true,
		Hashes:                    "sha256",
		ResponseInStdout:          true,
		ResponseHeadersInStdout:   true,
		MaxResponseBodySizeToRead: math.MaxInt32, // Max body size to read
		DisableStdin:              true,
		OnResult: func(r runner.Result) {
			// handle error
			if r.Err != nil {
				fmt.Printf("[Err] %s: %s\n", r.Input, r.Err)
				return
			}
			// mu.Lock()
			output = append(output, httpxOutput{
				Input:           r.Input,
				StatusCode:      r.StatusCode,
				Title:           r.Title,
				CDNName:         r.CDNName,
				CDNType:         r.CDNType,
				Technologies:    r.Technologies,
				Failed:          r.Failed,
				Words:           r.Words,
				Lines:           r.Lines,
				Port:            r.Port,
				Location:        r.Location,
				Hashes:          r.Hashes,
				ResponseHeaders: r.ResponseHeaders,
				ContentLength:   r.ContentLength,
			})
			/*fmt.Printf("Input: %#v\n", r.Input)
			fmt.Printf("Error string: %#v\n", r.Error)
			fmt.Printf("URL: %#v\n", r.URL)
			fmt.Printf("Port: %#v\n", r.Port)
			fmt.Printf("status-code: %#v\n", r.StatusCode)
			fmt.Printf("title: %#v\n", r.Title)
			fmt.Printf("cdn: %#v\n", r.CDN)
			fmt.Printf("cdn name: %#v\n", r.CDNName)
			fmt.Printf("cdn type: %#v\n", r.CDNType)
			fmt.Printf("location: %#v\n", r.Location)
			fmt.Printf("tech: %#v\n", r.Technologies)
			fmt.Printf("hashes: %#v\n", r.Hashes)
			fmt.Printf("failed: %#v\n", r.Failed)
			fmt.Printf("respones-headers: %#v\n", r.ResponseHeaders)
			fmt.Printf("words: %#v\n", r.Words)
			fmt.Printf("lines: %#v\n", r.Lines)
			fmt.Printf("content-length: %#v\n", r.ContentLength)*/
			// Silent:                    true,
			// NoColor:                   true,
			// Proxy:                   "socks5://127.0.0.1:1080",
			// Probe: true,
			// VHost: true,
			// ExtractFqdn:        true,
			// MaxResponseBodySizeToSave: 2147483647, // Max body size to save
			// OnClose: func() {
			// 	fmt.Println("ENDING")
			// },
			// fmt.Printf("fqdn: %#v\n", r.Fqdns)
			// fmt.Printf("tech details: %#v\n", r.TechnologyDetails)
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
			// fmt.Printf("\n\n")
			// mu.Unlock()

		},
	}

	// if err := options.ValidateOptions(); err != nil {
	// 	return fmt.Errorf("option validation failed: %v", err)
	// }

	httpxRunner, err := runner.New(&options)
	if err != nil {
		return nil, fmt.Errorf("creating httpx runner failed: %v", err)
	}

	httpxRunner.RunEnumeration()
	httpxRunner.Close()

	// Resetting logger level
	gologger.DefaultLogger.SetMaxLevel(levels.LevelVerbose)

	return output, nil
}
