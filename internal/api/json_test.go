package api

import (
	"encoding/json"
	"net"
	"reflect"
	"strconv"
	"testing"
	"time"
)

// utc generates a time.Time in the UTC timezone from the unix time provided.
func utc(unix int64) time.Time { return time.Unix(unix, 0).UTC() }

//nolint:golint,stylecheck
var smac_UnmarshalJSONTests = []struct {
	in   string
	want net.HardwareAddr
	err  bool
}{
	{err: true},
	{in: `"00:00:5e:00:53:01"`, want: []byte{0x00, 0x00, 0x5e, 0x00, 0x53, 0x01}},
	{in: `"02:00:5e:10:00:00:00:01"`, want: []byte{0x02, 0x00, 0x5e, 0x10, 0x00, 0x00, 0x00, 0x01}},
	{in: `"00:00:00:00:fe:80:00:00:00:00:00:00:02:00:5e:10:00:00:00:01"`, want: []byte{0x00, 0x00, 0x00, 0x00, 0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x5e, 0x10, 0x00, 0x00, 0x00, 0x01}},
	{in: `"00-00-5e-00-53-01"`, want: []byte{0x00, 0x00, 0x5e, 0x00, 0x53, 0x01}},
	{in: `"02-00-5e-10-00-00-00-01"`, want: []byte{0x02, 0x00, 0x5e, 0x10, 0x00, 0x00, 0x00, 0x01}},
	{in: `"00-00-00-00-fe-80-00-00-00-00-00-00-02-00-5e-10-00-00-00-01"`, want: []byte{0x00, 0x00, 0x00, 0x00, 0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x5e, 0x10, 0x00, 0x00, 0x00, 0x01}},
	{in: `"0000.5e00.5301"`, want: []byte{0x00, 0x00, 0x5e, 0x00, 0x53, 0x01}},
	{in: `"0200.5e10.0000.0001"`, want: []byte{0x02, 0x00, 0x5e, 0x10, 0x00, 0x00, 0x00, 0x01}},
	{in: `"0000.0000.fe80.0000.0000.0000.0200.5e10.0000.0001"`, want: []byte{0x00, 0x00, 0x00, 0x00, 0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x5e, 0x10, 0x00, 0x00, 0x00, 0x01}},
}

func TestSmac_UnmarshalJSON(t *testing.T) {
	for i, tt := range smac_UnmarshalJSONTests {
		tt := tt
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got net.HardwareAddr
			err := json.Unmarshal([]byte(tt.in), (*smac)(&got))
			switch {
			case tt.err && err == nil:
				t.Errorf("err should not be: <nil>")
			case !tt.err && err != nil:
				t.Errorf("err should be <nil>: %v", err)
			case !reflect.DeepEqual(tt.want, got):
				t.Errorf("want: %v, got: %v", tt.want, []byte(got))
			}
		})
	}
}

//nolint:golint,stylecheck
var sbool_UnmarshalJSONTests = []struct {
	in   string
	want bool
	err  bool
}{
	{err: true},
	{in: `"0"`},
	{in: `"1"`, want: true},
}

func TestSbool_UnmarshalJSON(t *testing.T) {
	for i, tt := range sbool_UnmarshalJSONTests {
		tt := tt
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got bool
			err := json.Unmarshal([]byte(tt.in), (*sbool)(&got))
			switch {
			case tt.err && err == nil:
				t.Errorf("err should not be: <nil>")
			case !tt.err && err != nil:
				t.Errorf("err should be <nil>: %v", err)
			case tt.want != got:
				t.Errorf("want: %v, got: %v", tt.want, got)
			}
		})
	}
}

//nolint:golint,stylecheck
var stime_UnmarshalJSONTests = []struct {
	in   string
	want time.Time
	err  bool
}{
	{err: true},
	{in: `"0000-00-00 00:00:00"`},
	{in: `"2006-01-02 15:04:05"`, want: utc(1136214245)},
}

func TestStime_UnmarshalJSON(t *testing.T) {
	for i, tt := range stime_UnmarshalJSONTests {
		tt := tt
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got time.Time
			err := json.Unmarshal([]byte(tt.in), (*stime)(&got))
			switch {
			case tt.err && err == nil:
				t.Errorf("err should not be: <nil>")
			case !tt.err && err != nil:
				t.Errorf("err should be <nil>: %v", err)
			case tt.want != got:
				t.Errorf("want: %v, got: %v", tt.want, got)
			}
		})
	}
}
