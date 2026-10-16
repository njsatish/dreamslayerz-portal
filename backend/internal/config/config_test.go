package config

import "testing"

func TestLoadDreamSlayerzConfig(t *testing.T) {
	cfg, err := Load("../../../config")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.ResolvedLayout != "single-provider" {
		t.Fatalf("layout=%s", cfg.ResolvedLayout)
	}
	if len(cfg.Services.Services) != 8 {
		t.Fatalf("services=%d", len(cfg.Services.Services))
	}
}
