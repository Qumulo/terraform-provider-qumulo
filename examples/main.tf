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

# Configuring the SMB server settings
resource "qumulo_smb_server" "update_smb" {
  session_encryption = "NONE"
  supported_dialects =["SMB2_DIALECT_2_002"]
  hide_shares_from_unauthorized_users = false
  hide_shares_from_unauthorized_hosts = false
  snapshot_directory_mode = "VISIBLE"
  bypass_traverse_checking = false
  signing_required = false
}

# Setting the server's time configuration
resource "qumulo_time_configuration" "time_config" {
    use_ad_for_primary = false
    ntp_servers = ["0.qumulo.pool.ntp.org", "1.qumulo.pool.ntp.org"]
}