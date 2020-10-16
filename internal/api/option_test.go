package api

import (
	"net/url"
	"reflect"
	"strconv"
	"testing"
)

//nolint:golint,stylecheck
var ParamIP_setTests = []struct {
	self ParamIP
	in   url.Values
	want url.Values
}{
	{in: url.Values{}, want: url.Values{paramIPv4Key: []string{"auto"}, paramIPv6Key: []string{"auto"}}},
	{self: ParamIP{}, want: url.Values{paramIPv4Key: []string{"auto"}, paramIPv6Key: []string{"auto"}}},
	{in: url.Values{paramIPv4Key: []string{"223.165.66.169"}, paramIPv6Key: []string{"2403:7000:8000:b00::4b"}}, want: url.Values{paramIPv4Key: []string{"auto"}, paramIPv6Key: []string{"auto"}}},
	{self: ParamIP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 223, 165, 66, 169}, want: url.Values{paramIPv4Key: []string{"223.165.66.169"}}},
	{self: ParamIP{36, 3, 112, 0, 128, 0, 11, 0, 0, 0, 0, 0, 0, 0, 0, 75}, want: url.Values{paramIPv6Key: []string{"2403:7000:8000:b00::4b"}}},
	{self: ParamIP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 223, 165, 66, 169}, in: url.Values{paramIPv6Key: []string{"2403:7000:8000:b00::4b"}}, want: url.Values{paramIPv4Key: []string{"223.165.66.169"}, paramIPv6Key: []string{"2403:7000:8000:b00::4b"}}},
	{self: ParamIP{36, 3, 112, 0, 128, 0, 11, 0, 0, 0, 0, 0, 0, 0, 0, 75}, in: url.Values{paramIPv4Key: []string{"223.165.66.169"}}, want: url.Values{paramIPv4Key: []string{"223.165.66.169"}, paramIPv6Key: []string{"2403:7000:8000:b00::4b"}}},
	{self: ParamIP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 223, 165, 66, 168}, in: url.Values{paramIPv4Key: []string{"223.165.66.169"}}, want: url.Values{paramIPv4Key: []string{"223.165.66.169", "223.165.66.168"}}},
	{self: ParamIP{36, 3, 112, 0, 128, 0, 11, 0, 0, 0, 0, 0, 0, 0, 0, 74}, in: url.Values{paramIPv6Key: []string{"2403:7000:8000:b00::4b"}}, want: url.Values{paramIPv6Key: []string{"2403:7000:8000:b00::4b", "2403:7000:8000:b00::4a"}}},
}

func TestParamIP_set(t *testing.T) {
	for i, tt := range ParamIP_setTests {
		tt := tt
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := url.Values{}
			if tt.in != nil {
				got = tt.in
			}
			tt.self.set(got)
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("want: %v, got: %v", tt.want, tt.in)
			}
		})
	}
}
