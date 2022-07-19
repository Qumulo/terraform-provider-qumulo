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
  port= "26064"
}

variable "some_cluster_name" {
  type    = string
  default = "SuperNewName"
}

data "qumulo_cluster_name" "all" {}

resource "qumulo_cluster_name" "update_name" {
  name = var.some_cluster_name
}

resource "qumulo_ad_settings" "ad_settings" {
  signing = "WANT_SIGNING"
  sealing = "WANT_SEALING"
  crypto = "WANT_AES"
}

output "some_name" {
  value = data.qumulo_cluster_name.all
}

output "updated_name" {
  value = qumulo_cluster_name.update_name
}