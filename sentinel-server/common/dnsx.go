package common

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/miekg/dns"
	"github.com/projectdiscovery/dnsx/libs/dnsx"
	sliceutil "github.com/projectdiscovery/utils/slice"
)

type dnsQueryOutput struct {
	domain  string
	records map[string][]string
}

func RunDnsx(domains, questionTypes []string, threads int) ([]dnsQueryOutput, error) {

	defaultOptions := dnsx.DefaultOptions

	var wg sync.WaitGroup
	numOfGoroutines := min(threads, len(domains))
	wg.Add(numOfGoroutines)

	// setting query types (A, CNAME, ...)
	defaultOptions.QuestionTypes = getQuesntionTypes(questionTypes)

	// create DNS Resolver with default options
	dnsClient, err := dnsx.New(defaultOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create dns client: %v", err)
	}

	// domains are given as input to the go routines
	input := make(chan string)

	// a struct containing a domain and multiple records are given to go routines as output
	output := make(chan dnsQueryOutput)

	for range numOfGoroutines {
		go dnsWorker(dnsClient, input, output, &wg)
	}

	// providing workers with domains
	go func() {
		for _, domain := range domains {
			input <- domain
		}
		close(input)
	}()

	// results variable is an array of domains and their corresponding records
	results := []dnsQueryOutput{}
	go func() {
		for queryOutput := range output {
			results = append(results, queryOutput)
		}
		close(output)
	}()

	// waiting for all the dns workers to be done
	wg.Wait()

	return results, nil
}

// dns workers resolve domains and write the records to the output chan as a dnsQueryOutput type
func dnsWorker(dnsClient *dnsx.DNSX, domains <-chan string, output chan<- dnsQueryOutput, wg *sync.WaitGroup) {
	defer wg.Done()
	// get requested record types for the given domains
	for domain := range domains {
		queryResponse := make(map[string][]string)
		rawResp, err := dnsClient.QueryMultiple(domain)
		if err != nil {
			log.Printf("failed resolving %v : %v\n", domain, err)
		}
		if 0 < len(rawResp.A) {
			queryResponse["a"] = rawResp.A
		}
		if 0 < len(rawResp.NS) {
			queryResponse["ns"] = rawResp.NS
		}
		if 0 < len(rawResp.AAAA) {
			queryResponse["aaaa"] = rawResp.AAAA
		}
		if 0 < len(rawResp.CNAME) {
			queryResponse["cname"] = rawResp.CNAME
		}
		if 0 < len(rawResp.PTR) {
			queryResponse["ptr"] = rawResp.PTR
		}
		if 0 < len(rawResp.TXT) {
			queryResponse["txt"] = rawResp.TXT
		}
		if 0 < len(rawResp.SRV) {
			queryResponse["srv"] = rawResp.SRV
		}
		if 0 < len(rawResp.MX) {
			queryResponse["mx"] = rawResp.MX
		}
		if 0 < len(rawResp.CAA) {
			queryResponse["caa"] = rawResp.CAA
		}
		output <- dnsQueryOutput{domain: domain, records: queryResponse}
	}
}

// converts a list of string record types to the corresponding integer types
// []string{"a", "cname", ...} -> []uint16{dns.TypeA, dns.TypeCNAME}
func getQuesntionTypes(questions []string) []uint16 {
	// convert all records to lowercase strings
	lowercaseQuestions := make([]string, 0)
	for _, question := range questions {
		lowercaseQuestions = append(lowercaseQuestions, strings.ToLower(question))
	}

	var questionTypes []uint16
	if sliceutil.Contains(lowercaseQuestions, "a") {
		questionTypes = append(questionTypes, dns.TypeA)
	}
	if sliceutil.Contains(lowercaseQuestions, "aaaa") {
		questionTypes = append(questionTypes, dns.TypeAAAA)
	}
	if sliceutil.Contains(lowercaseQuestions, "cname") {
		questionTypes = append(questionTypes, dns.TypeCNAME)
	}
	if sliceutil.Contains(lowercaseQuestions, "ptr") {
		questionTypes = append(questionTypes, dns.TypePTR)
	}
	if sliceutil.Contains(lowercaseQuestions, "soa") {
		questionTypes = append(questionTypes, dns.TypeSOA)
	}
	if sliceutil.Contains(lowercaseQuestions, "txt") {
		questionTypes = append(questionTypes, dns.TypeTXT)
	}
	if sliceutil.Contains(lowercaseQuestions, "srv") {
		questionTypes = append(questionTypes, dns.TypeSRV)
	}
	if sliceutil.Contains(lowercaseQuestions, "mx") {
		questionTypes = append(questionTypes, dns.TypeMX)
	}
	if sliceutil.Contains(lowercaseQuestions, "ns") {
		questionTypes = append(questionTypes, dns.TypeNS)
	}
	if sliceutil.Contains(lowercaseQuestions, "caa") {
		questionTypes = append(questionTypes, dns.TypeCAA)
	}

	return questionTypes
}
