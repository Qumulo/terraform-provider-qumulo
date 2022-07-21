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
  default = "SuperNewName"
}

resource "qumulo_cluster_name" "update_name" {
  name = var.some_cluster_name
}

resource "qumulo_ad_settings" "ad_settings" {
  signing = "WANT_SIGNING"
  sealing = "WANT_SEALING"
  crypto = "WANT_AES"
  domain = "veryfakesite"
  ad_username = "fake"
  ad_password = "fake"
}

output "updated_name" {
  value = qumulo_cluster_name.update_name
}