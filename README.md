# Terraform Provider for Qumulo

The Terraform Qumulo provider is a plugin for Terraform that allows for the updating of various settings of a Qumulo cluster. The provider makes use of the Qumulo REST API. The full list of supported features is below.

## Supported Features
- API Authentication
- Cluster Name
- Monitoring
- SSL/SSL CA
- LDAP

## To Run Acceptance Tests
First, set the following environment variables, specifying the host and port of the cluster you wish to test on, and the appropriate credentials to log in to that cluster

    export QUMULO_HOST={host}
    export QUMULO_PORT={port}
    export QUMULO_USERNAME={username}
    export QUMULO_PASSWORD={password}

Then, run the following command

    make testwip
