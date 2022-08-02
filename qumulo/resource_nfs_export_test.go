package qumulo

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestAccTestNfsExport(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: defaultAccNfsExport(defaultNfsExport),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNfsExport(defaultNfsExport, "qumulo_nfs_export.default_nfs_export")),
			},
			{
				Config: testAccNfsExport(testNfsExport),
				Check: resource.ComposeTestCheckFunc(
					testAccCompareNfsExports(testNfsExport),
					testAccCheckNfsExport(testNfsExport, "qumulo_nfs_export.some_nfs_export")),
			},
		},
	})
}

var defaultNfsExport = NfsExport{
	ExportPath:  "/my_very_own_export",
	FsPath:      "/home/testing/my_very_own_export",
	Description: "",
	Restrictions: []NfsRestriction{{
		HostRestrictions:      []string{},
		ReadOnly:              false,
		RequirePrivilegedPort: false,
		UserMapping:           "NFS_MAP_NONE",
	}},
	FieldsToPresentAs32Bit: []string{},
}

var testNfsExport = NfsExport{
	ExportPath:  "/my_very_own_export",
	FsPath:      "/home/testing/my_very_own_export",
	Description: "",
	Restrictions: []NfsRestriction{{
		HostRestrictions:      []string{"10.100.38.31"},
		ReadOnly:              false,
		RequirePrivilegedPort: false,
		UserMapping:           "NFS_MAP_ALL",
		MapToUser: map[string]interface{}{
			"id_type":  "NFS_UID",
			"id_value": "994",
		},
		MapToGroup: map[string]interface{}{
			"id_type":  "NFS_GID",
			"id_value": "994",
		},
	}},
	FieldsToPresentAs32Bit: []string{"FILE_IDS"},
}

func defaultAccNfsExport(ne NfsExport) string {
	return fmt.Sprintf(`
 resource "qumulo_nfs_export" "default_nfs_export" {
   export_path = %q
   fs_path = %q
   description = %q
   restrictions {
     host_restrictions = %v
     read_only = %v
     require_privileged_port = %v
     user_mapping = %q
   }
   fields_to_present_as_32_bit = %v
   allow_fs_path_create = true
 }`, ne.ExportPath, ne.FsPath, ne.Description, strings.ReplaceAll(fmt.Sprintf("%+q", ne.Restrictions[0].HostRestrictions), "\" \"", "\", \""), ne.Restrictions[0].ReadOnly,
		ne.Restrictions[0].RequirePrivilegedPort, ne.Restrictions[0].UserMapping,
		strings.ReplaceAll(fmt.Sprintf("%+q", ne.FieldsToPresentAs32Bit), "\" \"", "\", \""))
}

func testAccNfsExport(ne NfsExport) string {
	return fmt.Sprintf(`
 resource "qumulo_nfs_export" "some_nfs_export" {
   export_path = %q
   fs_path = %q
   description = %q
   restrictions {
     host_restrictions = %v
     read_only = %v
     require_privileged_port = %v
     user_mapping = %q
     map_to_user = {
       id_type =  %q
       id_value = %q
     }
     map_to_group = {
       id_type =  %q
       id_value = %q
     }
   }
   fields_to_present_as_32_bit = %v
   allow_fs_path_create = true
 }`, ne.ExportPath, ne.FsPath, ne.Description, strings.ReplaceAll(fmt.Sprintf("%+q", ne.Restrictions[0].HostRestrictions), "\" \"", "\", \""), ne.Restrictions[0].ReadOnly,
		ne.Restrictions[0].RequirePrivilegedPort, ne.Restrictions[0].UserMapping, ne.Restrictions[0].MapToUser["id_type"],
		ne.Restrictions[0].MapToUser["id_value"], ne.Restrictions[0].MapToGroup["id_type"], ne.Restrictions[0].MapToGroup["id_value"],
		strings.ReplaceAll(fmt.Sprintf("%+q", ne.FieldsToPresentAs32Bit), "\" \"", "\", \""))
}

func testAccCompareNfsExports(ne NfsExport) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("qumulo_nfs_export.some_nfs_export", "export_path", fmt.Sprintf("%v", ne.ExportPath)),
		resource.TestCheckResourceAttr("qumulo_nfs_export.some_nfs_export", "fs_path", fmt.Sprintf("%v", ne.FsPath)),
		resource.TestCheckResourceAttr("qumulo_nfs_export.some_nfs_export", "description", fmt.Sprintf("%v", ne.Description)),
		resource.TestCheckResourceAttr("qumulo_nfs_export.some_nfs_export", "restrictions.#", strconv.Itoa(len(ne.Restrictions))),
		resource.TestCheckResourceAttr("qumulo_nfs_export.some_nfs_export", "restrictions.0.host_restrictions.#", strconv.Itoa(len(ne.Restrictions[0].HostRestrictions))),
		resource.TestCheckResourceAttr("qumulo_nfs_export.some_nfs_export", "restrictions.0.read_only", fmt.Sprintf("%v", ne.Restrictions[0].ReadOnly)),
		resource.TestCheckResourceAttr("qumulo_nfs_export.some_nfs_export", "restrictions.0.require_privileged_port", fmt.Sprintf("%v", ne.Restrictions[0].RequirePrivilegedPort)),
		resource.TestCheckResourceAttr("qumulo_nfs_export.some_nfs_export", "restrictions.0.user_mapping", fmt.Sprintf("%v", ne.Restrictions[0].UserMapping)),
		resource.TestCheckResourceAttr("qumulo_nfs_export.some_nfs_export", "restrictions.0.map_to_user.id_type", fmt.Sprintf("%v", ne.Restrictions[0].MapToUser["id_type"])),
		resource.TestCheckResourceAttr("qumulo_nfs_export.some_nfs_export", "restrictions.0.map_to_user.id_value", fmt.Sprintf("%v", ne.Restrictions[0].MapToUser["id_value"])),
		resource.TestCheckResourceAttr("qumulo_nfs_export.some_nfs_export", "restrictions.0.map_to_group.id_type", fmt.Sprintf("%v", ne.Restrictions[0].MapToGroup["id_type"])),
		resource.TestCheckResourceAttr("qumulo_nfs_export.some_nfs_export", "restrictions.0.map_to_group.id_value", fmt.Sprintf("%v", ne.Restrictions[0].MapToGroup["id_value"])),
		resource.TestCheckResourceAttr("qumulo_nfs_export.some_nfs_export", "fields_to_present_as_32_bit.#", strconv.Itoa(len(ne.FieldsToPresentAs32Bit))))
}

func testAccCheckNfsExport(ne NfsExport, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*Client)
		ctx := context.Background()

		res, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("terraform resource not found: %s", resourceName)
		}

		nfsExportId := res.Primary.ID
		getNfsExportByIdUri := NfsExportsEndpoint + nfsExportId

		export, err := DoRequest[NfsExport, NfsExport](ctx, c, GET, getNfsExportByIdUri, nil)
		if err != nil {
			return err
		}
		ne.Id = nfsExportId
		if !reflect.DeepEqual(*export, ne) {
			return fmt.Errorf("NFS export mismatch: Expected %v, got %v", ne, *export)
		}
		return nil
	}
}
