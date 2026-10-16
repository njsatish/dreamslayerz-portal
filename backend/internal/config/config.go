package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Project struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}
type BusinessDetails struct {
	Industry string `json:"industry"`
	Timezone string `json:"timezone"`
	Currency string `json:"currency"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
}
type Website struct {
	ProductionURL string `json:"productionUrl"`
}
type Booksy struct {
	BusinessID int    `json:"businessId"`
	WidgetID   int    `json:"widgetId"`
	ProfileURL string `json:"profileUrl"`
	Country    string `json:"country"`
	Language   string `json:"language"`
}
type Booking struct {
	LayoutMode              string `json:"layoutMode"`
	SelectionOrder          string `json:"selectionOrder"`
	ProviderLabelSingular   string `json:"providerLabelSingular"`
	ProviderLabelPlural     string `json:"providerLabelPlural"`
	LiveAvailability        bool   `json:"liveAvailability"`
	SamePageWidget          bool   `json:"samePageWidget"`
	AvailabilityFailureMode string `json:"availabilityFailureMode"`
	AvailabilityWindowDays  int    `json:"availabilityWindowDays"`
}
type BusinessConfig struct {
	SchemaVersion int             `json:"schemaVersion"`
	Project       Project         `json:"project"`
	Business      BusinessDetails `json:"business"`
	Website       Website         `json:"website"`
	Booksy        Booksy          `json:"booksy"`
	Booking       Booking         `json:"booking"`
}
type Provider struct {
	ID              string   `json:"id"`
	BooksyStafferID int      `json:"booksyStafferId"`
	Name            string   `json:"name"`
	Title           string   `json:"title"`
	Enabled         bool     `json:"enabled"`
	OnlineBooking   bool     `json:"onlineBooking"`
	ServiceSlugs    []string `json:"serviceSlugs"`
}
type ProvidersConfig struct {
	SchemaVersion int        `json:"schemaVersion"`
	Providers     []Provider `json:"providers"`
}
type Availability struct {
	Enabled bool   `json:"enabled"`
	Route   string `json:"route"`
}
type Variant struct {
	ID              int          `json:"id"`
	Label           string       `json:"label"`
	DurationMinutes int          `json:"durationMinutes"`
	PriceLabel      string       `json:"priceLabel"`
	ProviderIDs     []string     `json:"providerIds"`
	Availability    Availability `json:"availability"`
}
type Service struct {
	ID       int       `json:"id"`
	Slug     string    `json:"slug"`
	Name     string    `json:"name"`
	Category string    `json:"category"`
	Enabled  bool      `json:"enabled"`
	Featured bool      `json:"featured"`
	Variants []Variant `json:"variants"`
}
type ServicesConfig struct {
	SchemaVersion int       `json:"schemaVersion"`
	Services      []Service `json:"services"`
}
type Theme struct {
	SchemaVersion int               `json:"schemaVersion"`
	Colors        map[string]string `json:"colors"`
	Hero          map[string]string `json:"hero"`
}
type Config struct {
	Business       BusinessConfig  `json:"business"`
	Providers      ProvidersConfig `json:"providers"`
	Services       ServicesConfig  `json:"services"`
	Theme          Theme           `json:"theme"`
	ResolvedLayout string          `json:"resolvedLayout"`
}

func Load(dir string) (Config, error) {
	var cfg Config
	for name, target := range map[string]any{"business.json": &cfg.Business, "providers.json": &cfg.Providers, "services.json": &cfg.Services, "theme.json": &cfg.Theme} {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return cfg, fmt.Errorf("read %s: %w", name, err)
		}
		if err := json.Unmarshal(data, target); err != nil {
			return cfg, fmt.Errorf("parse %s: %w", name, err)
		}
	}
	if err := cfg.Validate(); err != nil {
		return cfg, err
	}
	cfg.ResolvedLayout = cfg.resolveLayout()
	return cfg, nil
}

func (c Config) resolveLayout() string {
	if c.Business.Booking.LayoutMode != "auto" {
		return c.Business.Booking.LayoutMode
	}
	count := 0
	for _, p := range c.Providers.Providers {
		if p.Enabled && p.OnlineBooking && p.BooksyStafferID > 0 {
			count++
		}
	}
	if count > 1 {
		return "multi-provider"
	}
	return "single-provider"
}

func (c Config) Validate() error {
	var problems []string
	if c.Business.Project.Slug == "" || c.Business.Project.Name == "" {
		problems = append(problems, "project slug and name are required")
	}
	if c.Business.Booksy.BusinessID <= 0 || c.Business.Booksy.WidgetID <= 0 {
		problems = append(problems, "positive Booksy businessId and widgetId are required")
	}
	if c.Business.Booking.LayoutMode != "auto" && c.Business.Booking.LayoutMode != "single-provider" && c.Business.Booking.LayoutMode != "multi-provider" {
		problems = append(problems, "layoutMode must be auto, single-provider, or multi-provider")
	}
	providerIDs := map[string]bool{}
	eligible := 0
	for _, p := range c.Providers.Providers {
		if p.ID == "" || providerIDs[p.ID] {
			problems = append(problems, "provider IDs must be nonempty and unique")
		}
		providerIDs[p.ID] = true
		if p.Enabled && p.OnlineBooking && p.BooksyStafferID > 0 {
			eligible++
		}
	}
	if eligible == 0 {
		problems = append(problems, "at least one enabled online-booking provider is required")
	}
	slugs := map[string]bool{}
	variants := map[int]bool{}
	for _, s := range c.Services.Services {
		if !s.Enabled {
			continue
		}
		if s.Slug == "" || slugs[s.Slug] {
			problems = append(problems, "enabled service slugs must be nonempty and unique")
		}
		slugs[s.Slug] = true
		if len(s.Variants) == 0 {
			problems = append(problems, "service "+s.Slug+" has no variants")
		}
		for _, v := range s.Variants {
			if v.ID <= 0 || variants[v.ID] {
				problems = append(problems, fmt.Sprintf("variant IDs must be positive and unique: %d", v.ID))
			}
			variants[v.ID] = true
			if v.DurationMinutes < 1 || strings.TrimSpace(v.PriceLabel) == "" {
				problems = append(problems, "variant duration and price are required for "+s.Slug)
			}
			for _, id := range v.ProviderIDs {
				if !providerIDs[id] {
					problems = append(problems, "unknown provider "+id+" on "+s.Slug)
				}
			}
			if v.Availability.Enabled && v.Availability.Route == "" {
				problems = append(problems, "live availability route missing for "+s.Slug)
			}
		}
	}
	if len(problems) > 0 {
		return fmt.Errorf("configuration invalid: %s", strings.Join(problems, "; "))
	}
	return nil
}
