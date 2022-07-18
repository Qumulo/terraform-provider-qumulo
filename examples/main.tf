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
  host= "https://10.116.100.110:24100"
  port= "24100"
}

//variable "some_cluster_name" {
//  type    = string
//  default = "Some random name"
//}

data "qumulo_cluster_name" "all" {}

# Returns all coffees
//output "some_name" {
//  value = data.qumulo.all.cluster_name
//}