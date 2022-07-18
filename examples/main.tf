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
  host= "https://10.116.100.110"
  port= "26561"
}