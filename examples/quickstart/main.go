package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	sevalla "github.com/juststeveking/sevalla-go"
)

func main() {
	var (
		apiKey  = flag.String("api-key", os.Getenv("SEVALLA_API_KEY"), "Sevalla API key")
		action  = flag.String("action", "list", "Action to perform (list, create, deploy, logs)")
		appName = flag.String("app", "", "Application name")
		repoURL = flag.String("repo", "", "Repository URL for new app")
	)
	flag.Parse()

	if *apiKey == "" {
		log.Fatal("API key is required (set SEVALLA_API_KEY or use -api-key flag)")
	}

	// Create client
	client := sevalla.NewClient(
		sevalla.WithAPIKey(*apiKey),
	)

	ctx := context.Background()

	switch *action {
	case "list":
		listApplications(ctx, client)

	case "create":
		if *appName == "" || *repoURL == "" {
			log.Fatal("Both -app and -repo flags are required for create action")
		}
		createApplication(ctx, client, *appName, *repoURL)

	case "deploy":
		if *appName == "" {
			log.Fatal("-app flag is required for deploy action")
		}
		deployApplication(ctx, client, *appName)

	case "logs":
		if *appName == "" {
			log.Fatal("-app flag is required for logs action")
		}
		getApplicationLogs(ctx, client, *appName)

	default:
		log.Fatalf("Unknown action: %s", *action)
	}
}

func listApplications(ctx context.Context, client *sevalla.Client) {
	fmt.Println("📋 Listing applications...")

	apps, resp, err := client.Applications.List(ctx, nil)
	if err != nil {
		handleError(err)
	}

	fmt.Printf("\n✅ Found %d applications (Status: %d)\n\n", len(apps), resp.StatusCode)

	if len(apps) == 0 {
		fmt.Println("No applications found. Create one with: -action create -app myapp -repo https://github.com/user/repo")
		return
	}

	for _, app := range apps {
		fmt.Printf("📦 Application: %s\n", app.Name)
		fmt.Printf("   ID:       %s\n", app.ID)
		fmt.Printf("   Name:     %s\n", app.Name)
		fmt.Printf("   State:    %s\n", app.State)
		fmt.Printf("   Region: %s\n", app.Region)
		fmt.Printf("   Plan: %s\n", app.Plan)
		fmt.Printf("   Replicas: %d\n", app.Replicas)
		if app.URL != "" {
			fmt.Printf("   URL:      %s\n", app.URL)
		}
		fmt.Printf("   Created:  %s\n", app.CreatedAt.Format(time.RFC3339))
		fmt.Println()
	}
}

func createApplication(ctx context.Context, client *sevalla.Client, name, _ string) {
	fmt.Printf("🚀 Creating application '%s'...\n", name)

	createReq := &sevalla.CreateApplicationRequest{
		Name:          "quickstart-app",
		RepositoryURL: "https://github.com/vercel/next.js",
		Branch:        "canary",
		Region:        sevalla.RegionUSCentral,
		Plan:          sevalla.PlanHobby,
		Replicas:      1,
		EnvironmentVars: map[string]string{
			"NODE_ENV": "production",
		},
	}

	app, resp, err := client.Applications.Create(ctx, createReq)
	if err != nil {
		handleError(err)
	}

	fmt.Printf("\n✅ Application created successfully! (Status: %d)\n\n", resp.StatusCode)
	fmt.Printf("📦 Application Details:\n")
	fmt.Printf("   ID:       %s\n", app.ID)
	fmt.Printf("   Name:     %s\n", app.Name)
	fmt.Printf("   State:    %s\n", app.State)
	fmt.Printf("   Location: %s\n", app.Region)
	fmt.Printf("   URL:      %s\n", app.URL)
	fmt.Printf("   Created:  %s\n", app.CreatedAt.Format(time.RFC3339))

	if app.State == sevalla.StateDeploying {
		fmt.Println("\n⏳ Initial deployment in progress...")
		fmt.Println("   Check deployment status with: -action deploy -app", name)
	}
}

func deployApplication(ctx context.Context, client *sevalla.Client, appName string) {
	fmt.Printf("🔄 Deploying application '%s'...\n", appName)

	// First, find the application by name
	apps, _, err := client.Applications.List(ctx, nil)
	if err != nil {
		handleError(err)
	}

	var app *sevalla.Application
	for _, a := range apps {
		if a.Name == appName {
			app = a
			break
		}
	}

	if app == nil {
		log.Fatalf("❌ Application '%s' not found", appName)
	}

	// Trigger deployment
	deployment, resp, err := client.Applications.Deploy(ctx, app.ID)
	if err != nil {
		handleError(err)
	}

	fmt.Printf("\n✅ Deployment started! (Status: %d)\n\n", resp.StatusCode)
	fmt.Printf("🚀 Deployment Details:\n")
	fmt.Printf("   ID:            %s\n", deployment.ID)
	fmt.Printf("   Application:   %s\n", appName)
	fmt.Printf("   State:         %s\n", deployment.State)
	fmt.Printf("   Branch:        %s\n", deployment.Branch)
	fmt.Printf("   Started:       %s\n", deployment.StartedAt.Format(time.RFC3339))

	// Optionally wait for deployment to complete
	fmt.Println("\n⏳ Waiting for deployment to complete...")

	// Poll for deployment status
	for i := 0; i < 60; i++ { // Max 5 minutes
		time.Sleep(5 * time.Second)

		d, _, err := client.Deployments.Get(ctx, deployment.ID)
		if err != nil {
			fmt.Printf("Error checking deployment: %v\n", err)
			break
		}

		fmt.Printf("   Status: %s\n", d.State)

		if d.State == sevalla.StatusSuccess {
			fmt.Println("\n✅ Deployment completed successfully!")
			if d.CompletedAt != nil {
				duration := d.CompletedAt.Sub(d.StartedAt)
				fmt.Printf("   Duration: %s\n", duration.Round(time.Second))
			}
			break
		} else if d.State == sevalla.StatusFailed {
			fmt.Printf("\n❌ Deployment failed: %s\n", d.ErrorMessage)
			break
		}
	}
}

func getApplicationLogs(ctx context.Context, client *sevalla.Client, appName string) {
	fmt.Printf("📜 Getting logs for application '%s'...\n", appName)

	// First, find the application by name
	apps, _, err := client.Applications.List(ctx, nil)
	if err != nil {
		handleError(err)
	}

	var app *sevalla.Application
	for _, a := range apps {
		if a.Name == appName {
			app = a
			break
		}
	}

	if app == nil {
		log.Fatalf("❌ Application '%s' not found", appName)
	}

	// Get last 100 lines of logs
	logs, resp, err := client.Applications.GetLogs(ctx, app.ID, 100)
	if err != nil {
		handleError(err)
	}

	fmt.Printf("\n✅ Logs retrieved! (Status: %d)\n\n", resp.StatusCode)
	fmt.Println("📜 Application Logs (last 100 lines):")
	fmt.Println(strings.Repeat("=", 61))
	fmt.Println(logs)
	fmt.Println(strings.Repeat("=", 61))
}

func handleError(err error) {
	if sevalla.IsNotFound(err) {
		log.Fatal("❌ Resource not found")
	} else if sevalla.IsUnauthorized(err) {
		log.Fatal("❌ Unauthorized: Check your API key")
	} else if sevalla.IsBadRequest(err) {
		log.Fatalf("❌ Bad request: %v", err)
	} else if sevalla.IsRateLimited(err) {
		log.Fatal("❌ Rate limited: Please wait and try again")
	} else if sevalla.IsServerError(err) {
		log.Fatalf("❌ Server error: %v", err)
	} else {
		log.Fatalf("❌ Error: %v", err)
	}
}
