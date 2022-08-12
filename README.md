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

## To Run Acceptance Tests
First, set the following environment variables, specifying the host and port of the cluster you wish to test on, and the appropriate credentials to log in to that cluster

    export QUMULO_HOST={host}
    export QUMULO_PORT={port}
    export QUMULO_USERNAME={username}
    export QUMULO_PASSWORD={password}

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

More information coming soon!

Check out our [style guide](/STYLE.md) before contributing to ensure the code base remains unified.
