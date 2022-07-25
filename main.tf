terraform {
  required_providers {
    qumulo = {
      source = "qumulo.com/terraform-intern/qumulo"
      version = "0.2"
    }
  }
}

provider "qumulo" {
  username = "admin"
  password = "Admin123"
  host= "10.116.10.215"
  port= "17728"
}

resource "qumulo_cluster_name" "update_name" {
  name = "Inconceivable"
}