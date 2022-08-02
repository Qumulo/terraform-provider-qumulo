# **How to Import Resources**
There are a few steps necessary for importing resources into Terraform.

1. [Initialize the workspace](#initialize-workspace)
2. [Identify resources to be imported](#identify-resources)
3. [Import resources](#import-resources)
4. [Update configuration](#update-configuration)
5. [Apply changes](#apply-changes)

<a id="initialize-workspace"></a>
## Initialize the Workspace
First, create a .tf file with the qumulo provider configured. For example, create `main.tf` that includes

```
terraform {
  required_providers {
    qumulo = {
      source = "qumulo.com/terraform-intern/qumulo"
    }
  }
}
```

Then, run 

    $ terraform init

<a id="identify-resources"></a>
## Identify resources to be imported
Identify any resources you want to be imported into Terraform. For this tutorial, let's say we want to import the cluster name and the SMB share with ID 2. 

Add empty resource blocks corresponding to those resources to your .tf file.

```
resource "qumulo_cluster_name" "name" {}

resource "qumulo_smb_share" "share2" {}
```


<a id="import-resources"></a>
## Import resources
Now, import each resource using `terraform import`, making sure to include the ID. For resources where the ID does not matter, include the ID that you want terraform to associate with the resource. 
```
$ terraform import qumulo_cluster_name.name 1

$ terraform import qumulo_smb_share.share2 2
```
The resources are now imported into the Terraform state. However, you still need to update the configuration file.

<a id="update-configuration"></a>
## Update configuration
To view the current Terraform state, run
    
    $ terraform show

which should display something similar to

```
# qumulo_cluster_name.name:
resource "qumulo_cluster_name" "name" {
    cluster_name = "Buttercup"
    id           = "1"

    timeouts {}
}

# qumulo_smb_share.share2:
resource "qumulo_smb_share" "share2" {
    access_based_enumeration_enabled = false
    bytes_per_sector                 = "512"
    default_directory_create_mode    = "0755"
    default_file_create_mode         = "0644"
    description                      = "Decription"
    fs_path                          = "/"
    id                               = "2"
    require_encryption               = false
    share_name                       = "Files"

    network_permissions {
        address_ranges = []
        rights         = [
            "READ",
            "WRITE",
            "CHANGE_PERMISSIONS",
        ]
        type           = "ALLOWED"
    }

    permissions {
        rights = [
            "READ",
            "WRITE",
            "CHANGE_PERMISSIONS",
        ]
        type   = "DENIED"

        trustee {
            auth_id = "501"
            domain  = "LOCAL"
            gid     = 0
            name    = "guest"
            sid     = "S-1-5-21-1393870369-3041675342-41057371-501"
            uid     = 0
        }
    }
}
```

Copy that output to the .tf file, writing over the empty resource blocks we included earlier. Your .tf file should now look similar to this:

```
terraform {
  required_providers {
    qumulo = {
      source = "qumulo.com/terraform-intern/qumulo"
    }
  }
}

# qumulo_cluster_name.name:
resource "qumulo_cluster_name" "name" {
    cluster_name = "Buttercup"
}

# qumulo_smb_share.share2:
resource "qumulo_smb_share" "share2" {
    access_based_enumeration_enabled = false
    bytes_per_sector                 = "512"
    default_directory_create_mode    = "0755"
    default_file_create_mode         = "0644"
    description                      = "Description"
    fs_path                          = "/"
    require_encryption               = false
    share_name                       = "Files"

    network_permissions {
        address_ranges = []
        rights         = [
            "READ",
            "WRITE",
            "CHANGE_PERMISSIONS",
        ]
        type           = "ALLOWED"
    }

    permissions {
        rights = [
            "READ",
            "WRITE",
            "CHANGE_PERMISSIONS",
        ]
        type   = "DENIED"

        trustee {
            auth_id = "501"
            domain  = "LOCAL"
            gid     = 0
            name    = "guest"
            sid     = "S-1-5-21-1393870369-3041675342-41057371-501"
            uid     = 0
        }
    }
    permissions {
        rights = [
            "READ",
            "WRITE",
            "CHANGE_PERMISSIONS",
        ]
        type   = "ALLOWED"

        trustee {
            auth_id = "500"
            domain  = "LOCAL"
            gid     = 0
            name    = "admin"
            sid     = "S-1-5-21-1393870369-3041675342-41057371-500"
            uid     = 0
        }
    }
    permissions {
        rights = [
            "READ",
            "WRITE",
            "CHANGE_PERMISSIONS",
        ]
        type   = "ALLOWED"

        trustee {
            auth_id = "8589934592"
            domain  = "WORLD"
            gid     = 0
            name    = "Everyone"
            sid     = "S-1-1-0"
            uid     = 0
        }
    }
}
```

There are certain fields that are read only that you will have to manually remove, such as id. To view all such fields, run 

    $ terraform plan

This will show errors for all fields that need to be modified to have a valid configuration.

Once you have a valid configuration, `terraform plan` should not display any errors. 

This is also the time to adjust your configuration with any changes you wish to make.

<a id="apply-changes"></a>
## Apply changes
To apply the changes and sync with your cluster settings, run

    $ terraform apply

You should be good to go!

