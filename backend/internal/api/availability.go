package api

import (
	"github.com/njsatish/dreamslayerz-portal/backend/internal/booksy"
	appconfig "github.com/njsatish/dreamslayerz-portal/backend/internal/config"
	"net/http"
	"strconv"
	"time"
)

type offering struct {
	Provider appconfig.Provider
	Service  appconfig.Service
	Variant  appconfig.Variant
}

func (h Handler) resolveOffering(providerID, serviceSlug string, variantID int) (offering, bool) {
	var provider appconfig.Provider
	foundProvider := false
	for _, p := range h.Config.Providers.Providers {
		if p.ID == providerID && p.Enabled && p.OnlineBooking {
			provider = p
			foundProvider = true
			break
		}
	}
	if !foundProvider {
		return offering{}, false
	}
	for _, s := range h.Config.Services.Services {
		if s.Slug != serviceSlug || !s.Enabled {
			continue
		}
		for _, v := range s.Variants {
			if v.ID != variantID {
				continue
			}
			for _, id := range v.ProviderIDs {
				if id == providerID {
					return offering{provider, s, v}, true
				}
			}
		}
	}
	return offering{}, false
}

func (h Handler) serveAvailability(w http.ResponseWriter, r *http.Request) {
	providerID := r.URL.Query().Get("provider")
	serviceSlug := r.URL.Query().Get("service")
	variantID, err := strconv.Atoi(r.URL.Query().Get("variant"))
	if err != nil || variantID <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "valid variant is required"})
		return
	}
	selected, ok := h.resolveOffering(providerID, serviceSlug, variantID)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "provider, service, and variant combination not found"})
		return
	}
	if !h.Config.Business.Booking.LiveAvailability || !selected.Variant.Availability.Enabled {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"error": "live availability is not enabled", "fallback": "widget"})
		return
	}
	client, err := booksy.NewFromEnv()
	if err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"error": "live availability is not configured", "fallback": "widget"})
		return
	}
	start := time.Now().In(time.FixedZone("business", -4*60*60))
	days := h.Config.Business.Booking.AvailabilityWindowDays
	if days < 1 || days > 31 {
		days = 14
	}
	end := start.AddDate(0, 0, days)
	result, err := client.Availability(r.Context(), booksy.Request{BusinessID: h.Config.Business.Booksy.BusinessID, VariantID: selected.Variant.ID, StafferID: selected.Provider.BooksyStafferID, StartDate: start.Format("2006-01-02"), EndDate: end.Format("2006-01-02")})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": "live availability is temporarily unavailable", "fallback": "widget"})
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, map[string]any{"providerId": providerID, "serviceSlug": serviceSlug, "variantId": variantID, "live": true, "dates": result.Dates, "slots": result.Slots, "notice": "Booksy confirms the appointment after you continue."})
}
