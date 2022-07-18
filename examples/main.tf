terraform {
  required_providers {
    qumulo = {
      version = "0.2"
      source = "qumulo.com/terraform-intern/qumulo"
    }
  }
}

provider "qumulo" {
  username = "somerandomuser"
  password = "hdiosfhsdi"
  host= "10.116.100.110"
  port= "24100"
}

variable "cluster_name" {
  type    = string
  default = "Some random name"
}

data "only_one_name" "all" {}

# Returns all coffees
output "only_one_name" {
  value = data.hashicups_coffees.all.coffees
}