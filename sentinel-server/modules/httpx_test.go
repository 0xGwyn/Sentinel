package modules

import (
	"fmt"
	"testing"
)

func TestRunHttpx(t *testing.T) {
	// test is a file with multiple domains
	// file, err := os.Open("test")
	// if err != nil {
	// 	fmt.Println("Error opening file:", err)
	// 	return
	// }
	// defer file.Close()

	// scanner := bufio.NewScanner(file)
	// var lines []string
	// for scanner.Scan() {
	// 	lines = append(lines, scanner.Text())
	// }
	// if err := scanner.Err(); err != nil {
	// 	fmt.Println("Error scanning file:", err)
	// 	return
	// }
	subs := []string{"memoryleaks.ir:80", "https://walmart.com:8000", "docs.projectdiscovery.io", "www.cloudflare.com:8080"}
	// subs := []string{"www.cloudflare.com"}
	output, err := RunHttpx(subs, 50)
	if err != nil {
		panic(err)
	}
	for _, r := range output {
		fmt.Printf("Input: %#v\n", r.Input)
		fmt.Printf("status-code: %#v\n", r.StatusCode)
		fmt.Printf("title: %#v\n", r.Title)
		fmt.Printf("cdn name: %#v\n", r.CDNName)
		fmt.Printf("cdn type: %#v\n", r.CDNType)
		fmt.Printf("tech: %#v\n", r.Technologies)
		fmt.Printf("failed: %#v\n", r.Failed)
		fmt.Printf("words: %#v\n", r.Words)
		fmt.Printf("lines: %#v\n", r.Lines)
		fmt.Printf("Port: %#v\n", r.Port)
		fmt.Printf("location: %#v\n", r.Location)
		fmt.Printf("hashes: %#v\n", r.Hashes)
		fmt.Printf("respones-headers: %#v\n", r.ResponseHeaders)
		fmt.Printf("content-length: %#v\n", r.ContentLength)
		fmt.Printf("\n\n")
	}
}
