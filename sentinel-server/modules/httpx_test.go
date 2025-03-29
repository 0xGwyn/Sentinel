package modules

import (
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

	err := RunHttpx([]string{"memoryleaks.ir", "walmart.com"}, 50)
	if err != nil {
		panic(err)
	}
}
