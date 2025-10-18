# List Static Sites by Company ID

This example demonstrates how to list static sites filtered by a specific company ID using the Sevalla Go SDK.

## Usage

### Option 1: Using environment variables

```bash
export SEVALLA_API_KEY="your-api-key"
export SEVALLA_COMPANY_ID="your-company-id"
go run main.go
```

### Option 2: Using command line argument

```bash
export SEVALLA_API_KEY="your-api-key"
go run main.go your-company-id
```

## Example Code

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/juststeveking/sevalla-go"
)

func main() {
	// Create client
	client := sevalla.NewClient(
		sevalla.WithAPIKey("your-api-key"),
	)

	// List static sites for a specific company
	ctx := context.Background()
	opts := &sevalla.ListOptions{
		CompanyID: "your-company-id",
		PerPage:   10,
	}

	sites, resp, err := client.StaticSites.List(ctx, opts)
	if err != nil {
		log.Fatalf("Error listing static sites: %v", err)
	}

	// Process results
	for _, site := range sites {
		fmt.Printf("Site: %s (ID: %s)\n", site.Name, site.ID)
	}
}
```

## Features

The `ListOptions` struct now supports the following fields:

- `Page` - Page number for pagination
- `PerPage` - Number of items per page
- `Sort` - Field to sort by
- `Order` - Sort order (asc/desc)
- `CompanyID` - **NEW**: Filter results by company ID

## Notes

- The `CompanyID` parameter is optional. If not provided, all accessible static sites will be returned.
- The company ID filtering is implemented as a query parameter and works consistently across all list operations.
- Pagination information is available in the response to handle large result sets.
