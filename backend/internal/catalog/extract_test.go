package catalog

import "testing"

func TestExtract(t *testing.T) {
	raw := []byte(`{"services":[{"id":10,"name":"Deep Tissue Massage","category_name":"Massage","staffer_id":[44],"variants":[{"id":101,"duration":60,"price":90,"label":"60 minutes","staffer_id":[44]}]}],"staffers":[{"id":44,"name":"Sample Provider"}]}`)
	result, err := Extract(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Services) != 1 || result.Services[0].Variants[0].ID != 101 {
		t.Fatalf("unexpected services: %+v", result.Services)
	}
	if len(result.Providers) != 1 || result.Providers[0].BooksyStafferID != 44 {
		t.Fatalf("unexpected providers: %+v", result.Providers)
	}
}
