package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// Declare the constants that contain the strings to the colors
const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorReset  = "\033[0m"
)

func main() {
	filePath := getInput("Inform the filepath to the docs you want to check: ")
	file, err := openFile(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	content, err := readFileContent(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	text := string(content)

	urlPatternHttps := `https?://[^\s()]+`
	urlPatternEN := `/en/documentation/products/[^\s/()]+(/[^/\s/()]+)*`
	urlPatternPT := `/pt-br/documentacao/produtos/[^\s/()]+(/[^/\s/()]+)*`

	matchesHttp := findMatches(text, urlPatternHttps)
	matchesEN := findMatches(text, urlPatternEN)
	matchesPT := findMatches(text, urlPatternPT)

	for i := range matchesEN {
		matchesEN[i] = formatURL(matchesEN[i])
	}

	for i := range matchesPT {
		matchesPT[i] = formatURL(matchesPT[i])
	}

	allURLs := append(matchesEN, append(matchesHttp, matchesPT...)...)

	fmt.Println("Testing links")
	for _, link := range allURLs {
		statusCode, err := checkURL(link)
		if err != nil {
			fmt.Println(colorRed, "Link:", link, "- Error:", err.Error(), colorReset)
			continue
		}

		color := colorGreen
		if statusCode > 400 {
			color = colorRed
		} else if statusCode >= 300 && statusCode < 400 {
			color = colorYellow
		}

		fmt.Println(color, "Link:", link, "- Status Code:", statusCode, colorReset)
	}
}

// Get the input informed by the user
func getInput(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

// Open the file
func openFile(filePath string) (*os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// Read the content inside the file (.mdx)
func readFileContent(file *os.File) ([]byte, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	content := make([]byte, stat.Size())
	_, err = file.Read(content)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// Find mathches by the pattern (regex) informed
func findMatches(text, pattern string) []string {
	re := regexp.MustCompile(pattern)
	return re.FindAllString(text, -1)
}

// Format the URL adding the host
func formatURL(url string) string {
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	return "https://www.azion.com" + url
}

// Check if the link is working
func checkURL(link string) (int, error) {
	resp, err := http.Get(link)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}
