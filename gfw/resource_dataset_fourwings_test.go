package gfw

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/globalfishingwatch.org/terraform-provider-gfw/gfw/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestSchemaToFourwingsV1ConfigLayoutVersion(t *testing.T) {
	config := schemaToFourwingsV1Config(map[string]interface{}{"layout_version": "v2"})
	if config.LayoutVersion != "v2" {
		t.Fatalf("expected LayoutVersion v2, got %q", config.LayoutVersion)
	}
	body, err := json.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), `"layoutVersion":"v2"`) {
		t.Fatalf("expected layoutVersion in JSON body, got %s", body)
	}
}

func TestFlattenFourwingsV1ConfigLayoutVersion(t *testing.T) {
	cases := map[string]string{
		"":   "v1", // existing datasets: API returns no layoutVersion
		"v1": "v1",
		"v2": "v2",
	}
	for apiValue, expected := range cases {
		got := flattenFourwingsV1Config(api.FourwingsV1Config{LayoutVersion: apiValue})["layout_version"]
		if got != expected {
			t.Fatalf("API value %q: expected layout_version %q, got %v", apiValue, expected, got)
		}
	}
}

func TestFourwingsV1LayoutVersionSchema(t *testing.T) {
	configuration := resourceDataset().Schema["configuration"].Elem.(*schema.Resource)
	fourwings := configuration.Schema["fourwings_v1"].Elem.(*schema.Resource)
	s, ok := fourwings.Schema["layout_version"]
	if !ok {
		t.Fatal("fourwings_v1.layout_version not defined")
	}
	if !s.Optional || s.Default != "v1" {
		t.Fatalf("expected optional with default v1, got optional=%v default=%v", s.Optional, s.Default)
	}
	for _, v := range []string{"v1", "v2"} {
		if _, errs := s.ValidateFunc(v, "layout_version"); len(errs) > 0 {
			t.Fatalf("layout_version %q should be valid: %v", v, errs)
		}
	}
	for _, v := range []string{"v3", "V2", ""} {
		if _, errs := s.ValidateFunc(v, "layout_version"); len(errs) == 0 {
			t.Fatalf("layout_version %q should be rejected", v)
		}
	}
}

// Omitting layout_version in config must match a state populated with the v1 default.
func TestFourwingsV1LayoutVersionDefaultNoDiff(t *testing.T) {
	r := resourceDataset()
	state := &terraform.InstanceState{
		ID: "ds",
		Attributes: map[string]string{
			"id":                             "ds",
			"dataset_id":                     "ds",
			"name":                           "ds",
			"configuration.#":                "1",
			"configuration.0.fourwings_v1.#": "1",
			"configuration.0.fourwings_v1.0.layout_version": "v1",
		},
	}
	cfg := terraform.NewResourceConfigRaw(map[string]interface{}{
		"dataset_id": "ds",
		"name":       "ds",
		"configuration": []interface{}{map[string]interface{}{
			"fourwings_v1": []interface{}{map[string]interface{}{"table": "t"}},
		}},
	})
	diff, err := r.SimpleDiff(context.Background(), state, cfg, nil)
	if err != nil {
		t.Fatal(err)
	}
	if diff != nil {
		if a, ok := diff.Attributes["configuration.0.fourwings_v1.0.layout_version"]; ok {
			t.Fatalf("unexpected layout_version diff: %#v", a)
		}
	}
}
