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
  default = "WhatsInAName"
}

resource "qumulo_cluster_name" "update_name" {
  name = var.some_cluster_name
}

output "updated_name" {
  value = qumulo_cluster_name.update_name
}