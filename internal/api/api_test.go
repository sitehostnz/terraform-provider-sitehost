package api

import (
	"reflect"
	"strconv"
	"testing"
)

// nolint:golint,stylecheck
var APIClient_InfoTests = []struct {
	body    string
	modules []string
	err     bool
}{
	{body: `{"msg":"Unauthorised. The Api key is missing."}`, err: true},
	{body: `{"return":{"client_id":"970180","contact_id":null,"modules":["Job","Server"],"roles":["sitehost"]},"msg":"Successful.","status":true}`, modules: []string{"Job", "Server"}},
}

func TestAPIClient_Info(t *testing.T) {
	for i, tt := range APIClient_InfoTests {
		tt := tt
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			m, err := (*APIClient)(client(tt.body)).Info(ctx)
			switch {
			case tt.err && err == nil:
				t.Errorf("err should not be: <nil>")
			case !tt.err && err != nil:
				t.Errorf("err should be <nil>: %v", err)
			case !reflect.DeepEqual(m, tt.modules):
				t.Errorf("want: %v, got: %v", tt.modules, m)
			}
		})
	}
}
