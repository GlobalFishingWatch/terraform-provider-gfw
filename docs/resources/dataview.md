---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "gfw_dataview Resource - terraform-provider-gfw"
subcategory: ""
description: |-
  
---

# gfw_dataview (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `app` (String)
- `description` (String)
- `name` (String)
- `slug` (String)

### Optional

- `category` (String)
- `config` (Block List, Max: 1) (see [below for nested schema](#nestedblock--config))
- `created_at` (String)
- `datasets_config` (List of String)
- `info_config` (String)
- `updated_at` (String)

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--config"></a>
### Nested Schema for `config`

Optional:

- `color` (String)
- `color_ramp` (String)
- `datasets` (List of String)
- `type` (String)


