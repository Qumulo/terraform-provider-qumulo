terraform {
  required_providers {
    qumulo = {
      version = "0.2"
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
  default = "NewName"
}

variable "some_cert" {
  type    = string
  default = "randomcertauth"
}
variable "some_key" {
  type    = string
  default = "randomkey"
}

resource "qumulo_cluster_name" "update_name" {
  name = var.some_cluster_name
}

resource "qumulo_ldap_server" "some_ldap_server" {
  use_ldap = true
  bind_uri = "ldap://ldap.denvrdata.com"
  user = ""
  base_distinguished_names = "dc=cloud,dc=denvrdata,dc=com"
  ldap_schema = "CUSTOM"
  ldap_schema_description {
    group_member_attribute = "memberUid"
    user_group_identifier_attribute = "uid"
    login_name_attribute =  "uid"
    group_name_attribute = "cn"
    user_object_class = "posixAccount"
    group_object_class = "posixGroup"
    uid_number_attribute = "uidNumber"
    gid_number_attribute = "gidNumber"
  }
  encrypt_connection = false
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
