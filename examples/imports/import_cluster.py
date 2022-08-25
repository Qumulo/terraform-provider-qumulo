import argparse
import json
import logging
import os
import subprocess
import sys

import qumulo.lib.auth
import qumulo.lib.request

from qumulo.rest_client import RestClient

def main(opts):

    host = os.environ['QUMULO_HOST']
    port = os.environ['QUMULO_PORT']
    username = os.environ['QUMULO_USERNAME']
    password = os.environ['QUMULO_PASSWORD']

    if not host:
        print("Cluster host not configured, please set QUMULO_HOST environment variable")
        return
    if not port:
        print("Cluster port not configured, please set QUMULO_PORT environment variable")
        return
    if not username:
        print("Username not configured, please set QUMULO_USERNAME environment variable")
        return
    if not password:
        print("Password not configured, please set QUMULO_PASSWORD environment variable")
        return

    rc = RestClient(host, port)

    # Log in to Qumulo Core using local user or Active Directory credentials
    rc.login(username, password)
    
    if opts.dry:
        return

    while True:
        print(f"Warning: This will overwrite your config file ({opts.config_file}) and terraform state.")
        response = input("Do you wish to continue? (y/n) ")
        response = response.lower()
        if response == "y" or response == "yes":
            print("Starting import")
            break
        elif response == "n" or response == "no":
            print("Aborting import")
            return
        else:
            print("Invalid response, prompting again.")

    # TODO: implement selective importing
    enabled = {"cluster_name": True,
               "monitoring": True,
               "ssl_ca": True,
               "active_directory": True,
               "ldap": True,
               "time_config": True,
               "ftp": True,
               "fs_settings": True,
               "audit_log": True,
               "quotas": True,
               "local_users": True,
               "local_groups": True,
               "roles": True,
               "smb_settings": True,
               "smb_shares": True,
               "nfs_settings": True,
               "nfs_exports": True,
               "interfaces": True}

    if opts.enabled_features:
        for k in enabled:
            enabled[k] = False
        for feat in opts.enabled_features:
            if feat in enabled.keys():
                enabled[feat] = True
            else:
                print(f"Feature {feat} not supported. Supported features include {[f for f in enabled.keys()]}")
                return

    with open(opts.config_file, "w") as f:

        if opts.json_dump:
            # Dump settings to a json file instead
            settings = {}

            if enabled["cluster_name"]: settings["cluster_name"] = rc.cluster.get_cluster_conf()
            if enabled["monitoring"]: settings["monitoring"] = rc.support.get_config()
            if enabled["ssl_ca"]: 
                try:
                    settings["ssl_ca"] = rc.cluster.get_ssl_ca_certificate()
                except Exception as e:
                    logging.error('Unable to get SSL CA (%s)', e)
            if enabled["active_directory"]:
                settings["active_directory"] = {}
                settings["active_directory"]["settings"] = rc.ad.get_advanced_settings()
                settings["active_directory"]["ad"] = rc.ad.list_ad()
            if enabled["ldap"]: settings["ldap"] = rc.ldap.settings_get_v2()
            if enabled["time_config"]: settings["time_config"] = rc.time_config.get_time_status()["config"]
            if enabled["ftp"]: settings["ftp"] = rc.ftp.get_settings()
            if enabled["fs_settings"]: 
                settings["fs_settings"] = {}
                settings["fs_settings"]["perms"] = rc.fs.get_permissions_settings()
                settings["fs_settings"]["atime"] = rc.fs.get_atime_settings()
            if enabled["audit_log"]: settings["audit_log"] = rc.audit.get_syslog_config()
            if enabled["quotas"]: settings["quotas"] = next(rc.quota.get_all_quotas())["quotas"]
            if enabled["local_users"]: settings["local_users"] = rc.users.list_users()
            if enabled["local_groups"]: 
                settings["local_groups"] = rc.groups.list_groups()
                for g in settings["local_groups"]:
                    g["members"] = [str(m["id"]) for m in rc.groups.group_get_members(group_id=g["id"])]
            if enabled["roles"]: 
                settings["roles"] = rc.roles.list_roles()
                for name, r in settings["roles"].items():
                    r["members"] = [m for m in rc.roles.list_members(role_name=name)["members"]]
            if enabled["smb_settings"]: settings["smb_settings"] = rc.smb.get_smb_settings()
            if enabled["smb_shares"]: settings["smb_shares"] = rc.smb.smb_list_shares()
            if enabled["nfs_exports"]: settings["nfs_exports"] = rc.nfs.nfs_list_exports()
            if enabled["nfs_settings"]: settings["nfs_settings"] = rc.nfs.get_nfs_config()
            if enabled["interfaces"]: 
                settings["interfaces"] = rc.network.list_interfaces()
                settings["networks"] = {}
                for i in settings["interfaces"]:
                    settings["networks"][i["id"]] = rc.network.list_networks(interface_id=i["id"])

            
            f.write(json.dumps(settings, indent=4))

            return
                    

        # Setting up the provider
        f.write(getProviderBlock("qumulo.com/terraform-intern/qumulo"))
        f.write("\n")
        f.flush()

        os.system("rm terraform.tfstate")
        os.system("rm .terraform.lock.hcl")
        os.system("terraform init")

        if enabled["cluster_name"]: importClusterName(rc, f)
        if enabled["monitoring"]: importMonitoring(rc, f)
        if enabled["ssl_ca"]: 
            try:
                importSSLCA(rc, f)
            except Exception as e:
                logging.error('Unable to get SSL CA (%s)', e)
        if enabled["active_directory"]: importActiveDirectory(rc, f)
        if enabled["ldap"]: importLDAP(rc, f)
        if enabled["time_config"]: importTimeConfig(rc, f)
        if enabled["ftp"]: importFTPSettings(rc, f)
        if enabled["fs_settings"]: importFileSystemSettings(rc, f)
        if enabled["audit_log"]: importAuditLogSettings(rc, f)
        if enabled["quotas"]: importQuotas(rc, f)
        if enabled["local_users"]: importLocalUsers(rc, f)
        if enabled["local_groups"]: importLocalGroups(rc, f)
        if enabled["roles"]: importRoles(rc, f)
        if enabled["smb_settings"]: importSmbServerSettings(rc, f)
        if enabled["smb_shares"]: importSMBShares(rc, f)
        if enabled["nfs_settings"]: importNFSSettings(rc, f)
        if enabled["nfs_exports"]: importNFSExports(rc, f)
        if enabled["interfaces"]: importInterfaces(rc, f)
        

def importClusterName(rc, f):
    f.write(getClusterNameBlock(rc.cluster.get_cluster_conf()))
    f.write("\n")
    f.flush()

    os.system("terraform import qumulo_cluster_name.name 1")

def importMonitoring(rc, f):
    f.write(getMonitoringBlock(rc.support.get_config()))
    f.write("\n\n")
    f.flush()

    os.system("terraform import qumulo_monitoring.settings 1")

def importSSLCA(rc, f):
    f.write(getSSLCABlock(rc.cluster.get_ssl_ca_certificate()))
    f.write("\n")
    f.flush()

    os.system("terraform import qumulo_ssl_ca.certificate 1")

def importLDAP(rc, f):
    f.write(getLDAPBlock(rc.ldap.settings_get_v2()))
    f.write("\n")
    f.flush()

    os.system("terraform import qumulo_ldap_server.server 1")

def importTimeConfig(rc, f):
    f.write(getTimeConfigBlock(rc.time_config.get_time_status()["config"]))
    f.write("\n")
    f.flush()

    os.system("terraform import qumulo_time_configuration.time_config 1")

def importFTPSettings(rc, f):
    f.write(getFTPServerBlock(rc.ftp.get_settings()))
    f.write("\n")
    f.flush()

    os.system("terraform import qumulo_ftp_server.settings 1")

def importFileSystemSettings(rc, f):
    f.write(getFileSystemSettingsBlock(perms=rc.fs.get_permissions_settings(), 
                                       atime=rc.fs.get_atime_settings()))
    f.write("\n")
    f.flush()

    os.system("terraform import qumulo_file_system_settings.settings 1")

def importAuditLogSettings(rc, f):
    f.write(getAuditLogConfigBlock(rc.audit.get_syslog_config()))
    f.write("\n")
    f.flush()

    os.system("terraform import qumulo_syslog.config 1")

def importSmbServerSettings(rc, f):
    f.write(getSmbServerSettingsBlock(rc.smb.get_smb_settings()))
    f.write("\n")
    f.flush()

    os.system("terraform import qumulo_smb_server.settings 1")

def importRoles(rc, f):
    roles = rc.roles.list_roles()
    for name, role in roles.items():
        f.write(getRoleBlock(name=name, role=role))
        f.write("\n")
        f.flush()
        os.system(f"terraform import qumulo_role.{name} {name}")

        members = [m for m in rc.roles.list_members(role_name=name)["members"]]
        f.write(getRoleMembersBlock(role_name=name, members=members))
        f.write("\n")
        f.flush()
        for id in members:
            # we use subprocess here because os.system doesn't handle quotes well
            subprocess.call(["terraform", "import", f'qumulo_role_member.{name}["{str(id)}"]', 
                            f'{name}:{id}'])

def importQuotas(rc, f):
    quotas = next(rc.quota.get_all_quotas())["quotas"]
    for q in quotas:
        f.write(getQuotaBlock(q))
        f.write("\n")
        f.flush()
        os.system(f'terraform import qumulo_directory_quota.quota{q["id"]} {q["id"]}')

def importLocalGroups(rc, f):
    groups = rc.groups.list_groups()
    for g in groups:
        f.write(getLocalGroupBlock(g))
        f.write("\n")
        f.flush()
        os.system(f'terraform import qumulo_local_group.{g["name"]} {g["id"]}')

        member_ids = [str(m["id"]) for m in rc.groups.group_get_members(group_id=g["id"])]
        f.write(getLocalGroupMembersBlock(group=g, ids=member_ids))
        f.write("\n")
        f.flush()
        for id in member_ids:
            # we use subprocess here because os.system doesn't handle quotes well
            subprocess.call(["terraform", "import", f'qumulo_local_group_member.{g["name"]}["{str(id)}"]', 
                            f'{g["id"]}:{id}'])


def importLocalUsers(rc, f):
    users = rc.users.list_users()
    for u in users:
        f.write(getLocalUserBlock(u))
        f.write("\n")
        f.flush()
        os.system(f'terraform import qumulo_local_user.{u["name"]} {u["id"]}')

def importSMBShares(rc, f):
    shares = rc.smb.smb_list_shares()
    for s in shares:
        f.write(getSMBShareBlock(s))
        f.write("\n")
        f.flush()
        os.system(f'terraform import qumulo_smb_share.{s["share_name"]} {s["id"]}')

def importNFSSettings(rc, f):
    f.write(getNFSSettingsBlock(rc.nfs.get_nfs_config()))
    f.write("\n")
    f.flush()

    os.system("terraform import qumulo_nfs_settings.settings 1")

def importNFSExports(rc, f):
    exports = rc.nfs.nfs_list_exports()
    for x in exports:
        f.write(getNFSExportBlock(x))
        f.write("\n")
        f.flush()
        os.system(f'terraform import qumulo_nfs_export.export{x["id"]} {x["id"]}')

def importActiveDirectory(rc, f):
    f.write(getActiveDirectoryBlock(ad=rc.ad.list_ad(), settings=rc.ad.get_advanced_settings()))
    f.write("\n")
    f.flush()

    os.system("terraform import qumulo_ad_settings.ad_settings 1")

def importInterfaces(rc, f):
    interfaces = rc.network.list_interfaces()
    for i in interfaces:
        f.write(getInterfaceBlock(i))
        f.write("\n")
        f.flush()
        os.system(f'terraform import qumulo_interface_configuration.{i["name"]} {i["id"]}')

        networks = rc.network.list_networks(interface_id=i["id"])
        for n in networks:
            f.write(getNetworkConfigBlock(network=n, interface_id=i["id"]))
            f.write("\n")
            f.flush()
            os.system(f'terraform import qumulo_network_configuration.i{i["id"]}n{n["id"]} {i["id"]}:{n["id"]}')



def getProviderBlock(source: str) -> str:
    return  """terraform {{
  required_providers {{
    qumulo = {{
      source = "{0}"
    }}
  }}
}}
""".format(source)

def getClusterNameBlock(cluster_conf) -> str:
    return f"""resource "qumulo_cluster_name" "name" {{
  cluster_name = "{cluster_conf["cluster_name"]}"
}}
"""

def getMonitoringBlock(monitoring) -> str:
    return f"""resource "qumulo_monitoring" "settings" {{
  enabled = {str(monitoring["enabled"]).lower()}
  mq_host = "{monitoring["mq_host"]}"
  mq_port = {monitoring["mq_port"]}
  mq_proxy_host = "{monitoring["mq_proxy_host"]}"
  mq_proxy_port = {monitoring["mq_proxy_port"]}
  s3_proxy_host = "{monitoring["s3_proxy_host"]}"
  s3_proxy_port = {monitoring["s3_proxy_port"]}
  s3_proxy_disable_https = {str(monitoring["s3_proxy_disable_https"]).lower()}
  vpn_host = "{monitoring["vpn_host"]}"
  vpn_enabled = {str(monitoring["vpn_enabled"]).lower()}
  period = {monitoring["period"]}
}}
"""

def getSSLCABlock(ca) -> str:
    return f"""resource "qumulo_ssl_ca" "certificate" {{
  ca_certificate = <<CERTDELIM
{ca["ca_certificate"]}CERTDELIM
}}
"""

def getLDAPBlock(ldap) -> str:
    return f"""resource "qumulo_ldap_server" "server" {{
  use_ldap = {str(ldap["use_ldap"]).lower()}
  bind_uri = "{ldap["bind_uri"]}"
  user = "{ldap["user"]}"
  base_distinguished_names = "{ldap["base_distinguished_names"]}"
  ldap_schema = "{ldap["ldap_schema"]}"
  ldap_schema_description {{
    group_member_attribute = "{ldap["ldap_schema_description"]["group_member_attribute"]}"
    user_group_identifier_attribute = "{ldap["ldap_schema_description"]["user_group_identifier_attribute"]}"
    login_name_attribute =  "{ldap["ldap_schema_description"]["login_name_attribute"]}"
    group_name_attribute = "{ldap["ldap_schema_description"]["group_name_attribute"]}"
    user_object_class = "{ldap["ldap_schema_description"]["user_object_class"]}"
    group_object_class = "{ldap["ldap_schema_description"]["group_object_class"]}"
    uid_number_attribute = "{ldap["ldap_schema_description"]["uid_number_attribute"]}"
    gid_number_attribute = "{ldap["ldap_schema_description"]["gid_number_attribute"]}"
  }}
  encrypt_connection = {str(ldap["encrypt_connection"]).lower()}
}}
"""

def getTimeConfigBlock(conf) -> str:
    return f"""resource "qumulo_time_configuration" "time_config" {{
  use_ad_for_primary = {str(conf["use_ad_for_primary"]).lower()}
  ntp_servers = {str(conf["ntp_servers"]).replace("'", '"')}
}}
"""

def getFTPServerBlock(ftp) -> str:
    anonymous_user = ""
    if ftp["anonymous_user"]:
        anonymous_user = f"""  anonymous_user = {{
    id_type = "{ftp["anonymous_user"]["id_type"]}"
    id_value = "{ftp["anonymous_user"]["id_value"]}"
  }}\n"""

    return f"""resource "qumulo_ftp_server" "settings" {{
  enabled = {str(ftp["enabled"]).lower()}
  check_remote_host = {str(ftp["check_remote_host"]).lower()}
  log_operations = {str(ftp["log_operations"]).lower()}
  chroot_users = {str(ftp["chroot_users"]).lower()}
  allow_unencrypted_connections = {str(ftp["allow_unencrypted_connections"]).lower()}
  expand_wildcards = {str(ftp["expand_wildcards"]).lower()}
{anonymous_user}\
  greeting = "{ftp["greeting"]}"
}}
"""

def getFileSystemSettingsBlock(perms, atime) -> str:
    return f"""resource "qumulo_file_system_settings" "settings" {{
  permissions_mode = "{perms["mode"]}"
  atime_enabled = {str(atime["enabled"]).lower()}
  atime_granularity = "{atime["granularity"]}"
}}
"""

def getAuditLogConfigBlock(config) -> str:
    return f"""resource "qumulo_syslog" "config" {{
  enabled = {str(config["enabled"]).lower()}
  server_address = "{config["server_address"]}"
  server_port = {config["server_port"]}
}}
"""

def getRoleBlock(role, name) -> str:
    privileges = str(role["privileges"]).replace("'", '"').replace(" ", "\n\t\t\t\t")
    return f"""resource "qumulo_role" "{name}" {{
  description = "{role["description"]}"
  name        = "{name}"
  privileges  = {privileges}
}}
"""

def getRoleMembersBlock(role_name, members) -> str:
    return f"""resource qumulo_role_member "{role_name}" {{
  for_each = toset( {str(members).replace("'", '"')} )
  auth_id = each.key
  role_name = qumulo_role.{role_name}.id
}}
"""

def getSmbServerSettingsBlock(smb) -> str:
    return f"""resource "qumulo_smb_server" "settings" {{
  session_encryption = "{smb["session_encryption"]}"
  supported_dialects = {str(smb["supported_dialects"]).replace("'", '"')}
  hide_shares_from_unauthorized_users = {str(smb["hide_shares_from_unauthorized_users"]).lower()}
  hide_shares_from_unauthorized_hosts = {str(smb["hide_shares_from_unauthorized_hosts"]).lower()}
  snapshot_directory_mode = "{smb["snapshot_directory_mode"]}"
  bypass_traverse_checking = {str(smb["bypass_traverse_checking"]).lower()}
  signing_required = {str(smb["signing_required"]).lower()}
}}
"""

def getQuotaBlock(quota) -> str:
    return f"""resource "qumulo_directory_quota" "quota{quota["id"]}" {{
  directory_id = "{quota["id"]}"
  limit = "{quota["limit"]}"
}}
"""

def getLocalGroupBlock(group) -> str:
    gid = ''
    if group["gid"]:
        gid = f'\n  gid = "{group["gid"]}"'

    return f"""resource "qumulo_local_group" "{group["name"]}" {{
  name = "{group["name"]}"{gid}
}}
"""

def getLocalGroupMembersBlock(group, ids) -> str:
    return f"""resource "qumulo_local_group_member" "{group["name"]}" {{
  for_each = toset( {str(ids).replace("'", '"')} )
  member_id = each.key
  group_id = qumulo_local_group.{group["name"]}.id
}}
"""


def getLocalUserBlock(user) -> str:
    uid = ''
    if user["uid"]:
        uid = f'\n  uid = "{user["uid"]}"'

    home_directory = ''
    if user["home_directory"]:
        home_directory = f'\n  home_directory = "{user["home_directory"]}"'

    return f"""resource "qumulo_local_user" "{user["name"]}" {{
  name = "{user["name"]}"
  primary_group = {user["primary_group"]}{uid}{home_directory}
}}
"""

def getSMBShareBlock(share) -> str:
    permissions = ""
    network_permissions = ""
    default_file_create_mode = ""
    default_directory_create_mode = ""
    bytes_per_sector = ""

    if len(share["permissions"]) > 0:
        for p in share["permissions"]:
            permissions += getSMBPermissionsBlock(p)
        permissions = permissions.replace("\n", "\n  ") #indent properly

    if len(share["network_permissions"]) > 0:
        for p in share["network_permissions"]:
            network_permissions += getSMBNetworkPermissionsBlock(p)
        network_permissions = network_permissions.replace("\n", "\n  ") #indent properly
        network_permissions += "\n"
        

    if v := share["default_file_create_mode"]:
        default_file_create_mode = f'  default_file_create_mode = "{v}"\n'

    if v := share["default_directory_create_mode"]:
        default_directory_create_mode = f'  default_directory_create_mode = "{v}"\n'

    if v := share["bytes_per_sector"]:
        bytes_per_sector = f'  bytes_per_sector = "{v}"\n'

    return f"""resource "qumulo_smb_share" "{share["share_name"]}" {{
  share_name = "{share["share_name"]}"
  fs_path = "{share["fs_path"]}"
  description = "{share["description"]}"\
  {permissions}\
  {network_permissions}\
  access_based_enumeration_enabled = {str(share["access_based_enumeration_enabled"]).lower()}
  require_encryption = {str(share["require_encryption"]).lower()}
{default_file_create_mode}{default_directory_create_mode}{bytes_per_sector}\
}}
"""

def getSMBPermissionsBlock(perm) -> str:
    domain = ""
    name = ""
    auth_id = ""
    uid = ""
    gid = ""
    sid = ""
    if v := perm["trustee"]["domain"]:
        domain = f'    domain = "{v}"\n'
    if v := perm["trustee"]["name"]:
        name = f'    name = "{v}"\n'
    if v := perm["trustee"]["auth_id"]:
        auth_id = f'    auth_id = "{v}"\n'
    if v := perm["trustee"]["uid"]:
        uid = f'    uid = "{v}"\n'
    if v := perm["trustee"]["gid"]:
        gid = f'    gid = "{v}"\n'
    if v := perm["trustee"]["sid"]:
        sid = f'    sid = "{v}"\n'
    
    return f"""
permissions {{
  type = "{perm["type"]}"
  trustee {{
{domain}\
{name}\
{auth_id}\
{uid}\
{gid}\
{sid}\
  }}
  rights = {str(perm["rights"]).replace("'", '"')}
}}"""

def getSMBNetworkPermissionsBlock(perm) -> str:
    return f"""
network_permissions {{
  type = "{perm["type"]}"
  address_ranges = {str(perm["address_ranges"])}
  rights = {str(perm["rights"]).replace("'", '"')}
}}"""

def getNFSSettingsBlock(settings) -> str:
    return f"""resource "qumulo_nfs_settings" "settings" {{
  v4_enabled = {str(settings["v4_enabled"]).lower()}
  krb5_enabled = {str(settings["krb5_enabled"]).lower()}
  auth_sys_enabled = {str(settings["auth_sys_enabled"]).lower()}
}}
"""

def getNFSExportBlock(export) -> str:
    restrictions = ""
    for r in export["restrictions"]:
        restrictions += getNFSRestrictionsBlock(r)

    allow_fs_path_create = ""
    if "allow_fs_path_create" in export.keys():
        allow_fs_path_create = f'  allow_fs_path_create = {str(export["allow_fs_path_create"]).lower()}'

    return f"""resource "qumulo_nfs_export" "export{export["id"]}" {{
  export_path = "{export["export_path"]}"
  fs_path = "{export["fs_path"]}"
  description = "{export["description"]}"
{restrictions}
  fields_to_present_as_32_bit = {str(export["fields_to_present_as_32_bit"]).replace("'", '"')}\
{allow_fs_path_create}
}}
"""

def getNFSRestrictionsBlock(restriction) -> str:
    map_to_user = ""
    map_to_group = ""
    if "map_to_user" in restriction.keys():
        v = restriction["map_to_user"]
        map_to_user = f"""    map_to_user = {{
      id_type = "{v["id_type"]}"
      id_value = "{v["id_value"]}"
    }}\n"""

    if "map_to_group" in restriction.keys():
        v = restriction["map_to_group"]
        map_to_group = f"""    map_to_group = {{
      id_type = "{v["id_type"]}"
      id_value = "{v["id_value"]}"
    }}\n"""

    return f"""  restrictions {{
    host_restrictions = {str(restriction["host_restrictions"]).replace("'", '"')}
    read_only = {str(restriction["read_only"]).lower()}
    require_privileged_port = {str(restriction["require_privileged_port"]).lower()}
    user_mapping = "{restriction["user_mapping"]}"
{map_to_user}{map_to_group}\
  }}"""

def getActiveDirectoryBlock(ad, settings) -> str:
    domain_netbios = ""
    ou = ""
    base_dn = ""

    if ad["domain_netbios"]:
        domain_netbios = f'  domain_netbios = "{ad["domain_netbios"]}"\n'

    if ad["ou"]:
        ou = f'  ou = "{ad["ou"]}"\n'

    if ad["base_dn"]:
        ou = f'  base_dn = "{ad["base_dn"]}"\n'

    return f"""resource "qumulo_ad_settings" "ad_settings" {{
  signing = "{settings["signing"]}"
  sealing = "{settings["sealing"]}"
  crypto = "{settings["crypto"]}"
  domain = "{ad["domain"]}"
{ou}\
{domain_netbios}\
{base_dn}\
  ad_username = "insert_ad_username_here"
  ad_password = "insert_ad_password_here"
  use_ad_posix_attributes = {str(ad["use_ad_posix_attributes"]).lower()}
}}
"""

def getInterfaceBlock(interface) -> str:
    default_ipv6 = ""
    if interface["default_gateway_ipv6"]:
        default_ipv6 = f'  default_gateway_ipv6 = "{interface["default_gateway_ipv6"]}"\n'

    return f"""resource "qumulo_interface_configuration" "{interface["name"]}" {{
  name = "{interface["name"]}"
  default_gateway = "{interface["default_gateway"]}"
{default_ipv6}\
  bonding_mode = "{interface["bonding_mode"]}"
  mtu = {interface["mtu"]}
  interface_id = "{interface["id"]}"
}}
"""

def getNetworkConfigBlock(network, interface_id) -> str:
    return f"""resource "qumulo_network_configuration" "i{interface_id}n{network["id"]}" {{
  interface_id = "{interface_id}"
  assigned_by        = "{network["assigned_by"]}"
  dns_search_domains = {str(network["dns_search_domains"]).replace("'", '"')}
  dns_servers        = {str(network["dns_servers"]).replace("'", '"')}
  floating_ip_ranges = {str(network["floating_ip_ranges"]).replace("'", '"')}
  ip_ranges          = {str(network["ip_ranges"]).replace("'", '"')}
  mtu                = {network["mtu"]}
  name               = "{network["name"]}"
  vlan_id            = {network["vlan_id"]}
  netmask = "{network["netmask"]}"
  network_id = "{network["id"]}"
}}
"""



if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Import Qumulo cluster settings into a Terraform configuration")
    parser.add_argument('--config_file', '--file', dest='config_file', type=str, help='Configuration file to write to',
                        default='main.tf')
    parser.add_argument('--dry','-d', dest='dry', action='store_true', help='Dry run, test rc commands',
                        default=False)
    parser.add_argument('--enable', '-e', action='extend', dest='enabled_features', nargs="+", type=str, 
                        help="Enable specific features to be imported. If this flag is not present, all features are imported.")
    parser.add_argument('--json','-j', dest='json_dump', action='store_true',
                        help='Dump to a JSON file instead of a Terraform config', default=False)
    args = parser.parse_args()
    main(args)
    # main(config_file=args.config_file, dry=args.dry, enabled_features=args.enable, json_dump=args.json)
