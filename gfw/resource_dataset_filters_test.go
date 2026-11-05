package gfw

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Filters populated by the API but absent from config must not plan a removal.
func TestDatasetFiltersOmittedInConfigNoDiff(t *testing.T) {
	r := resourceDataset()
	state := &terraform.InstanceState{
		ID: "ds",
		Attributes: map[string]string{
			"id":                                 "ds",
			"dataset_id":                         "ds",
			"name":                               "ds",
			"filters.#":                          "1",
			"filters.0.user_context_layers.#":    "1",
			"filters.0.user_context_layers.0.id": "name",
		},
	}
	cfg := terraform.NewResourceConfigRaw(map[string]interface{}{"dataset_id": "ds", "name": "ds"})
	diff, err := r.SimpleDiff(context.Background(), state, cfg, nil)
	if err != nil {
		t.Fatal(err)
	}
	if diff != nil {
		for k, a := range diff.Attributes {
			if len(k) >= 7 && k[:7] == "filters" {
				t.Fatalf("unexpected filters diff %s: %#v", k, a)
			}
		}
	}
}
