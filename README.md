# Terraform Provider for Qumulo

The Terraform Qumulo provider is a plugin for Terraform that allows for the updating of various settings of a Qumulo cluster. The provider makes use of the Qumulo REST API. The full list of supported features is below.

## Supported Features
- API Authentication
- Active Directory
- Audit Log
- Cluster Name
- Directory Quotas
- File System Settings
- FTP Server
- Interface and Network Configuration
- LDAP Server
- Local Users & Groups
- Monitoring (MQ)
- Network & Interface Configuration
- NFS Exports & Settings
- Roles
- SMB Server & Shares
- SSL & SSL CA
- Time Configuration
- Web UI Settings

## Starting Out
### Connecting to a Cluster
First, set the following environment variables, specifying the host and port of the cluster you wish to connect, and the appropriate credentials to log in to that cluster

    export QUMULO_HOST={host}
    export QUMULO_PORT={port}
    export QUMULO_USERNAME={username}
    export QUMULO_PASSWORD={password}

### Creating a Terraform Config File
Create a folder in which you want to initialize your Terraform workspace. Then, create a main.tf file within that folder, with the following header:

    terraform {
      required_providers {
        qumulo = {
          source = "Qumulo/qumulo"
          version = "0.1.1"
        }
      }
    }

Run `terraform init` to initialize the workspace.

If you'd like to import some or all of the current state of your cluster, you can either use [import_cluster.py](/examples/imports/import_cluster.py) (run `python3 import_cluster.py -h` to view usage), or import resources manually as shown [here](/IMPORT.md)

Then, add resources that you want to manage with Terraform. Examples of resources can be found [here](/examples/main.tf)

Now, run `terraform apply` to create those resources with the REST API. You're good to go!

## To Run Acceptance Tests
Make sure the environment variables as mentioned above are set. 

Make sure the TF_ACC environment variable is set to enable acceptance testing

    export TF_ACC=1

Then, run the following command

    make testwip

The above command runs all the acceptance tests in the provider.
To run a specific test, set the TESTNAME environment variable to the corresponding test name:

    export TESTNAME={testname}

Then, run the following command

    make runtest

## Developing the Qumulo Provider

For contributing to the Qumulo terraform provider, refer to the [development docs](https://github.com/Qumulo/terraform-provider-qumulo/blob/main/docs/TF-RESOURCE.md). It provides a brief overview for getting started on adding new Qumulo resources to be managed via Terraform and adding tests for the same. 

Check out our [style guide](/STYLE.md) before contributing to ensure the code base remains unified.

For further help on writing custom providers, refer to the [official Terraform documentation](https://www.hashicorp.com/blog/writing-custom-terraform-providers).
