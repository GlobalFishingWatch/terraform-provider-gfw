package gfw

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/globalfishingwatch.org/terraform-provider-gfw/gfw/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestProviderSchema(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatal(err)
	}
}

// Workspace sent to the API and read back must flatten to the same config, or plans drift
func TestReportWorkspaceRoundTrip(t *testing.T) {
	workspace := map[string]interface{}{
		"id":          "cmar-public",
		"name":        "CMAR",
		"description": "",
		"app":         "fishing-map",
		"category":    "fishing-activity",
		"state":       `{"reportCategory":"activity","sidebarOpen":true}`,
		"start_at":    "2025-01-19T23:41:59.000Z",
		"end_at":      "2025-04-19T23:41:59.000Z",
		"viewport":    []interface{}{map[string]interface{}{"zoom": 4.99, "latitude": 2.9, "longitude": -85.5}},
		"dataview_instances": []interface{}{map[string]interface{}{
			"id":              "basemap",
			"category":        "",
			"dataview_id":     "basemap",
			"config":          `{"visible":false}`,
			"datasets_config": []interface{}{`{"datasetId":"x","endpoint":"context-tiles","params":[]}`},
		}},
	}
	d := schema.TestResourceDataRaw(t, resourceReport().Schema, map[string]interface{}{
		"name": "CMAR", "description": "d", "area_ids": []interface{}{"1"}, "dataset_ids": []interface{}{"x"},
		"workspace": []interface{}{workspace},
	})
	sent, err := schemaToReportWorkspace(d)
	if err != nil {
		t.Fatal(err)
	}
	body, _ := json.Marshal(sent)
	var read api.ReportWorkspace
	if err := json.Unmarshal(body, &read); err != nil {
		t.Fatal(err)
	}
	flat, err := flattenReportWorkspace(read)
	if err != nil {
		t.Fatal(err)
	}
	if err := d.Set("workspace", []interface{}{flat}); err != nil {
		t.Fatal(err)
	}
	if got := d.Get("workspace").([]interface{})[0]; !reflect.DeepEqual(got, workspace) {
		t.Fatalf("round trip mismatch\n got: %#v\nwant: %#v", got, workspace)
	}
}
