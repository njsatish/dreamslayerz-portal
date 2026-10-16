package catalog

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type Provider struct {
	ID              string   `json:"id"`
	BooksyStafferID int      `json:"booksyStafferId"`
	Name            string   `json:"name"`
	Title           string   `json:"title"`
	Enabled         bool     `json:"enabled"`
	OnlineBooking   bool     `json:"onlineBooking"`
	ServiceSlugs    []string `json:"serviceSlugs"`
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
type Result struct {
	Providers []Provider `json:"providers"`
	Services  []Service  `json:"services"`
	Warnings  []string   `json:"warnings"`
}

var nonSlug = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = nonSlug.ReplaceAllString(value, "-")
	return strings.Trim(value, "-")
}
func asMap(v any) (map[string]any, bool) { m, ok := v.(map[string]any); return m, ok }
func asSlice(v any) ([]any, bool)        { s, ok := v.([]any); return s, ok }
func asString(v any) string {
	if s, ok := v.(string); ok {
		return strings.TrimSpace(s)
	}
	return ""
}
func asInt(v any) int {
	switch n := v.(type) {
	case float64:
		return int(n)
	case json.Number:
		i, _ := strconv.Atoi(n.String())
		return i
	case int:
		return n
	case string:
		i, _ := strconv.Atoi(n)
		return i
	}
	return 0
}
func firstString(m map[string]any, keys ...string) string {
	for _, k := range keys {
		if value := asString(m[k]); value != "" {
			return value
		}
	}
	return ""
}
func firstInt(m map[string]any, keys ...string) int {
	for _, k := range keys {
		if value := asInt(m[k]); value > 0 {
			return value
		}
	}
	return 0
}
func titleName(value string) string {
	parts := strings.FieldsFunc(value, func(r rune) bool { return r == '-' || r == '_' })
	for i, p := range parts {
		r := []rune(p)
		if len(r) > 0 {
			r[0] = unicode.ToUpper(r[0])
			parts[i] = string(r)
		}
	}
	return strings.Join(parts, " ")
}

func staffIDs(v any) []int {
	seen := map[int]bool{}
	var out []int
	var walk func(any)
	walk = func(value any) {
		switch x := value.(type) {
		case []any:
			for _, item := range x {
				walk(item)
			}
		case map[string]any:
			for _, key := range []string{"id", "staffer_id", "staff_id"} {
				if id := asInt(x[key]); id > 0 && !seen[id] {
					seen[id] = true
					out = append(out, id)
				}
			}
		default:
			if id := asInt(x); id > 0 && !seen[id] {
				seen[id] = true
				out = append(out, id)
			}
		}
	}
	walk(v)
	sort.Ints(out)
	return out
}

func Extract(raw []byte) (Result, error) {
	decoder := json.NewDecoder(strings.NewReader(string(raw)))
	decoder.UseNumber()
	var root any
	if err := decoder.Decode(&root); err != nil {
		return Result{}, fmt.Errorf("parse Booksy JSON: %w", err)
	}
	providerNames := map[int]string{}
	providerServices := map[int]map[string]bool{}
	serviceByKey := map[string]Service{}
	var warnings []string
	var walk func(any, string)
	walk = func(value any, category string) {
		switch node := value.(type) {
		case []any:
			for _, item := range node {
				walk(item, category)
			}
		case map[string]any:
			currentCategory := category
			if c := firstString(node, "category_name", "category", "group_name"); c != "" {
				currentCategory = c
			}
			if id := firstInt(node, "id", "staffer_id"); id > 0 {
				if name := firstString(node, "display_name", "full_name", "name"); name != "" {
					if _, hasVariants := node["variants"]; !hasVariants {
						providerNames[id] = name
					}
				}
			}
			variantsRaw, hasVariants := asSlice(node["variants"])
			serviceID := firstInt(node, "id", "service_id")
			serviceName := firstString(node, "name", "service_name")
			if hasVariants && serviceID > 0 && serviceName != "" {
				slug := slugify(serviceName)
				if slug == "" {
					slug = fmt.Sprintf("service-%d", serviceID)
				}
				service := Service{ID: serviceID, Slug: slug, Name: serviceName, Category: currentCategory, Enabled: true, Variants: []Variant{}}
				for _, rawVariant := range variantsRaw {
					vm, ok := asMap(rawVariant)
					if !ok {
						continue
					}
					variantID := firstInt(vm, "id", "variant_id", "service_variant_id")
					if variantID <= 0 {
						continue
					}
					duration := firstInt(vm, "duration", "duration_minutes")
					price := vm["price"]
					priceLabel := firstString(vm, "price_label", "priceLabel")
					if priceLabel == "" && price != nil {
						switch p := price.(type) {
						case json.Number:
							priceLabel = "$" + p.String()
						case float64:
							priceLabel = fmt.Sprintf("$%.2f", p)
						}
					}
					if priceLabel == "" {
						priceLabel = "Price varies"
					}
					ids := staffIDs(vm["staffer_id"])
					if len(ids) == 0 {
						ids = staffIDs(node["staffer_id"])
					}
					providerIDs := []string{}
					for _, id := range ids {
						pid := fmt.Sprintf("provider-%d", id)
						providerIDs = append(providerIDs, pid)
						if providerServices[id] == nil {
							providerServices[id] = map[string]bool{}
						}
						providerServices[id][slug] = true
					}
					if len(providerIDs) == 0 {
						warnings = append(warnings, fmt.Sprintf("service %s variant %d has no staff assignment", serviceName, variantID))
					}
					service.Variants = append(service.Variants, Variant{ID: variantID, Label: firstString(vm, "label", "name"), DurationMinutes: duration, PriceLabel: priceLabel, ProviderIDs: providerIDs, Availability: Availability{Enabled: false, Route: "/availability/" + slug}})
				}
				if len(service.Variants) > 0 {
					serviceByKey[fmt.Sprintf("%d:%s", serviceID, slug)] = service
				}
			}
			for key, child := range node {
				if key == "variants" {
					continue
				}
				walk(child, currentCategory)
			}
		}
	}
	walk(root, "")
	services := make([]Service, 0, len(serviceByKey))
	for _, s := range serviceByKey {
		services = append(services, s)
	}
	sort.Slice(services, func(i, j int) bool { return services[i].Name < services[j].Name })
	providerIDs := map[int]bool{}
	for id := range providerNames {
		providerIDs[id] = true
	}
	for id := range providerServices {
		providerIDs[id] = true
	}
	providers := make([]Provider, 0, len(providerIDs))
	for id := range providerIDs {
		name := providerNames[id]
		if name == "" {
			name = titleName(fmt.Sprintf("provider-%d", id))
		}
		slugs := []string{}
		for slug := range providerServices[id] {
			slugs = append(slugs, slug)
		}
		sort.Strings(slugs)
		providers = append(providers, Provider{ID: fmt.Sprintf("provider-%d", id), BooksyStafferID: id, Name: name, Title: "Service provider", Enabled: true, OnlineBooking: true, ServiceSlugs: slugs})
	}
	sort.Slice(providers, func(i, j int) bool { return providers[i].Name < providers[j].Name })
	if len(services) == 0 {
		warnings = append(warnings, "no service objects with variants were discovered")
	}
	if len(providers) == 0 {
		warnings = append(warnings, "no provider assignments were discovered")
	}
	return Result{Providers: providers, Services: services, Warnings: warnings}, nil
}
