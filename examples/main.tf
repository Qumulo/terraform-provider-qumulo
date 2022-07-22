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
  host= "10.116.100.110"
  port= "19299"
}

variable "some_cluster_name" {
  type    = string
  default = "qfsd"
}

resource "qumulo_cluster_name" "update_name" {
  name = var.some_cluster_name
}

resource "qumulo_ad_settings" "ad_settings" {
  signing = "WANT_SIGNING"
  sealing = "WANT_SEALING"
  crypto = "WANT_AES"
   domain = "ad.eng.qumulo.com"
   ad_username = "Administrator"
   ad_password = "a"
   use_ad_posix_attributes = false
   base_dn = "CN=Users,DC=ad,DC=eng,DC=qumulo,DC=com"
}

output "updated_name" {
  value = qumulo_cluster_name.update_name
}