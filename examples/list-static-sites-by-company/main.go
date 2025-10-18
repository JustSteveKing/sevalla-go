package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/juststeveking/sevalla-go"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("SEVALLA_API_KEY")
	if apiKey == "" {
		log.Fatal("SEVALLA_API_KEY environment variable is required")
	}

	// Get company ID from environment variable or command line argument
	companyID := os.Getenv("SEVALLA_COMPANY_ID")
	if companyID == "" && len(os.Args) > 1 {
		companyID = os.Args[1]
	}
	if companyID == "" {
		log.Fatal("SEVALLA_COMPANY_ID environment variable or command line argument is required")
	}

	// Create client
	client := sevalla.NewClient(
		sevalla.WithAPIKey(apiKey),
	)

	// List static sites for a specific company
	ctx := context.Background()
	opts := &sevalla.ListOptions{
		CompanyID: companyID,
		PerPage:   10,
	}

	fmt.Printf("Fetching static sites for company: %s\n\n", companyID)

	sites, resp, err := client.StaticSites.List(ctx, opts)
	if err != nil {
		log.Fatalf("Error listing static sites: %v", err)
	}

	// Print pagination info
	if resp.NextPage > 0 {
		fmt.Printf("Pagination: Page %d (Next: %d, Last: %d)\n\n",
			resp.NextPage-1, resp.NextPage, resp.LastPage)
	}

	// Display results
	if len(sites) == 0 {
		fmt.Println("No static sites found for this company.")
		return
	}

	fmt.Printf("Found %d static site(s):\n\n", len(sites))

	for i, site := range sites {
		fmt.Printf("%d. %s (ID: %s)\n", i+1, site.Name, site.ID)
		fmt.Printf("   State: %s\n", site.State)
		fmt.Printf("   Region: %s\n", site.Region)
		if site.URL != "" {
			fmt.Printf("   URL: %s\n", site.URL)
		}
		if site.RepositoryURL != "" {
			fmt.Printf("   Repository: %s (Branch: %s)\n", site.RepositoryURL, site.Branch)
		}
		fmt.Printf("   Auto Deploy: %v\n", site.AutoDeploy)
		fmt.Printf("   CDN Enabled: %v\n", site.CDNEnabled)
		fmt.Printf("   SSL Enabled: %v\n", site.SSLEnabled)
		fmt.Println()
	}
}
