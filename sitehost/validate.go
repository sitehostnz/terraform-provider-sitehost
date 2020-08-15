package sitehost

import "errors"

// unsigned validates that an incoming schema.TypeInt is positive and can be
// safely mapped to a uint for the api client.
//
// unsigned is a schema.SchemaValidateFunc.
func unsigned(v interface{}, k string) ([]string, []error) {
	if v.(int) < 0 {
		return nil, []error{errors.New(k + ": must be a natural number")}
	}
	return nil, nil
}
