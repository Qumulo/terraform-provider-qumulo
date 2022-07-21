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

output "updated_name" {
  value = qumulo_cluster_name.update_name
}