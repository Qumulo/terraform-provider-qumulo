# Configuring the required provider

terraform {
  required_providers {
    qumulo = {
      source = "qumulo.com/terraform-intern/qumulo"
    }
  }
}

# Optional: Configuring the provider credentials. This will override any environment variables set

provider "qumulo" {
  username = "admin"
  password = "Admin123"
  host= "10.116.10.215"
  port= "22848"
}

# Setting the cluster name and SSL Certificate Authority
resource "qumulo_cluster_name" "update_name" {
  cluster_name = "InigiMontoya"
}

# Configuring the monitoring settings
resource "qumulo_monitoring" "update_monitoring" {
  mq_host = "missionq.qumulo.com"
  mq_port = 443
  mq_proxy_host = ""
  mq_proxy_port = 32
  s3_proxy_host = "monitor.qumulo.com"
  s3_proxy_port = 445
  vpn_host = "ep1.qumulo.com"
  period = 60
}

# Configuring NFS settings
resource "qumulo_nfs_settings" "my_new_settings" {
  v4_enabled = false
  krb5_enabled = false
  krb5p_enabled = false
  auth_sys_enabled = true
}