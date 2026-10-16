package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/njsatish/dreamslayerz-portal/backend/internal/catalog"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	businessID := flag.Int("business-id", 0, "Booksy business ID (required)")
	businessName := flag.String("business-name", "", "Business display name")
	projectSlug := flag.String("project-slug", "", "Project slug")
	profileURL := flag.String("profile-url", "", "Public Booksy profile URL")
	output := flag.String("output", "", "Output directory; defaults to generated/booksy-<id>")
	input := flag.String("input", "", "Parse a previously saved Booksy business JSON file instead of making a request")
	flag.Parse()
	if *businessID <= 0 {
		fatal("--business-id must be a positive integer")
	}
	if *output == "" {
		*output = filepath.Join("generated", fmt.Sprintf("booksy-%d", *businessID))
	}
	var raw []byte
	var err error
	if *input != "" {
		raw, err = os.ReadFile(*input)
	} else {
		raw, err = fetchBusiness(context.Background(), *businessID)
	}
	if err != nil {
		fatal(err.Error())
	}
	result, err := catalog.Extract(raw)
	if err != nil {
		fatal(err.Error())
	}
	if err := os.MkdirAll(*output, 0755); err != nil {
		fatal(err.Error())
	}
	mustWrite(filepath.Join(*output, "raw-business.json"), raw)
	mustWriteJSON(filepath.Join(*output, "providers.candidate.json"), map[string]any{"schemaVersion": 1, "providers": result.Providers})
	mustWriteJSON(filepath.Join(*output, "services.candidate.json"), map[string]any{"schemaVersion": 1, "services": result.Services})
	name := strings.TrimSpace(*businessName)
	if name == "" {
		name = fmt.Sprintf("Booksy Business %d", *businessID)
	}
	slug := strings.TrimSpace(*projectSlug)
	if slug == "" {
		slug = fmt.Sprintf("booksy-%d", *businessID)
	}
	mustWriteJSON(filepath.Join(*output, "business.candidate.json"), map[string]any{"schemaVersion": 1, "project": map[string]any{"slug": slug, "name": name}, "booksy": map[string]any{"businessId": *businessID, "widgetId": *businessID, "profileUrl": *profileURL, "country": "us", "language": "en"}, "booking": map[string]any{"layoutMode": "auto", "selectionOrder": "provider-first", "providerLabelSingular": "Provider", "providerLabelPlural": "Providers", "liveAvailability": false, "samePageWidget": true, "availabilityFailureMode": "widget", "availabilityWindowDays": 14}})
	mustWriteJSON(filepath.Join(*output, "capture-report.json"), map[string]any{"businessId": *businessID, "capturedAt": time.Now().UTC().Format(time.RFC3339), "providers": len(result.Providers), "services": len(result.Services), "warnings": result.Warnings, "status": "candidate-review-required"})
	fmt.Printf("Capture complete: %s\nProviders: %d\nServices: %d\n", *output, len(result.Providers), len(result.Services))
	for _, warning := range result.Warnings {
		fmt.Printf("WARNING: %s\n", warning)
	}
	fmt.Println("No live config was overwritten. Review candidate files before applying them.")
}

func fetchBusiness(ctx context.Context, businessID int) ([]byte, error) {
	apiKey := strings.TrimSpace(os.Getenv("BOOKSY_API_KEY"))
	appVersion := strings.TrimSpace(os.Getenv("BOOKSY_APP_VERSION"))
	fingerprint := strings.TrimSpace(os.Getenv("BOOKSY_FINGERPRINT"))
	if apiKey == "" || appVersion == "" || fingerprint == "" {
		return nil, fmt.Errorf("BOOKSY_API_KEY, BOOKSY_APP_VERSION, and BOOKSY_FINGERPRINT are required for live capture")
	}
	template := strings.TrimSpace(os.Getenv("BOOKSY_CATALOG_URL_TEMPLATE"))
	if template == "" {
		template = "https://us.booksy.com/core/v2/customer_api/me/businesses/%d"
	}
	url := fmt.Sprintf(template, businessID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Origin", "https://booksy.com")
	req.Header.Set("Referer", "https://booksy.com/")
	req.Header.Set("X-Api-Key", apiKey)
	req.Header.Set("X-App-Version", appVersion)
	req.Header.Set("X-Fingerprint", fingerprint)
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 25<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Booksy catalog request returned HTTP %d", resp.StatusCode)
	}
	return body, nil
}
func mustWrite(path string, data []byte) {
	if err := os.WriteFile(path, data, 0644); err != nil {
		fatal(err.Error())
	}
}
func mustWriteJSON(path string, value any) {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fatal(err.Error())
	}
	mustWrite(path, append(data, '\n'))
}
func fatal(message string) { fmt.Fprintln(os.Stderr, "ERROR:", message); os.Exit(1) }
