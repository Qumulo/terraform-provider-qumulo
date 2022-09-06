# Defining Qumulo resources

The terraform provider for Qumulo manages creating, updating, reading and destroying infrastructure on a Qumulo's cluster.
Each of the infrastructure configurations are modeled as terraform resources. 
The below guidelines can be useful for adding support for new features (a.k.a terraform resources). A prerequisite is to have a corresponding set of REST endpoints that can be used for some/all of the CRUD operations on the resource.

1. As a convention, each of the terraform resources are defined in their own file, named after the resource with a prefix ``resource_``. For example,
to create an ``example_settings`` resource, create a file named ``resource_example_settings.go``.

2. Given that each of the Terraform resources talk to a specific REST endpoint, we define the REST endpoint prefix as a constant at the top of the file.

3. Next, we define any enum values (using Go string slices and constants) that might be present as part of the Terraform schema. 

4. The next step is to define the Go structs for the HTTP request and response bodies. In the case the request body of the CREATE/UPDATE endpoints is similar to the response body of the READ endpoint, we define a common struct for both of them. Refer to the [STYLE guide](https://github.com/Qumulo/terraform-provider-qumulo/blob/dev-docs/STYLE.md#structs) to understand the conventions followed while defining structs.
5. The resource structure and data schema along with the CRUD operations are defined as a function returning a ``schema.Resource`` type. For example,
          
      ```golang
       package qumulo
   
       func resourceExampleSettings() *schema.Resource {
           return &schema.Resource{
               CreateContext: resourceExampleSettingsCreate,
               ReadContext:   resourceExampleSettingsRead,
               UpdateContext: resourceExampleSettingsUpdate,
               DeleteContext: resourceExampleSettingsDelete,
               
               Schema: map[string]*schema.Schema {
                   "name": &schema.Schema{
                         Type: schema.TypeString,
                         Required: true,
                  },
               },
           }
       }
    ```
    The data schema includes the definitions of all (required and optional) the properties required to specify a Terraform resource. Specifically, these are the properties specified in the terraform configuration block. For example, 
    ```terraform
    resource "example_settings" "some_example" {
      #define properties here
      name = "nameMe!"
   }
    ```
6. The four fields ``CreateContext``, ``ReadContext``, ``UpdateContext`` (optional if all the fields are marked ForceNew), and ``DeleteContext`` are mandatory for the management of the resource via Terraform. There are other functions, like ``Importer`` which are optional and have been defined for most/all of Qumulo's resources. Based on the schema and current state of the resource, Terraform determines which of the functions to call. You can refer to any of the existing resource_.go files for what goes inside each of these function definitions!
7. In order to register the new terraform resource, update ``provider.go`` with the new resource.
      ```golang 
       func Provider() *schema.Provider {
           return &schema.Provider{
               ResourcesMap: map[string]*schema.Resource{
                   "example_settings": resourceExampleSettings(),
               },
          }
       }   
       ```
8. To add tests corresponding to the resource, create a new file with the resource name suffixed by ``_test`` for the Go plugin to identify tests to be run as part of ``make test``. For our example, it would be ``resource_example_settings_test.go``.
    Check out any of the resource test files for further details on the structure and implementation details related to testing Qumulo's terraform resources.