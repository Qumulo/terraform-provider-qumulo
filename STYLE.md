# Style Guide for the Qumulo Terraform Provider

We (the initial developers at Qumulo) wrote a style guide to help with development--it's a lot easier to read the code if everything looks the same. We generally aim to follow the style espoused by Terraform and Go (including some forced style like capitalization and Terraform resource naming). Any contributions to this repository are strongly encouraged to follow these guidelines and may be rejected if they do not. (Ideally, the entire repository could be automatically linted to ensure uniformity in the future.)

This guide is always under development! Please let us know of any suggestions/changes/additions to make our codebase even more readable!

## General Naming
- Follow the names specified in the Qumulo REST API (i.e. `Active Directory` instead of `AD`, `LDAP` instead of `Lightweight Directory Access Protocol`)
    - **EXCEPTION**: For constants that represent API endpoints, use the format of the endpoint, omitting API versions. This is usually the same as the name of the section in the REST API, but not always; for example, the constant for `/v1/ad/settings` is `AdSettingsEndpoint` instead of `ActiveDirectorySettingsEndpoint`.
- Enforce MixedCaps or mixedCaps for all names, as is Go style, including with acronyms (i.e. `Ad` or `Ldap` instead of `AD` or `LDAP`).
- Terraform fields should be the same as their JSON equivalents in all lowercase (including underscores).
    - **EXCEPTION**: if there are username or password fields, append `resource_` in the configuration file name to help make it clear what the username or password are used for in the configuration file. Use the short form of the resource name, if it exists. Do not use `resource_` internally, since the field is located in a resource file that is named appropriately. (i.e. there is an `ad_password` field in the configuration file which is internally mapped to a `Password` field in a struct.)

## Variables
- Client objects should be named `c`, not `client`. In general, use the shortest-possible name for common variables such as `d` for a `ResourceData`, etc. (This suggestion follows Terraform's preferred format.)
- Constants, including lists of valid strings, begin with a capital letter.
- Variables with a struct type should have names which mirror the struct name as much as possible. Drop `Body` if it is part of the struct name, and optionally drop the name of the resource if the name is too long (since the file disambiguates).
    - Variables in update functions can optionally prepend `updated` to denote new values
- Variables internal to a function should be `mixedCaps`, not `MixedCaps`.

## Functions
- Enforce Terraform syntax for function names, or `resourceNameCRUD` (i.e. `resourceActiveDirectoryCreate`)
- All helper functions should be `mixedCaps` (begin with a lowercase letter).
- Helper functions should be named `operationResourceDetail`, such as `updateActiveDirectorySettings`.
    - The `Detail` portion of a helper function name should mirror the API endpoint that is hit as much as possible.

## Structs
- Structs should be named `ResourceDetail[Request/Response/Body]` where `Detail` mirrors the usage of the struct (exactly as function names should).
    - Append `Request` or `Response` if the struct is used solely for one API call.
    - Append `Body` if the struct is used for both the request and the response of an API call.
- Structs which are a part of other structs (i.e. not directly exposed to the REST API) should not have an ending of `Request/Response/Body`.
- Fields in a struct should be `MixedCaps` spellings of their JSON equivalents *without* punctuation. Do not omit or shorten any parts of the name.

## Enums
- Only use enums when we need to use a hardcoded value in the code, such as a default value for a field, and the type of the field on the REST API is also enumerated. If the field is just used for parsing, use a normal `string` type with a valid-list (which are the enum values). (Note that you need this valid list for enums anyway, so all you lose is the internal typing.)
- Enums should be named `ResourceField` where `Field` is the name of the field in the struct to which the enum is related, i.e. `ActiveDirectorySigning`. (Note that the actual struct field will have type `string`, since parsing to an Enum takes more work than it's worth to implement.)
- Enum values should be `MixedCaps` spellings of their API equivalents without punctuation.
- Enums start from `iota + 1` to allow for 0 to still be a special value (in our case, an error).
- There exists a variable named `EnumValues` of type `[]string` with all string representations of the enum in the correct order (i.e. `ActiveDirectorySigningValues`). (If you're not using an enum, this is the list you should use, and all you need.)
- The `String()` method for enums should index into `EnumValues`.