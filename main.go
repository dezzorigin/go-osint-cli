package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Platform struct {
	Name string
	URL  string
}

func main() {
	var username string
	fmt.Print("Enter username to check: ")
	fmt.Scanln(&username)

	if username == "" {
		fmt.Println("Error: Username cannot be empty.")
		return
	}

	platforms := []Platform{
		{Name: "GitHub", URL: "https://github.com/%s"},
		{Name: "GitLab", URL: "https://gitlab.com/%s"},
		{Name: "Dev.to", URL: "https://dev.to/%s"},
		{Name: "Pinterest", URL: "https://www.pinterest.com/%s/"},
		{Name: "Reddit", URL: "https://www.reddit.com/user/%s/"},
	}

	fmt.Printf("\nScanning for '%s' across %d platforms...\n\n", username, len(platforms))

	var wg sync.WaitGroup
	startTime := time.Now()

	for _, p := range platforms {
		wg.Add(1)

		go func(target Platform) {
			defer wg.Done()
			urlFull := fmt.Sprintf(target.URL, username)

			client := http.Client{
				Timeout: 10 * time.Second,
			}

			resp, err := client.Get(urlFull)
			if err != nil {
				fmt.Printf("[Error] [%s]: %v\n", target.Name, err)
				return
			}
			defer resp.Body.Close()

			colorReset := "\033[0m"
			colorGreen := "\033[32m"
			colorRed := "\033[31m"

			// Use switch for cleaner idiomatic Go
			switch resp.StatusCode {
			case http.StatusOK:
				fmt.Printf("[%s] %sFound%s -> %s\n", target.Name, colorGreen, colorReset, urlFull)
			case http.StatusNotFound:
				fmt.Printf("[%s] %sNot Found%s (Available)\n", target.Name, colorRed, colorReset)
			default:
				fmt.Printf("[%s] Unexpected Response (Status: %d)\n", target.Name, resp.StatusCode)
			}
		}(p)
	}

	wg.Wait()
	duration := time.Since(startTime)

	fmt.Printf("\nScan completed in: %v\n", duration)
}
