---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "qumulo_ad_settings Resource - terraform-provider-qumulo"
subcategory: ""
description: |-
  
---

# qumulo_ad_settings (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `ad_password` (String, Sensitive)
- `ad_username` (String)
- `domain` (String)

### Optional

- `base_dn` (String)
- `crypto` (String)
- `domain_netbios` (String)
- `ou` (String)
- `sealing` (String)
- `signing` (String)
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `use_ad_posix_attributes` (Boolean)

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `update` (String)


