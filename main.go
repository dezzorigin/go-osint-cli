// Command go-osint-cli checks whether a username exists on several public
// platforms. It requests each profile URL and reads the HTTP status code:
//
//	200 -> profile exists ("Found")
//	404 -> profile does not exist ("Not Found", username available)
//	3xx -> the platform redirected (usually a login wall), so the
//	       username cannot be verified from an anonymous request
//
// Usage:
//
//	go-osint-cli <username>  # scan a username directly
//	go-osint-cli             # ask for a username interactively
package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const userAgent = "go-osint-cli/1.1.0 (+https://github.com/dezzorigin/go-osint-cli)"

// Platform is one site whose profile URL follows a fixed pattern.
type Platform struct {
	Name string
	URL  string
}

var platforms = []Platform{
	{Name: "GitHub", URL: "https://github.com/%s"},
	{Name: "GitLab", URL: "https://gitlab.com/%s"},
	{Name: "Dev.to", URL: "https://dev.to/%s"},
	{Name: "Pinterest", URL: "https://www.pinterest.com/%s/"},
	{Name: "Reddit", URL: "https://www.reddit.com/user/%s/"},
}

func main() {
	username := ""
	if len(os.Args) > 1 {
		username = strings.TrimSpace(os.Args[1])
	}
	if username == "" {
		fmt.Print("Enter username to check: ")
		line, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		username = strings.TrimSpace(line)
	}
	if username == "" {
		fmt.Println("Error: Username cannot be empty.")
		os.Exit(1)
	}

	fmt.Printf("\nScanning for '%s' across %d platforms...\n\n", username, len(platforms))

	results := make([]string, len(platforms))
	var wg sync.WaitGroup
	startTime := time.Now()

	for i, p := range platforms {
		wg.Add(1)
		go func(i int, target Platform) {
			defer wg.Done()
			results[i] = check(target, username)
		}(i, p)
	}

	wg.Wait()

	// Print in platform order (not in arrival order) so output is stable.
	for _, r := range results {
		fmt.Println(r)
	}

	fmt.Printf("\nScan completed in: %v\n", time.Since(startTime))
}

// check requests one profile URL and turns the status code into one
// human-readable result line.
func check(p Platform, username string) string {
	urlFull := fmt.Sprintf(p.URL, username)

	client := &http.Client{
		Timeout: 10 * time.Second,
		// Do not follow redirects: a redirect (e.g. to a sign-in page)
		// is itself the answer — the profile cannot be verified.
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest(http.MethodGet, urlFull, nil)
	if err != nil {
		return fmt.Sprintf("[%s] Error: %v", p.Name, err)
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("[%s] Error: %v", p.Name, err)
	}
	defer resp.Body.Close()

	return formatResult(p, urlFull, resp.StatusCode)
}

// formatResult maps an HTTP status code to the result line shown to the
// user. Pure function, covered by main_test.go.
func formatResult(p Platform, urlFull string, status int) string {
	const (
		colorReset  = "\033[0m"
		colorGreen  = "\033[32m"
		colorRed    = "\033[31m"
		colorYellow = "\033[33m"
	)

	switch {
	case status == http.StatusOK:
		return fmt.Sprintf("[%s] %sFound%s -> %s", p.Name, colorGreen, colorReset, urlFull)
	case status == http.StatusNotFound:
		return fmt.Sprintf("[%s] %sNot Found%s (Available)", p.Name, colorRed, colorReset)
	case status >= 300 && status < 400:
		return fmt.Sprintf("[%s] %sRedirect (Status: %d)%s — cannot verify (login wall?)", p.Name, colorYellow, status, colorReset)
	default:
		return fmt.Sprintf("[%s] Unexpected Response (Status: %d)", p.Name, status)
	}
}
