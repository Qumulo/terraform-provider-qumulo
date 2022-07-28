package qumulo

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// Wrapper for diag.Diagnostics to provide some useful methods
type ErrorCollection struct {
	diags diag.Diagnostics
}

func (coll *ErrorCollection) addMaybeError(err error) {
	if err != nil {
		coll.diags = append(coll.diags, diag.FromErr(err)...)
	}
}
