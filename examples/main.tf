# Configuring the required provider

terraform {
  required_providers {
    qumulo = {
      source = "Qumulo/qumulo"
      version = "0.1.1"
    }
  }
}

# Optional: Configuring the provider credentials. This will override any environment variables set

# provider "qumulo" {
#   username = "<username>"
#   password = "<password>"
#   host= "<hostname>"
#   port= "<port>"
# }

# Optional: Setting up some variables. These can instead be put directly into the resource body

variable "some_cluster_name" {
  type    = string
  default = "InigoMontoya"
}
 variable "some_cert" {
   type    = string
   default = <<CERTDELIM
 -----BEGIN CERTIFICATE-----
 MIICIDCCAYmgAwIBAgIUZcdqCxZB1O4RD548ygFhGBXxQdQwDQYJKoZIhvcNAQEL
 BQAwIjEPMA0GA1UEAwwGVGVzdENBMQ8wDQYDVQQKDAZRdW11bG8wHhcNMjIwNzIy
 MTcwOTI4WhcNMzIwNzE5MTcwOTI4WjAiMQ8wDQYDVQQDDAZUZXN0Q0ExDzANBgNV
 BAoMBlF1bXVsbzCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEAv9Xupp43GfpI
 0bVkB1BIa0ZBt5hpjxgee5PKwn3pbcg/M0M4qGhtX9/DR4utMqMib+X517hyo18E
 Vd+gZa0plafaPfwzz8YkO2EovYEFIaBxgqYkTQ0YZVt40cWEMMCWuyPndX0bvOrW
 1f5zvOcc0+dDXoiqbhUDKiXBfzK745UCAwEAAaNTMFEwHQYDVR0OBBYEFKYiYrFK
 cZcR+gDTAqxV6u81B9htMB8GA1UdIwQYMBaAFKYiYrFKcZcR+gDTAqxV6u81B9ht
 MA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADgYEAjPXNGT38WwyWu4Xe
 Wngxmk0OIKZthsbZVDxSti3mse7KWadb6EkaRM/ZIO9CFPyB67zh3KAwhKiMbPVE
 JH62qN5t5xoqdDzzuOUHw1SSF78lfMAWk84TplzXegdysXjYFVhxvqYV9DIEhsTw
 HjX0jrbwN2tDfjTKNQwi7P7RPDY=
 -----END CERTIFICATE-----
 CERTDELIM
 }

# Setting the cluster name and SSL Certificate Authority
resource "qumulo_cluster_name" "update_name" {
  cluster_name = var.some_cluster_name
}

resource "qumulo_ssl_ca" "update_ssl_ca" {
  ca_certificate = var.some_cert
}

# Configuring the monitoring settings
resource "qumulo_monitoring" "update_monitoring" {
  mq_host = "missionq.qumulo.com"
  mq_port = 443
  mq_proxy_host = ""
  mq_proxy_port = 32
  s3_proxy_host = "monitor.qumulo.com"
  s3_proxy_port = 443
  vpn_host = "ep1.qumulo.com"
  period = 60
}

# Configuring the Active Directory settings and joining to an AD domain
resource "qumulo_ad_settings" "ad_settings" {
  signing = "WANT_SIGNING"
  sealing = "WANT_SEALING"
  crypto = "WANT_AES"
  domain = "ad.eng.qumulo.com"
  domain_netbios = "AD"
  ad_username = "Administrator"
  ad_password = "a"
  use_ad_posix_attributes = false
  base_dn = "CN=Users,DC=ad,DC=eng,DC=qumulo,DC=com"
}

# Configuring NFS settings
resource "qumulo_nfs_settings" "my_new_settings" {
  v4_enabled = false
  krb5_enabled = true
  auth_sys_enabled = true
}

# Setting an NFS export
resource "qumulo_nfs_export" "new_nfs_export" {
  export_path = "/lib"
  fs_path = "/testing"
  description = "testing nfs export via terraform"
  restrictions {
    host_restrictions = ["10.100.38.31"]
    read_only = true
    require_privileged_port = false
    user_mapping = "NFS_MAP_ALL"
    map_to_user = {
      id_type = "LOCAL_USER"
      id_value = "admin"
    }
  }
  fields_to_present_as_32_bit = []
  allow_fs_path_create = true
}

resource "qumulo_nfs_export" "some_nfs_export" {
  export_path = "/tmp"
  fs_path = "/home/pthathamanjunatha"
  description = "testing nfs export via terraform"
  restrictions {
    host_restrictions = ["10.100.38.31"]
    read_only = false
    require_privileged_port = false
    user_mapping = "NFS_MAP_ALL"
    map_to_user = {
      id_type =  "NFS_UID"
      id_value = "994"
    }
    map_to_group = {
      id_type =  "NFS_GID"
      id_value = "994"
    }
  }
  fields_to_present_as_32_bit = ["FILE_IDS"]
  allow_fs_path_create = true
}

# Configuring the LDAP server settings
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

# Configuring the SMB server settings
resource "qumulo_smb_server" "update_smb" {
  session_encryption = "NONE"
  supported_dialects =["SMB2_DIALECT_2_002", "SMB2_DIALECT_2_1"]
  hide_shares_from_unauthorized_users = false
  hide_shares_from_unauthorized_hosts = true
  snapshot_directory_mode = "VISIBLE"
  bypass_traverse_checking = false
  signing_required = false
}

# Creating an SMB share
resource "qumulo_smb_share" "share1" {
  share_name = "TestingShareHi344"
  fs_path = "/"
  description = "This is a share used for testing purposes"
  permissions {
    type = "ALLOWED"
    trustee {
      domain = "LOCAL"
      name = "admin"
    }
    rights = ["READ", "WRITE", "CHANGE_PERMISSIONS"]
  }
  permissions {
    type = "DENIED"
    trustee {
      domain = "LOCAL"
      uid = 65534
    }
    rights = ["WRITE"]
  }
  network_permissions {
    type = "ALLOWED"
    address_ranges = []
    rights = ["READ", "WRITE", "CHANGE_PERMISSIONS"]
  }
  access_based_enumeration_enabled = false
  require_encryption = false
}

# Setting the server's time configuration
resource "qumulo_time_configuration" "time_config" {
    use_ad_for_primary = false
    ntp_servers = ["0.qumulo.pool.ntp.org", "1.qumulo.pool.ntp.org"]
}

# Setting an interface configuration
resource "qumulo_interface_configuration" "interface_config" {
  name = "bond0"
  default_gateway = "10.220.0.1"
  bonding_mode = "IEEE_8023AD"
  mtu = 1500
  interface_id = "1"
}

# Setting a network configuration
resource "qumulo_network_configuration" "network_config" {
  interface_id = "1"
  assigned_by        = "DHCP"
  dns_search_domains = []
  dns_servers        = []
  floating_ip_ranges = []
  ip_ranges          = []
  mtu                = 1500
  name               = "Default"
  vlan_id            = 0
  netmask = ""
  network_id = "1"
}

# Setting a directory quota for the directory with ID 2
resource "qumulo_directory_quota" "new_quota" {
    directory_id = "2"
    limit = "1000000000"
}

# Creating 3 local users with the same group, password, and home directory
resource "qumulo_local_user" "test_user" {
  for_each = toset( ["testuser1", "testuser2", "testuser3"] )
  name = each.key
  primary_group = 514
  password = "Test1234"
  home_directory = "/"
}

# Creating a single test user
resource "qumulo_local_user" "test_user2" {
  name = "mytestuser32"
  primary_group = 514
  password = "Test1234"
  home_directory = "/"
}

# Creating a local group
resource "qumulo_local_group" "test_group" {
    name = "testgroup2"
}

# Adding local users from the test_user resource to a local group
resource "qumulo_local_group_member" "test_member" {
  for_each = qumulo_local_user.test_user
  member_id = each.value.id
  group_id = qumulo_local_group.test_group.id
}

# Creating a role
resource "qumulo_role" "actors" {
  description = "Testing testing 123"
  name        = "Actors"
  privileges  = [
    "PRIVILEGE_AD_READ",
    "PRIVILEGE_AD_USE",
    "PRIVILEGE_AD_WRITE",
  ]
}

# Adding users to a role
resource "qumulo_role_member" "actors_member" {
  for_each = qumulo_local_user.test_user
  name = each.value.name
  role_name = qumulo_role.actors.name
}

# Configuring the file system settings
resource "qumulo_file_system_settings" "fs_settings" {
  permissions_mode = "CROSS_PROTOCOL"
  atime_enabled = false
  atime_granularity = "HOUR"
}

# Configuring the web UI settings
resource "qumulo_web_ui" "settings" {
  inactivity_timeout {
    nanoseconds = 900000000001
  }
  login_banner = "SampleBanner"
}

# Configuring the audit log settings
resource "qumulo_syslog" "syslog_audit_log" {
  enabled = false
  server_address = ""
  server_port = 0
}
resource "qumulo_cloudwatch" "cloudwatch_audit_log" {
  enabled = false
  log_group_name = ""
  region = ""
}

# Configuring the FTP server settings
resource "qumulo_ftp_server" "some_ftp_server" {
  enabled = true
  check_remote_host = false
  log_operations = true
  chroot_users = true
  allow_unencrypted_connections = true
  expand_wildcards = false
  greeting = "Hello!"
}

