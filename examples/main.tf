terraform {
  required_providers {
    qumulo = {
      source = "qumulo.com/terraform-intern/qumulo"
    }
  }
}

# provider "qumulo" {
#   username = "admin"
#   password = "Admin123"
#   host= "10.116.10.215"
#   port= "17728"
# }

variable "some_cluster_name" {
  type    = string
  default = "InigoMontoya"
}
# variable "some_cert" {
#   type    = string
#   default = <<CERTDELIM
# -----BEGIN CERTIFICATE-----
# MIICIDCCAYmgAwIBAgIUZcdqCxZB1O4RD548ygFhGBXxQdQwDQYJKoZIhvcNAQEL
# BQAwIjEPMA0GA1UEAwwGVGVzdENBMQ8wDQYDVQQKDAZRdW11bG8wHhcNMjIwNzIy
# MTcwOTI4WhcNMzIwNzE5MTcwOTI4WjAiMQ8wDQYDVQQDDAZUZXN0Q0ExDzANBgNV
# BAoMBlF1bXVsbzCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEAv9Xupp43GfpI
# 0bVkB1BIa0ZBt5hpjxgee5PKwn3pbcg/M0M4qGhtX9/DR4utMqMib+X517hyo18E
# Vd+gZa0plafaPfwzz8YkO2EovYEFIaBxgqYkTQ0YZVt40cWEMMCWuyPndX0bvOrW
# 1f5zvOcc0+dDXoiqbhUDKiXBfzK745UCAwEAAaNTMFEwHQYDVR0OBBYEFKYiYrFK
# cZcR+gDTAqxV6u81B9htMB8GA1UdIwQYMBaAFKYiYrFKcZcR+gDTAqxV6u81B9ht
# MA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADgYEAjPXNGT38WwyWu4Xe
# Wngxmk0OIKZthsbZVDxSti3mse7KWadb6EkaRM/ZIO9CFPyB67zh3KAwhKiMbPVE
# JH62qN5t5xoqdDzzuOUHw1SSF78lfMAWk84TplzXegdysXjYFVhxvqYV9DIEhsTw
# HjX0jrbwN2tDfjTKNQwi7P7RPDY=
# -----END CERTIFICATE-----
# CERTDELIM
# }

resource "qumulo_cluster_name" "update_name" {
  cluster_name = var.some_cluster_name
}

# resource "qumulo_ad_settings" "ad_settings" {
#   signing = "WANT_SIGNING"
#   sealing = "WANT_SEALING"
#   crypto = "WANT_AES"
#   domain = "ad.eng.qumulo.com"
#   ad_username = "Administrator"
#   ad_password = "a"
#   use_ad_posix_attributes = false
#   base_dn = "CN=Users,DC=ad,DC=eng,DC=qumulo,DC=com"
# }

# resource "qumulo_ldap_server" "some_ldap_server" {
#   use_ldap = true
#   bind_uri = "ldap://ldap.denvrdata.com"
#   user = ""
#   base_distinguished_names = "dc=cloud,dc=denvrdata,dc=com"
#   ldap_schema = "CUSTOM"
#   ldap_schema_description {
#     group_member_attribute = "memberUid"
#     user_group_identifier_attribute = "uid"
#     login_name_attribute =  "uid"
#     group_name_attribute = "cn"
#     user_object_class = "posixAccount"
#     group_object_class = "posixGroup"
#     uid_number_attribute = "uidNumber"
#     gid_number_attribute = "gidNumber"
#   }
#   encrypt_connection = false
# }

# resource "qumulo_nfs_export" "new_nfs_export" {
#   export_path = "/lib"
#   fs_path = "/testing"
#   description = "testing nfs export via terraform"
#   restrictions {
#     host_restrictions = ["10.100.38.31"]
#     read_only = true
#     require_privileged_port = false
#     user_mapping = "NFS_MAP_ALL"
#     map_to_user = {
#       id_type = "LOCAL_USER"
#       id_value = "admin"
#     }
#   }
#   fields_to_present_as_32_bit = []
#   allow_fs_path_create = true
# }

#  resource "qumulo_nfs_export" "some_nfs_export" {
#    export_path = "/tmp"
#    fs_path = "/home/pthathamanjunatha"
#    description = "testing nfs export via terraform"
#    restrictions {
#      host_restrictions = ["10.100.38.31"]
#      read_only = false
#      require_privileged_port = false
#      user_mapping = "NFS_MAP_ALL"
#      map_to_user = {
#        id_type =  "NFS_UID"
#        id_value = "994"
#      }
#      map_to_group = {
#        id_type =  "NFS_GID"
#        id_value = "994"
#      }
#    }
#    fields_to_present_as_32_bit = ["FILE_IDS"]
#    allow_fs_path_create = true
#  }

//resource "qumulo_ssl_cert" "update_ssl" {
//  certificate = var.some_cert
//  private_key = var.some_key
//}

# resource "qumulo_ssl_ca" "update_ssl_ca" {
#   ca_certificate = var.some_cert
# }

# resource "qumulo_monitoring" "update_monitoring" {
#   enabled = false
#   mq_host = "missionq.qumulo.com"
#   mq_port = 443
#   mq_proxy_host = ""
#   mq_proxy_port = 32
#   s3_proxy_host = "monitor.qumulo.com"
#   s3_proxy_port = 443
#   s3_proxy_disable_https = false
#   vpn_enabled = false
#   vpn_host = "ep1.qumulo.com"
#   period = 60
# }

# resource "qumulo_smb_server" "update_smb" {
#   session_encryption = "NONE"
#   supported_dialects =["SMB2_DIALECT_2_002", "SMB2_DIALECT_2_1"]
#   hide_shares_from_unauthorized_users = false
#   hide_shares_from_unauthorized_hosts = true
#   snapshot_directory_mode = "VISIBLE"
#   bypass_traverse_checking = false
#   signing_required = false
# }

resource "qumulo_smb_share" "share1" {
  share_name = "TestingShare123"
  fs_path = "/"
  description = "This is a share used for testing purposes"
  permissions {
    type = "ALLOWED"
    trustee {
        domain = "LOCAL"
        auth_id = "501"
    }
    rights = ["READ", "WRITE", "CHANGE_PERMISSIONS"]
  }
  network_permissions {
    type = "ALLOWED"
    address_ranges = []
    rights = ["READ", "WRITE", "CHANGE_PERMISSIONS"]
  }
  access_based_enumeration_enabled = false
  bytes_per_sector = 512
  require_encryption = false
}

# output "some_smb_server" {
#   value = qumulo_smb_server.update_smb
# }

# output "some_monitoring_config" {
#   value = qumulo_monitoring.update_monitoring
# }

# output "some_name" {
#   value = qumulo_cluster_name.update_name
# }

# output "some_authority" {
#   value = qumulo_ssl_ca.update_ssl_ca
# }

//output "some_ssl" {
//  value = qumulo_ssl_cert.update_ssl
//}



  # permissions {
  #   type = "ALLOWED"
  #   trustee = {
  #       domain = "LOCAL"
  #       auth_id = "500"
  #       uid = null
  #       gid = null
  #       sid = "S-1-5-21-2821601516-3888733394-2762749966-500"
  #       name = null
  #   }
  #   rights = [
  #       "READ",
  #       "WRITE",
  #       "CHANGE_PERMISSIONS"
  #   ]
  # }
  # permissions = 
  #     {
  #         type = "ALLOWED"
  #         trustee = {
  #             domain = "WORLD"
  #             auth_id = "8589934592"
  #             sid = "S-1-1-0"
  #             name = null
  #         }
  #         rights = [
  #             "READ",
  #             "WRITE",
  #             "CHANGE_PERMISSIONS"
  #         ]
  #     }
