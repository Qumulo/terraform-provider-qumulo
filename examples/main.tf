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
  default = "NewName"
}

variable "some_cert" {
  type    = string
  default = "randomcert"
}
variable "some_key" {
  type    = string
  default = "randomkey"
}

data "qumulo_cluster_name" "all" {}

resource "qumulo_cluster_name" "update_name" {
  name = var.some_cluster_name
}

resource "qumulo_ssl_cert" "update_ssl" {
  certificate = var.some_cert
  private_key = var.some_key
}

output "some_ssl" {
  value = qumulo_ssl_cert.update_ssl
}

output "some_name" {
  value = data.qumulo_cluster_name.all
}

output "updated_name" {
  value = qumulo_cluster_name.update_name
}