terraform {
  required_providers {
    qumulo = {
      source = "qumulo.com/terraform-intern/qumulo"
    }
  }
}

provider "qumulo" {
  username = "admin"
  password = "Admin123"
  host= "10.116.10.215"
  port= "17728"
}

variable "some_cluster_name" {
  type    = string
  default = "InigoMontoya"
}

resource "qumulo_cluster_name" "update_name" {
  name = var.some_cluster_name
}

//resource "qumulo_ssl_cert" "update_ssl" {
//  certificate = var.some_cert
//  private_key = var.some_key
//}

//resource "qumulo_ssl_ca" "update_ssl_ca" {
  //ca_certificate = var.some_cert
//}

resource "qumulo_monitoring" "update_monitoring" {
  enabled = false
  mq_host = "missionq.qumulo.com"
  mq_port = 443
  mq_proxy_host = ""
  mq_proxy_port = 17
  s3_proxy_host = "monitor.qumulo.com"
  s3_proxy_port = 443
  s3_proxy_disable_https = false
  vpn_enabled = false
  vpn_host = "ep1.qumulo.com"
  period = 60
}

resource "qumulo_smb_server" "update_smb" {
  session_encryption = "NONE"
  supported_dialects =["SMB2_DIALECT_2_002", "SMB2_DIALECT_2_1"]
  hide_shares_from_unauthorized_users = false
  hide_shares_from_unauthorized_hosts = true
  snapshot_directory_mode = "VISIBLE"
  bypass_traverse_checking = false
  signing_required = false
}

output "some_smb_server" {
  value = qumulo_smb_server.update_smb
}

output "some_monitoring_config" {
  value = qumulo_monitoring.update_monitoring
}

output "some_name" {
  value = qumulo_cluster_name.update_name
}

//output "some_authority" {
//  value = qumulo_ssl_ca.update_ssl_ca
//}

//output "some_ssl" {
//  value = qumulo_ssl_cert.update_ssl
//}