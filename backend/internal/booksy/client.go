package booksy

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

type Client struct {
	HTTP        *http.Client
	APIKey      string
	AppVersion  string
	Fingerprint string
}

type Request struct {
	BusinessID int
	VariantID  int
	StafferID  int
	StartDate  string
	EndDate    string
}

type upstreamSlot struct {
	Time string `json:"t"`
}

type upstreamDay struct {
	Date  string         `json:"date"`
	Slots []upstreamSlot `json:"slots"`
}

type upstreamResponse struct {
	TimeSlots []upstreamDay `json:"time_slots"`
}

type Date struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type Slot struct {
	Date   string `json:"date"`
	Value  string `json:"value"`
	Label  string `json:"label"`
	Period string `json:"period"`
}

type Availability struct {
	Dates []Date `json:"dates"`
	Slots []Slot `json:"slots"`
}

func NewFromEnv() (*Client, error) {
	client := &Client{
		HTTP:        &http.Client{Timeout: 12 * time.Second},
		APIKey:      strings.TrimSpace(os.Getenv("BOOKSY_API_KEY")),
		AppVersion:  strings.TrimSpace(os.Getenv("BOOKSY_APP_VERSION")),
		Fingerprint: strings.TrimSpace(os.Getenv("BOOKSY_FINGERPRINT")),
	}
	if client.APIKey == "" || client.AppVersion == "" || client.Fingerprint == "" {
		return nil, errors.New("Booksy availability credentials are not configured")
	}
	return client, nil
}

func (c *Client) Availability(ctx context.Context, input Request) (Availability, error) {
	if input.BusinessID <= 0 || input.VariantID <= 0 || input.StafferID == 0 {
		return Availability{}, errors.New("invalid Booksy identifiers")
	}
	body, err := json.Marshal(map[string]any{
		"subbookings": []map[string]any{{
			"service_variant_id": input.VariantID,
			"staffer_id":         input.StafferID,
			"combo_children":     []any{},
		}},
		"start_date": input.StartDate,
		"end_date":   input.EndDate,
	})
	if err != nil {
		return Availability{}, err
	}
	url := fmt.Sprintf("https://us.booksy.com/core/v2/customer_api/me/businesses/%d/appointments/time_slots", input.BusinessID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return Availability{}, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://booksy.com")
	req.Header.Set("Referer", "https://booksy.com/")
	req.Header.Set("X-Api-Key", c.APIKey)
	req.Header.Set("X-App-Version", c.AppVersion)
	req.Header.Set("X-Fingerprint", c.Fingerprint)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return Availability{}, fmt.Errorf("Booksy request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Availability{}, fmt.Errorf("Booksy returned HTTP %d", resp.StatusCode)
	}
	var payload upstreamResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return Availability{}, fmt.Errorf("decode Booksy response: %w", err)
	}
	return normalize(payload), nil
}

func normalize(payload upstreamResponse) Availability {
	seenDates := map[string]bool{}
	seenSlots := map[string]bool{}
	result := Availability{Dates: []Date{}, Slots: []Slot{}}
	for _, day := range payload.TimeSlots {
		parsedDate, err := time.Parse("2006-01-02", day.Date)
		if err != nil {
			continue
		}
		for _, raw := range day.Slots {
			parsedTime, err := time.Parse("15:04", raw.Time)
			if err != nil {
				continue
			}
			key := day.Date + "T" + raw.Time
			if seenSlots[key] {
				continue
			}
			seenSlots[key] = true
			if !seenDates[day.Date] {
				seenDates[day.Date] = true
				result.Dates = append(result.Dates, Date{Value: day.Date, Label: parsedDate.Format("Mon, Jan 2")})
			}
			period := "evening"
			if parsedTime.Hour() < 12 {
				period = "morning"
			} else if parsedTime.Hour() < 17 {
				period = "afternoon"
			}
			result.Slots = append(result.Slots, Slot{Date: day.Date, Value: raw.Time, Label: parsedTime.Format("3:04 PM"), Period: period})
		}
	}
	sort.Slice(result.Dates, func(i, j int) bool { return result.Dates[i].Value < result.Dates[j].Value })
	sort.Slice(result.Slots, func(i, j int) bool {
		if result.Slots[i].Date == result.Slots[j].Date {
			return result.Slots[i].Value < result.Slots[j].Value
		}
		return result.Slots[i].Date < result.Slots[j].Date
	})
	return result
}
