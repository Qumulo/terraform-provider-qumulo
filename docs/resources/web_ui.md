---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "qumulo_web_ui Resource - terraform-provider-qumulo"
subcategory: ""
description: |-
  
---

# qumulo_web_ui (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `inactivity_timeout` (Block List, Min: 1, Max: 1) (see [below for nested schema](#nestedblock--inactivity_timeout))

### Optional

- `login_banner` (String)
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--inactivity_timeout"></a>
### Nested Schema for `inactivity_timeout`

Required:

- `nanoseconds` (String)


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `update` (String)

