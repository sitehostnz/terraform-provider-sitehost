package api

import (
	"net"
	"strconv"
	"time"
)

// atou is equivalent to ParseUint(s, 10, 0), converted to type uint.
//
// It is based on strconv.Atoi.
func atou(s string) (uint, error) {
	u, err := strconv.ParseUint(s, 10, 0)
	return uint(u), err
}

// stime is a time.Time that can unmarshal the format provided by the api.
//
// Specifically, it maps "0000-00-00 00:00:00" to time.Time{}, something
// time.Parse does not do, and it supports the custom format used by the
// server: "2006-01-02 15:04:05".
type stime time.Time

// UnmarshalJSON implements json.Unmarshaler.
func (t *stime) UnmarshalJSON(b []byte) error {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	if s == "0000-00-00 00:00:00" {
		return nil
	}
	tt, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return err
	}
	*t = stime(tt)
	return nil
}

// smac is a net.HardwareAddr with UnmarshalJSON support.
//
// It leverages net.ParseMAC and may be sunset if golang/go#29678 is resolved.
type smac net.HardwareAddr

// UnmarshalJSON implements json.Unmarshaler.
func (m *smac) UnmarshalJSON(b []byte) error {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	a, err := net.ParseMAC(s)
	*m = smac(a)
	return err
}

// sbool is a bool derived from a string wrapped number.
//
// Invalid responses yield the default, false, and only "1" can make true.
type sbool bool

// UnmarshalJSON implements json.Unmarshaler.
func (s *sbool) UnmarshalJSON(b []byte) error { *s = string(b) == `"1"`; return nil }
