package gfw

import (
	"context"
	"time"

	"github.com/globalfishingwatch.org/terraform-provider-gfw/gfw/api"
	"github.com/globalfishingwatch.org/terraform-provider-gfw/gfw/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var DATASET_TYPES []string = []string{
	"4wings:v1",
	"bulk-download:v1",
	"context-layer:v1",
	"data-download:v1",
	"events:v1",
	"insights:v1",
	"pm-tiles:v1",
	"ports:v1",
	"temporal-context-layer:v1",
	"thumbnails:v1",
	"tracks:v1",
	"user-context-layer:v1",
	"user-tracks:v1",
	"vessels:v1",
}

var DATASET_CATEGORIES []string = []string{
	"activity",
	"context-layer",
	"context",
	"detections",
	"environment",
	"event",
	"gap",
	"vessel",
}

var DATASET_SUBCATEGORIES []string = []string{
	"ais_false",
	"animal",
	"chlorophyl",
	"currents",
	"encounter",
	"fishing",
	"gap_start",
	"gap",
	"gaps",
	"info",
	"insight",
	"loitering",
	"nitrate",
	"oxygen",
	"ph",
	"phosphate",
	"port_visit",
	"ports",
	"presence",
	"salinity",
	"sar",
	"sentinel-2",
	"track",
	"user",
	"viirs",
	"water-temperature",
	"waves",
	"winds",
}

var DATASET_UNITS []string = []string{
	"days",
	"detections",
	"habitat suitability",
	"hours",
	"mg/m^3",
	"NA",
	"ºC",
	"probability",
	"PSU",
	"TBD",
	"unit",
}

var DATASET_STATUSES []string = []string{
	"deprecated",
	"done",
	"error",
	"importing",
}

var DATASET_CONFIGURATION_GEOMETRY_TYPES []string = []string{"tracks", "polygons", "points"}

var DATASET_CONTEXT_LAYER_FORMATS []string = []string{"csv", "geojson", "pmtile"}
var DATASET_BULK_DOWNLOAD_FORMATS []string = []string{"CSV", "JSON"}
var DATASET_4WINGS_INTERVALS []string = []string{"HOUR", "DAY", "MONTH", "YEAR"}
var DATASET_FRONTEND_FORMATS []string = []string{"GeoJSON", "Shapefile", "CSV", "KML"}
var DATASET_4WINGS_REPORT_GROUPINGS []string = []string{"id", "mmsi", "geartype", "flag", "flagAndGearType"}
var DATASET_SOURCE_TYPES []string = []string{"gcs", "bigquery", "clickhouse"}

func filterConfigSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"label": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"type": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"required": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"array": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"enum": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"format": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"max_length": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"min_length": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"max": {
			Type:     schema.TypeFloat,
			Optional: true,
		},
		"min": {
			Type:     schema.TypeFloat,
			Optional: true,
		},
		"single_selection": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"operation": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

func resourceDataset() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatasetCreate,
		ReadContext:   resourceDatasetRead,
		UpdateContext: resourceDatasetUpdate,
		DeleteContext: resourceDatasetDelete,
		Schema: map[string]*schema.Schema{
			"dataset_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(DATASET_TYPES, false),
			},
			"alias": {
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Type:     schema.TypeList,
				Optional: true,
			},
			"description": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"start_date": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: utils.IsISOTime,
			},
			"end_date": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: utils.IsISOTime,
			},
			"unit": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"category": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(DATASET_CATEGORIES, false),
			},
			"subcategory": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(DATASET_SUBCATEGORIES, false),
			},
			"source": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(DATASET_STATUSES, false),
			},
			"related_datasets": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(DATASET_TYPES, false),
						},
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_supported_versions": {
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Type:     schema.TypeList,
							Optional: true,
						},
						"context_layer_v1": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"import_logs": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"srid": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"format": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(DATASET_CONTEXT_LAYER_FORMATS, false),
									},
									"fields": {
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Type:     schema.TypeList,
										Optional: true,
									},
									"file_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"id_property": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"user_context_layer_v1": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"table": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"import_logs": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"srid": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"format": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(DATASET_CONTEXT_LAYER_FORMATS, false),
									},
									"fields": {
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Type:     schema.TypeList,
										Optional: true,
									},
									"file_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"id_property": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"value_property_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"temporal_context_layer_v1": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dataset": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"project": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"source": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"table": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"user_tracks_v1": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"file_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"id_property": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"pm_tiles_v1": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"file_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"id_property": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"events_v1": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"table": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"dataset": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"project": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"function": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ttl": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"max_zoom": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  12,
									},
									"source": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(DATASET_SOURCE_TYPES, false),
									},
								},
							},
						},
						"fourwings_v1": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"report_groupings": {
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringInSlice(DATASET_4WINGS_REPORT_GROUPINGS, false),
										},
										Type:     schema.TypeList,
										Optional: true,
									},
									"table": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"dataset": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"max_zoom": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  12,
									},
									"project": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"function": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"intervals": {
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringInSlice(DATASET_4WINGS_INTERVALS, false),
										},
										Type:     schema.TypeList,
										Optional: true,
									},
									"ttl": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"max": {
										Type:     schema.TypeFloat,
										Optional: true,
									},
									"min": {
										Type:     schema.TypeFloat,
										Optional: true,
									},
									"tile_scale": {
										Type:     schema.TypeFloat,
										Optional: true,
									},
									"tile_offset": {
										Type:     schema.TypeFloat,
										Optional: true,
									},
									"internal_scale": {
										Type:     schema.TypeFloat,
										Optional: true,
									},
									"internal_offset": {
										Type:     schema.TypeFloat,
										Optional: true,
									},
									"gee_band": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"gee_images": {
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Type:     schema.TypeList,
										Optional: true,
									},
									"interaction_columns": {
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Type:     schema.TypeList,
										Optional: true,
									},
									"interaction_group_columns": {
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Type:     schema.TypeList,
										Optional: true,
									},
									"temporal_aggregation": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"source": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"bucket": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"folder": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"tracks_v1": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"folder": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"database_instance": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"table": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"frontend": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"max_zoom": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  12,
									},
									"translate": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"max": {
										Type:     schema.TypeFloat,
										Optional: true,
									},
									"min": {
										Type:     schema.TypeFloat,
										Optional: true,
									},
									"disable_interaction": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"latitude": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"longitude": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"start_time": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"end_time": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"timestamp": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"geometry_type": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(DATASET_CONFIGURATION_GEOMETRY_TYPES, false),
									},
									"source_format": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(DATASET_FRONTEND_FORMATS, false),
									},
									"time_filter_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"value_properties": {
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Type:     schema.TypeList,
										Optional: true,
									},
									"polygon_color": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"point_size": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"line_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"segment_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"vessels_v1": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"index": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"index_boost": {
										Type:     schema.TypeFloat,
										Optional: true,
									},
								},
							},
						},
						"insights_v1": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sources": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeString,
													Required: true,
												},
												"type": {
													Type:     schema.TypeString,
													Required: true,
												},
												"insight": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
								},
							},
						},
						"bulk_download_v1": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"gcs_uri": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"format": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(DATASET_BULK_DOWNLOAD_FORMATS, false),
									},
									"compressed": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"data_download_v1": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"email_groups": {
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Type:     schema.TypeList,
										Optional: true,
									},
									"gcs_folder": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"doi": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"concept_doi": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"thumbnails_v1": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"extensions": {
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Type:     schema.TypeList,
										Optional: true,
									},
									"bucket": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"folder": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"scale": {
										Type:     schema.TypeFloat,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"filters": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fourwings": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: filterConfigSchema(),
							},
						},
						"events": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: filterConfigSchema(),
							},
						},
						"vessels": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: filterConfigSchema(),
							},
						},
						"tracks": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: filterConfigSchema(),
							},
						},
					},
				},
			},
			"documentation": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"enable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"status": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"queries": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"provider": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceDatasetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	c := m.(*api.GFWClient)
	var diags diag.Diagnostics

	id := d.Get("dataset_id").(string)
	dataset, err := schemaToDataset(d)
	if err != nil {
		return diag.FromErr(err)
	}
	dataset.ID = id
	datasetCreated, err := c.CreateDataset(dataset)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(datasetCreated.ID)
	resourceDatasetRead(ctx, d, m)

	return diags
}

func resourceDatasetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	datasetId := d.Id()
	c := m.(*api.GFWClient)
	dataset, err := c.GetDataset(datasetId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", dataset.Name)
	d.Set("description", dataset.Description)
	d.Set("created_at", dataset.CreatedAt)
	d.Set("type", dataset.Type)
	d.Set("alias", dataset.Alias)
	d.Set("start_date", dataset.StartDate)
	if d.Get("end_date") != nil && d.Get("end_date").(string) != "" {
		d.Set("end_date", dataset.EndDate)
	}
	d.Set("unit", dataset.Unit)
	d.Set("category", dataset.Category)
	d.Set("subcategory", dataset.Subcategory)
	if d.Get("status") != nil && d.Get("status").(string) != "" {
		d.Set("status", dataset.Status)
	}
	d.Set("source", dataset.Source)
	d.Set("type", dataset.Type)

	if dataset.Filters != nil {
		filters := flattenDatasetFilters(*dataset.Filters)
		if err := d.Set("filters", []interface{}{filters}); err != nil {
			return diag.FromErr(err)
		}
	}

	if dataset.Configuration != nil {
		configuration := flattenDatasetConfiguration(*dataset.Configuration)
		if err := d.Set("configuration", []interface{}{configuration}); err != nil {
			return diag.FromErr(err)
		}
	}

	if dataset.Documentation != nil {
		documentation := flattenDatasetDocumentation(*dataset.Documentation)
		if err := d.Set("documentation", []interface{}{documentation}); err != nil {
			return diag.FromErr(err)
		}
	}

	relatedDatasets := flattenRelatedDatasets(dataset.RelatedDatasets)
	if err := d.Set("related_datasets", relatedDatasets); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceDatasetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	dataset, err := schemaToDataset(d)
	if err != nil {
		return diag.FromErr(err)
	}
	datasetId := d.Id()
	c := m.(*api.GFWClient)
	err = c.UpdateDataset(datasetId, dataset)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceDatasetRead(ctx, d, m)
}

func schemaToDataset(d *schema.ResourceData) (api.CreateDataset, error) {
	dataset := api.CreateDataset{}
	dataset.Name = d.Get("name").(string)

	dataset.Description = d.Get("description").(string)

	dataset.Type = d.Get("type").(string)

	dataset.Unit = d.Get("unit").(string)

	dataset.Category = d.Get("category").(string)

	dataset.Subcategory = d.Get("subcategory").(string)

	dataset.Status = d.Get("status").(string)

	dataset.Source = d.Get("source").(string)

	dataset.StartDate = d.Get("start_date").(string)

	dataset.EndDate = d.Get("end_date").(string)

	if d.Get("alias") != nil {
		dataset.Alias = utils.ConvertArrayInterfaceToArrayString(d.Get("alias").([]interface{}))
	}
	if d.Get("filters") != nil {
		filtersList := d.Get("filters").([]interface{})
		if len(filtersList) > 0 {
			filters := schemaToDatasetFilters(filtersList[0].(map[string]interface{}))
			dataset.Filters = &filters
		}
	}
	if d.Get("configuration") != nil {
		configuration := d.Get("configuration").([]interface{})
		if len(configuration) > 0 {
			config := schemaToDatasetConfiguration(configuration[0].(map[string]interface{}), d.Get("name").(string))
			dataset.Configuration = &config
		}
	}
	if d.Get("documentation") != nil {
		documentationList := d.Get("documentation").([]interface{})
		if len(documentationList) > 0 {
			documentation := schemaToDatasetDocumentation(documentationList[0].(map[string]interface{}))
			dataset.Documentation = &documentation
		}
	}
	if d.Get("related_datasets") != nil {
		relatedDatasets := schemaToRelatedDatasets(d.Get("related_datasets").([]interface{}))
		dataset.RelatedDatasets = relatedDatasets
	}

	return dataset, nil
}

func schemaToDatasetConfiguration(schema map[string]interface{}, name string) api.DatasetConfiguration {
	config := api.DatasetConfiguration{}

	// API Supported Versions
	if val, ok := schema["api_supported_versions"]; ok {
		config.ApiSupportedVersions = utils.ConvertArrayInterfaceToArrayString(val.([]interface{}))
	}

	// Context Layer V1
	if val, ok := schema["context_layer_v1"]; ok {
		configArray := val.([]interface{})
		if len(configArray) > 0 {
			contextConfig := schemaToContextLayerV1Config(configArray[0].(map[string]interface{}))
			config.ContextLayerV1 = &contextConfig
		}
	}

	// User Context Layer V1
	if val, ok := schema["user_context_layer_v1"]; ok {
		configArray := val.([]interface{})
		if len(configArray) > 0 {
			userContextConfig := schemaToUserContextLayerV1Config(configArray[0].(map[string]interface{}))
			config.UserContextLayerV1 = &userContextConfig
		}
	}

	// Temporal Context Layer V1
	if val, ok := schema["temporal_context_layer_v1"]; ok {
		configArray := val.([]interface{})
		if len(configArray) > 0 {
			temporalContextConfig := schemaToTemporalContextLayerV1Config(configArray[0].(map[string]interface{}))
			config.TemporalContextLayerV1 = &temporalContextConfig
		}
	}

	// User Tracks V1
	if val, ok := schema["user_tracks_v1"]; ok {
		configArray := val.([]interface{})
		if len(configArray) > 0 {
			userTracksConfig := schemaToUserTracksV1Config(configArray[0].(map[string]interface{}))
			config.UserTracksV1 = &userTracksConfig
		}
	}

	// PM Tiles V1
	if val, ok := schema["pm_tiles_v1"]; ok {
		configArray := val.([]interface{})
		if len(configArray) > 0 {
			pmTilesConfig := schemaToPmTilesV1Config(configArray[0].(map[string]interface{}))
			config.PmTilesV1 = &pmTilesConfig
		}
	}

	// Events V1
	if val, ok := schema["events_v1"]; ok {
		configArray := val.([]interface{})
		if len(configArray) > 0 {
			eventsConfig := schemaToEventsV1Config(configArray[0].(map[string]interface{}))
			config.EventsV1 = &eventsConfig
		}
	}

	// Fourwings V1
	if val, ok := schema["fourwings_v1"]; ok {
		configArray := val.([]interface{})
		if len(configArray) > 0 {
			fourwingsConfig := schemaToFourwingsV1Config(configArray[0].(map[string]interface{}))
			config.FourwingsV1 = &fourwingsConfig
		}
	}

	// Tracks V1
	if val, ok := schema["tracks_v1"]; ok {
		configArray := val.([]interface{})
		if len(configArray) > 0 {
			tracksConfig := schemaToTracksV1Config(configArray[0].(map[string]interface{}))
			config.TracksV1 = &tracksConfig
		}
	}

	// Frontend
	if val, ok := schema["frontend"]; ok {
		configArray := val.([]interface{})
		if len(configArray) > 0 {
			frontendConfig := schemaToFrontendConfig(configArray[0].(map[string]interface{}))
			config.Frontend = &frontendConfig
		}
	}

	// Vessels V1
	if val, ok := schema["vessels_v1"]; ok {
		configArray := val.([]interface{})
		if len(configArray) > 0 {
			vesselsConfig := schemaToVesselsV1Config(configArray[0].(map[string]interface{}))
			config.VesselsV1 = &vesselsConfig
		}
	}

	// Insights V1
	if val, ok := schema["insights_v1"]; ok {
		configArray := val.([]interface{})
		if len(configArray) > 0 {
			insightsConfig := schemaToInsightsV1Config(configArray[0].(map[string]interface{}))
			config.InsightsV1 = &insightsConfig
		}
	}

	// Bulk Download V1
	if val, ok := schema["bulk_download_v1"]; ok {
		configArray := val.([]interface{})
		if len(configArray) > 0 {
			bulkDownloadConfig := schemaToBulkDownloadV1Config(configArray[0].(map[string]interface{}))
			config.BulkDownloadV1 = &bulkDownloadConfig
		}
	}

	// Data Download V1
	if val, ok := schema["data_download_v1"]; ok {
		configArray := val.([]interface{})
		if len(configArray) > 0 {
			dataDownloadConfig := schemaToDataDownloadV1Config(configArray[0].(map[string]interface{}))
			config.DataDownloadV1 = &dataDownloadConfig
		}
	}

	// Thumbnails V1
	if val, ok := schema["thumbnails_v1"]; ok {
		configArray := val.([]interface{})
		if len(configArray) > 0 {
			thumbnailsConfig := schemaToThumbnailsV1Config(configArray[0].(map[string]interface{}))
			config.ThumbnailsV1 = &thumbnailsConfig
		}
	}

	return config
}

// Helper functions for nested configurations

func schemaToContextLayerV1Config(schema map[string]interface{}) api.ContextLayerV1Config {
	config := api.ContextLayerV1Config{}
	if val, ok := schema["import_logs"]; ok {
		config.ImportLogs = val.(string)
	}
	if val, ok := schema["srid"]; ok {
		config.Srid = val.(string)
	}
	if val, ok := schema["format"]; ok {
		config.Format = val.(string)
	}
	if val, ok := schema["fields"]; ok {
		config.Fields = utils.ConvertArrayInterfaceToArrayString(val.([]interface{}))
	}
	if val, ok := schema["file_path"]; ok {
		config.FilePath = val.(string)
	}
	if val, ok := schema["id_property"]; ok {
		config.IDProperty = val.(string)
	}
	return config
}

func schemaToUserContextLayerV1Config(schema map[string]interface{}) api.UserContextLayerV1Config {
	config := api.UserContextLayerV1Config{}
	if val, ok := schema["table"]; ok {
		config.Table = val.(string)
	}
	if val, ok := schema["import_logs"]; ok {
		config.ImportLogs = val.(string)
	}
	if val, ok := schema["srid"]; ok {
		config.Srid = val.(string)
	}
	if val, ok := schema["format"]; ok {
		config.Format = val.(string)
	}
	if val, ok := schema["fields"]; ok {
		config.Fields = utils.ConvertArrayInterfaceToArrayString(val.([]interface{}))
	}
	if val, ok := schema["file_path"]; ok {
		config.FilePath = val.(string)
	}
	if val, ok := schema["id_property"]; ok {
		config.IDProperty = val.(string)
	}
	if val, ok := schema["value_property_id"]; ok {
		config.ValuePropertyID = val.(string)
	}
	return config
}

func schemaToTemporalContextLayerV1Config(schema map[string]interface{}) api.TemporalContextLayerV1Config {
	config := api.TemporalContextLayerV1Config{}
	if val, ok := schema["dataset"]; ok {
		config.Dataset = val.(string)
	}
	if val, ok := schema["project"]; ok {
		config.Project = val.(string)
	}
	if val, ok := schema["source"]; ok {
		config.Source = val.(string)
	}
	if val, ok := schema["table"]; ok {
		config.Table = val.(string)
	}
	return config
}

func schemaToUserTracksV1Config(schema map[string]interface{}) api.UserTracksV1Config {
	config := api.UserTracksV1Config{}
	if val, ok := schema["file_path"]; ok {
		config.FilePath = val.(string)
	}
	if val, ok := schema["id_property"]; ok {
		config.IDProperty = val.(string)
	}
	return config
}

func schemaToPmTilesV1Config(schema map[string]interface{}) api.PmTilesV1Config {
	config := api.PmTilesV1Config{}
	if val, ok := schema["file_path"]; ok {
		config.FilePath = val.(string)
	}
	if val, ok := schema["id_property"]; ok {
		config.IDProperty = val.(string)
	}
	return config
}

func schemaToEventsV1Config(schema map[string]interface{}) api.EventsV1Config {
	config := api.EventsV1Config{}
	if val, ok := schema["table"]; ok {
		config.Table = val.(string)
	}
	if val, ok := schema["dataset"]; ok {
		config.Dataset = val.(string)
	}
	if val, ok := schema["project"]; ok {
		config.Project = val.(string)
	}
	if val, ok := schema["function"]; ok {
		config.Function = val.(string)
	}
	if val, ok := schema["ttl"]; ok {
		config.TTL = val.(int)
	}
	if val, ok := schema["max_zoom"]; ok {
		config.MaxZoom = val.(int)
	}
	if val, ok := schema["source"]; ok {
		config.Source = val.(string)
	}
	return config
}

func schemaToFourwingsV1Config(schema map[string]interface{}) api.FourwingsV1Config {
	config := api.FourwingsV1Config{}
	if val, ok := schema["report_groupings"]; ok {
		config.ReportGroupings = utils.ConvertArrayInterfaceToArrayString(val.([]interface{}))
	}
	if val, ok := schema["table"]; ok {
		config.Table = val.(string)
	}
	if val, ok := schema["dataset"]; ok {
		config.Dataset = val.(string)
	}
	if val, ok := schema["max_zoom"]; ok {
		config.MaxZoom = val.(int)
	}
	if val, ok := schema["project"]; ok {
		config.Project = val.(string)
	}
	if val, ok := schema["function"]; ok {
		config.Function = val.(string)
	}
	if val, ok := schema["intervals"]; ok {
		config.Intervals = utils.ConvertArrayInterfaceToArrayString(val.([]interface{}))
	}
	if val, ok := schema["ttl"]; ok {
		config.TTL = val.(int)
	}
	if val, ok := schema["max"]; ok {
		config.Max = val.(float64)
	}
	if val, ok := schema["min"]; ok {
		config.Min = val.(float64)
	}
	if val, ok := schema["tile_scale"]; ok {
		config.TileScale = val.(float64)
	}
	if val, ok := schema["tile_offset"]; ok {
		config.TileOffset = val.(float64)
	}
	if val, ok := schema["internal_scale"]; ok {
		config.InternalScale = val.(float64)
	}
	if val, ok := schema["internal_offset"]; ok {
		config.InternalOffset = val.(float64)
	}
	if val, ok := schema["gee_band"]; ok {
		config.GeeBand = val.(string)
	}
	if val, ok := schema["gee_images"]; ok {
		config.GeeImages = utils.ConvertArrayInterfaceToArrayString(val.([]interface{}))
	}
	if val, ok := schema["interaction_columns"]; ok {
		config.InteractionColumns = utils.ConvertArrayInterfaceToArrayString(val.([]interface{}))
	}
	if val, ok := schema["interaction_group_columns"]; ok {
		config.InteractionGroupColumns = utils.ConvertArrayInterfaceToArrayString(val.([]interface{}))
	}
	if val, ok := schema["temporal_aggregation"]; ok {
		config.TemporalAggregation = val.(bool)
	}
	if val, ok := schema["source"]; ok {
		config.Source = val.(string)
	}
	if val, ok := schema["bucket"]; ok {
		config.Bucket = val.(string)
	}
	if val, ok := schema["folder"]; ok {
		config.Folder = val.(string)
	}
	return config
}

func schemaToTracksV1Config(schema map[string]interface{}) api.TracksV1Config {
	config := api.TracksV1Config{}
	if val, ok := schema["bucket"]; ok {
		config.Bucket = val.(string)
	}
	if val, ok := schema["folder"]; ok {
		config.Folder = val.(string)
	}
	if val, ok := schema["database_instance"]; ok {
		config.DatabaseInstance = val.(string)
	}
	if val, ok := schema["table"]; ok {
		config.Table = val.(string)
	}
	return config
}

func schemaToFrontendConfig(schema map[string]interface{}) api.FrontendConfig {
	config := api.FrontendConfig{}
	if val, ok := schema["max_zoom"]; ok {
		config.MaxZoom = val.(int)
	}
	if val, ok := schema["translate"]; ok {
		config.Translate = val.(bool)
	}
	if val, ok := schema["max"]; ok {
		config.Max = val.(float64)
	}
	if val, ok := schema["min"]; ok {
		config.Min = val.(float64)
	}
	if val, ok := schema["disable_interaction"]; ok {
		config.DisableInteraction = val.(bool)
	}
	if val, ok := schema["latitude"]; ok {
		config.Latitude = val.(string)
	}
	if val, ok := schema["longitude"]; ok {
		config.Longitude = val.(string)
	}
	if val, ok := schema["start_time"]; ok {
		config.StartTime = val.(string)
	}
	if val, ok := schema["end_time"]; ok {
		config.EndTime = val.(string)
	}
	if val, ok := schema["timestamp"]; ok {
		config.Timestamp = val.(string)
	}
	if val, ok := schema["geometry_type"]; ok {
		config.GeometryType = val.(string)
	}
	if val, ok := schema["source_format"]; ok {
		config.SourceFormat = val.(string)
	}
	if val, ok := schema["time_filter_type"]; ok {
		config.TimeFilterType = val.(string)
	}
	if val, ok := schema["value_properties"]; ok {
		config.ValueProperties = utils.ConvertArrayInterfaceToArrayString(val.([]interface{}))
	}
	if val, ok := schema["polygon_color"]; ok {
		config.PolygonColor = val.(string)
	}
	if val, ok := schema["point_size"]; ok {
		config.PointSize = val.(string)
	}
	if val, ok := schema["line_id"]; ok {
		config.LineID = val.(string)
	}
	if val, ok := schema["segment_id"]; ok {
		config.SegmentID = val.(string)
	}
	return config
}

func schemaToVesselsV1Config(schema map[string]interface{}) api.VesselsV1Config {
	config := api.VesselsV1Config{}
	if val, ok := schema["index"]; ok {
		config.Index = val.(string)
	}
	if val, ok := schema["index_boost"]; ok {
		config.IndexBoost = val.(float64)
	}
	return config
}

func schemaToInsightsV1Config(schema map[string]interface{}) api.InsightsV1Config {
	config := api.InsightsV1Config{}
	if val, ok := schema["sources"]; ok {
		sourcesArray := val.([]interface{})
		if len(sourcesArray) > 0 {
			array := make([]api.InsightSources, len(sourcesArray))
			for i, source := range sourcesArray {
				array[i] = schemaToDatasetInsightSource(source.(map[string]interface{}))
			}
			config.Sources = array
		}
	}
	return config
}

func schemaToBulkDownloadV1Config(schema map[string]interface{}) api.BulkDownloadV1Config {
	config := api.BulkDownloadV1Config{}
	if val, ok := schema["gcs_uri"]; ok {
		config.GcsUri = val.(string)
	}
	if val, ok := schema["path"]; ok {
		config.Path = val.(string)
	}
	if val, ok := schema["format"]; ok {
		config.Format = val.(string)
	}
	if val, ok := schema["compressed"]; ok {
		config.Compressed = val.(bool)
	}
	return config
}

func schemaToDataDownloadV1Config(schema map[string]interface{}) api.DataDownloadV1Config {
	config := api.DataDownloadV1Config{}
	if val, ok := schema["email_groups"]; ok {
		config.EmailGroups = utils.ConvertArrayInterfaceToArrayString(val.([]interface{}))
	}
	if val, ok := schema["gcs_folder"]; ok {
		config.GcsFolder = val.(string)
	}
	if val, ok := schema["doi"]; ok {
		config.Doi = val.(string)
	}
	if val, ok := schema["concept_doi"]; ok {
		config.ConceptDOI = val.(int)
	}
	return config
}

func schemaToThumbnailsV1Config(schema map[string]interface{}) api.ThumbnailsV1Config {
	config := api.ThumbnailsV1Config{}
	if val, ok := schema["extensions"]; ok {
		config.Extensions = utils.ConvertArrayInterfaceToArrayString(val.([]interface{}))
	}
	if val, ok := schema["bucket"]; ok {
		config.Bucket = val.(string)
	}
	if val, ok := schema["folder"]; ok {
		config.Folder = val.(string)
	}
	if val, ok := schema["scale"]; ok {
		config.Scale = val.(float64)
	}
	return config
}

func schemaToDatasetInsightSource(schema map[string]interface{}) api.InsightSources {
	doc := api.InsightSources{
		ID:      schema["id"].(string),
		Type:    schema["type"].(string),
		Insight: schema["insight"].(string),
	}

	return doc
}

func schemaToRelatedDatasets(schema []interface{}) []api.RelatedDataset {
	relatedDatasets := make([]api.RelatedDataset, len(schema))
	for i, s := range schema {

		relatedDatasets[i] = api.RelatedDataset{
			ID:   s.(map[string]interface{})["id"].(string),
			Type: s.(map[string]interface{})["type"].(string),
		}
	}
	return relatedDatasets
}

func resourceDatasetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	datasetId := d.Id()

	c := m.(*api.GFWClient)
	_, err := c.DeleteDataset(datasetId)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func flattenDatasetConfiguration(config api.DatasetConfiguration) interface{} {
	a := make(map[string]interface{})

	// API Supported Versions
	a["api_supported_versions"] = config.ApiSupportedVersions

	// Context Layer V1
	if config.ContextLayerV1 != nil {
		a["context_layer_v1"] = []interface{}{flattenContextLayerV1Config(*config.ContextLayerV1)}
	}

	// User Context Layer V1
	if config.UserContextLayerV1 != nil {
		a["user_context_layer_v1"] = []interface{}{flattenUserContextLayerV1Config(*config.UserContextLayerV1)}
	}

	// Temporal Context Layer V1
	if config.TemporalContextLayerV1 != nil {
		a["temporal_context_layer_v1"] = []interface{}{flattenTemporalContextLayerV1Config(*config.TemporalContextLayerV1)}
	}

	// User Tracks V1
	if config.UserTracksV1 != nil {
		a["user_tracks_v1"] = []interface{}{flattenUserTracksV1Config(*config.UserTracksV1)}
	}

	// PM Tiles V1
	if config.PmTilesV1 != nil {
		a["pm_tiles_v1"] = []interface{}{flattenPmTilesV1Config(*config.PmTilesV1)}
	}

	// Events V1
	if config.EventsV1 != nil {
		a["events_v1"] = []interface{}{flattenEventsV1Config(*config.EventsV1)}
	}

	// Fourwings V1
	if config.FourwingsV1 != nil {
		a["fourwings_v1"] = []interface{}{flattenFourwingsV1Config(*config.FourwingsV1)}
	}

	// Tracks V1
	if config.TracksV1 != nil {
		a["tracks_v1"] = []interface{}{flattenTracksV1Config(*config.TracksV1)}
	}

	// Frontend
	if config.Frontend != nil {
		a["frontend"] = []interface{}{flattenFrontendConfig(*config.Frontend)}
	}

	// Vessels V1
	if config.VesselsV1 != nil {
		a["vessels_v1"] = []interface{}{flattenVesselsV1Config(*config.VesselsV1)}
	}

	// Insights V1
	if config.InsightsV1 != nil {
		a["insights_v1"] = []interface{}{flattenInsightsV1Config(*config.InsightsV1)}
	}

	// Bulk Download V1
	if config.BulkDownloadV1 != nil {
		a["bulk_download_v1"] = []interface{}{flattenBulkDownloadV1Config(*config.BulkDownloadV1)}
	}

	// Data Download V1
	if config.DataDownloadV1 != nil {
		a["data_download_v1"] = []interface{}{flattenDataDownloadV1Config(*config.DataDownloadV1)}
	}

	// Thumbnails V1
	if config.ThumbnailsV1 != nil {
		a["thumbnails_v1"] = []interface{}{flattenThumbnailsV1Config(*config.ThumbnailsV1)}
	}

	return a
}

// Flatten functions for nested configurations

func flattenContextLayerV1Config(config api.ContextLayerV1Config) map[string]interface{} {
	a := make(map[string]interface{})
	a["import_logs"] = config.ImportLogs
	a["srid"] = config.Srid
	a["format"] = config.Format
	a["fields"] = config.Fields
	a["file_path"] = config.FilePath
	a["id_property"] = config.IDProperty
	return a
}

func flattenUserContextLayerV1Config(config api.UserContextLayerV1Config) map[string]interface{} {
	a := make(map[string]interface{})
	a["table"] = config.Table
	a["import_logs"] = config.ImportLogs
	a["srid"] = config.Srid
	a["format"] = config.Format
	a["fields"] = config.Fields
	a["file_path"] = config.FilePath
	a["id_property"] = config.IDProperty
	a["value_property_id"] = config.ValuePropertyID
	return a
}

func flattenTemporalContextLayerV1Config(config api.TemporalContextLayerV1Config) map[string]interface{} {
	a := make(map[string]interface{})
	a["dataset"] = config.Dataset
	a["project"] = config.Project
	a["source"] = config.Source
	a["table"] = config.Table
	return a
}

func flattenUserTracksV1Config(config api.UserTracksV1Config) map[string]interface{} {
	a := make(map[string]interface{})
	a["file_path"] = config.FilePath
	a["id_property"] = config.IDProperty
	return a
}

func flattenPmTilesV1Config(config api.PmTilesV1Config) map[string]interface{} {
	a := make(map[string]interface{})
	a["file_path"] = config.FilePath
	a["id_property"] = config.IDProperty
	return a
}

func flattenEventsV1Config(config api.EventsV1Config) map[string]interface{} {
	a := make(map[string]interface{})
	a["table"] = config.Table
	a["dataset"] = config.Dataset
	a["project"] = config.Project
	a["function"] = config.Function
	a["ttl"] = config.TTL
	a["max_zoom"] = config.MaxZoom
	a["source"] = config.Source
	return a
}

func flattenFourwingsV1Config(config api.FourwingsV1Config) map[string]interface{} {
	a := make(map[string]interface{})
	a["report_groupings"] = config.ReportGroupings
	a["table"] = config.Table
	a["dataset"] = config.Dataset
	a["max_zoom"] = config.MaxZoom
	a["project"] = config.Project
	a["function"] = config.Function
	a["intervals"] = config.Intervals
	a["ttl"] = config.TTL
	a["max"] = config.Max
	a["min"] = config.Min
	a["tile_scale"] = config.TileScale
	a["tile_offset"] = config.TileOffset
	a["internal_scale"] = config.InternalScale
	a["internal_offset"] = config.InternalOffset
	a["gee_band"] = config.GeeBand
	a["gee_images"] = config.GeeImages
	a["interaction_columns"] = config.InteractionColumns
	a["interaction_group_columns"] = config.InteractionGroupColumns
	a["temporal_aggregation"] = config.TemporalAggregation
	a["source"] = config.Source
	a["bucket"] = config.Bucket
	a["folder"] = config.Folder
	return a
}

func flattenTracksV1Config(config api.TracksV1Config) map[string]interface{} {
	a := make(map[string]interface{})
	a["bucket"] = config.Bucket
	a["folder"] = config.Folder
	a["database_instance"] = config.DatabaseInstance
	a["table"] = config.Table
	return a
}

func flattenFrontendConfig(config api.FrontendConfig) map[string]interface{} {
	a := make(map[string]interface{})
	a["max_zoom"] = config.MaxZoom
	a["translate"] = config.Translate
	a["max"] = config.Max
	a["min"] = config.Min
	a["disable_interaction"] = config.DisableInteraction
	a["latitude"] = config.Latitude
	a["longitude"] = config.Longitude
	a["start_time"] = config.StartTime
	a["end_time"] = config.EndTime
	a["timestamp"] = config.Timestamp
	a["geometry_type"] = config.GeometryType
	a["source_format"] = config.SourceFormat
	a["time_filter_type"] = config.TimeFilterType
	a["value_properties"] = config.ValueProperties
	a["polygon_color"] = config.PolygonColor
	a["point_size"] = config.PointSize
	a["line_id"] = config.LineID
	a["segment_id"] = config.SegmentID
	return a
}

func flattenVesselsV1Config(config api.VesselsV1Config) map[string]interface{} {
	a := make(map[string]interface{})
	a["index"] = config.Index
	a["index_boost"] = config.IndexBoost
	return a
}

func flattenInsightsV1Config(config api.InsightsV1Config) map[string]interface{} {
	a := make(map[string]interface{})
	if config.Sources != nil {
		a["sources"] = flattenDatasetInsightSources(config.Sources)
	}
	return a
}

func flattenBulkDownloadV1Config(config api.BulkDownloadV1Config) map[string]interface{} {
	a := make(map[string]interface{})
	a["gcs_uri"] = config.GcsUri
	a["path"] = config.Path
	a["format"] = config.Format
	a["compressed"] = config.Compressed
	return a
}

func flattenDataDownloadV1Config(config api.DataDownloadV1Config) map[string]interface{} {
	a := make(map[string]interface{})
	a["email_groups"] = config.EmailGroups
	a["gcs_folder"] = config.GcsFolder
	a["doi"] = config.Doi
	a["concept_doi"] = config.ConceptDOI
	return a
}

func flattenThumbnailsV1Config(config api.ThumbnailsV1Config) map[string]interface{} {
	a := make(map[string]interface{})
	a["extensions"] = config.Extensions
	a["bucket"] = config.Bucket
	a["folder"] = config.Folder
	a["scale"] = config.Scale
	return a
}

func flattenRelatedDatasets(relatedDatasets []api.RelatedDataset) []map[string]interface{} {
	list := make([]map[string]interface{}, len(relatedDatasets))
	for i, rd := range relatedDatasets {
		list[i] = map[string]interface{}{}
		list[i]["id"] = rd.ID
		list[i]["type"] = rd.Type

	}

	return list
}

func flattenDatasetInsightSources(docs []api.InsightSources) interface{} {

	array := make([]map[string]interface{}, len(docs))

	for i, doc := range docs {
		a := make(map[string]interface{})

		a["id"] = doc.ID
		a["type"] = doc.Type
		a["insight"] = doc.Insight

		array[i] = a
	}

	return array
}

func schemaToFilterConfig(data map[string]interface{}) api.FilterConfig {
	filter := api.FilterConfig{}

	if v, ok := data["label"].(string); ok {
		filter.Label = v
	}
	if v, ok := data["id"].(string); ok {
		filter.ID = v
	}
	if v, ok := data["type"].(string); ok {
		filter.Type = v
	}
	if v, ok := data["required"].(bool); ok {
		filter.Required = v
	}
	if v, ok := data["array"].(bool); ok {
		filter.Array = v
	}
	if v, ok := data["enum"].([]interface{}); ok && len(v) > 0 {
		filter.Enum = utils.ConvertArrayInterfaceToArrayString(v)
	}
	if v, ok := data["enabled"].(bool); ok {
		filter.Enabled = v
	}
	if v, ok := data["format"].(string); ok {
		filter.Format = v
	}
	if v, ok := data["max_length"].(int); ok {
		filter.MaxLength = v
	}
	if v, ok := data["min_length"].(int); ok {
		filter.MinLength = v
	}
	if v, ok := data["max"].(float64); ok {
		filter.Max = v
	}
	if v, ok := data["min"].(float64); ok {
		filter.Min = v
	}
	if v, ok := data["single_selection"].(bool); ok {
		filter.SingleSelection = v
	}
	if v, ok := data["operation"].(string); ok {
		filter.Operation = v
	}

	return filter
}

func schemaToDatasetFilters(data map[string]interface{}) api.DatasetFilters {
	filters := api.DatasetFilters{}

	if v, ok := data["fourwings"].([]interface{}); ok && len(v) > 0 {
		fourwingsFilters := make([]api.FilterConfig, len(v))
		for i, item := range v {
			fourwingsFilters[i] = schemaToFilterConfig(item.(map[string]interface{}))
		}
		filters.Fourwings = fourwingsFilters
	}

	if v, ok := data["events"].([]interface{}); ok && len(v) > 0 {
		eventsFilters := make([]api.FilterConfig, len(v))
		for i, item := range v {
			eventsFilters[i] = schemaToFilterConfig(item.(map[string]interface{}))
		}
		filters.Events = eventsFilters
	}

	if v, ok := data["vessels"].([]interface{}); ok && len(v) > 0 {
		vesselsFilters := make([]api.FilterConfig, len(v))
		for i, item := range v {
			vesselsFilters[i] = schemaToFilterConfig(item.(map[string]interface{}))
		}
		filters.Vessels = vesselsFilters
	}

	if v, ok := data["tracks"].([]interface{}); ok && len(v) > 0 {
		tracksFilters := make([]api.FilterConfig, len(v))
		for i, item := range v {
			tracksFilters[i] = schemaToFilterConfig(item.(map[string]interface{}))
		}
		filters.Tracks = tracksFilters
	}

	return filters
}

func flattenFilterConfig(filter api.FilterConfig) map[string]interface{} {
	result := make(map[string]interface{})

	if filter.Label != "" {
		result["label"] = filter.Label
	}
	if filter.ID != "" {
		result["id"] = filter.ID
	}
	if filter.Type != "" {
		result["type"] = filter.Type
	}
	result["required"] = filter.Required
	result["array"] = filter.Array
	if len(filter.Enum) > 0 {
		result["enum"] = filter.Enum
	}
	result["enabled"] = filter.Enabled
	if filter.Format != "" {
		result["format"] = filter.Format
	}
	if filter.MaxLength != 0 {
		result["max_length"] = filter.MaxLength
	}
	if filter.MinLength != 0 {
		result["min_length"] = filter.MinLength
	}
	if filter.Max != 0 {
		result["max"] = filter.Max
	}
	if filter.Min != 0 {
		result["min"] = filter.Min
	}
	result["single_selection"] = filter.SingleSelection
	if filter.Operation != "" {
		result["operation"] = filter.Operation
	}

	return result
}

func flattenDatasetFilters(filters api.DatasetFilters) map[string]interface{} {
	result := make(map[string]interface{})

	if len(filters.Fourwings) > 0 {
		fourwingsArray := make([]map[string]interface{}, len(filters.Fourwings))
		for i, filter := range filters.Fourwings {
			fourwingsArray[i] = flattenFilterConfig(filter)
		}
		result["fourwings"] = fourwingsArray
	}

	if len(filters.Events) > 0 {
		eventsArray := make([]map[string]interface{}, len(filters.Events))
		for i, filter := range filters.Events {
			eventsArray[i] = flattenFilterConfig(filter)
		}
		result["events"] = eventsArray
	}

	if len(filters.Vessels) > 0 {
		vesselsArray := make([]map[string]interface{}, len(filters.Vessels))
		for i, filter := range filters.Vessels {
			vesselsArray[i] = flattenFilterConfig(filter)
		}
		result["vessels"] = vesselsArray
	}

	if len(filters.Tracks) > 0 {
		tracksArray := make([]map[string]interface{}, len(filters.Tracks))
		for i, filter := range filters.Tracks {
			tracksArray[i] = flattenFilterConfig(filter)
		}
		result["tracks"] = tracksArray
	}

	return result
}

func schemaToDatasetDocumentation(data map[string]interface{}) api.DatasetDocumentation {
	doc := api.DatasetDocumentation{}

	if v, ok := data["type"].(string); ok {
		doc.Type = v
	}
	if v, ok := data["enable"].(bool); ok {
		doc.Enable = v
	}
	if v, ok := data["status"].(string); ok {
		doc.Status = v
	}
	if v, ok := data["queries"].([]interface{}); ok && len(v) > 0 {
		doc.Queries = utils.ConvertArrayInterfaceToArrayString(v)
	}
	if v, ok := data["provider"].(string); ok {
		doc.Provider = v
	}

	return doc
}

func flattenDatasetDocumentation(doc api.DatasetDocumentation) map[string]interface{} {
	result := make(map[string]interface{})

	if doc.Type != "" {
		result["type"] = doc.Type
	}
	result["enable"] = doc.Enable
	if doc.Status != "" {
		result["status"] = doc.Status
	}
	if len(doc.Queries) > 0 {
		result["queries"] = doc.Queries
	}
	if doc.Provider != "" {
		result["provider"] = doc.Provider
	}

	return result
}
