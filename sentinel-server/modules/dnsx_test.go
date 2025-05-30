package modules

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestRunDnsx(t *testing.T) {
	// test is a file with multiple domains
	file, err := os.Open("test")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
		return
	}

	results, err := RunDnsx(lines, []string{"a", "cname", "txt"}, 25)
	if err != nil {
		panic(err)
	}

	for _, result := range results {
		fmt.Printf("%v %v\n", result.Domain, result.Records)
	}

}
