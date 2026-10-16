package booksy

import "testing"

func TestNormalize(t *testing.T) {
	got := normalize(upstreamResponse{TimeSlots: []upstreamDay{{Date: "2026-10-05", Slots: []upstreamSlot{{Time: "13:30"}, {Time: "09:00"}, {Time: "09:00"}}}}})
	if len(got.Dates) != 1 || len(got.Slots) != 2 {
		t.Fatalf("dates=%d slots=%d", len(got.Dates), len(got.Slots))
	}
	if got.Slots[0].Value != "09:00" || got.Slots[0].Label != "9:00 AM" || got.Slots[0].Period != "morning" {
		t.Fatalf("unexpected first slot: %+v", got.Slots[0])
	}
}
