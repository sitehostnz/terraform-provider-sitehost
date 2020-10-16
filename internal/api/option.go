package api

import (
	"errors"
	"net"
	"net/url"
	"strconv"
	"strings"
)

type callOption interface{ set(url.Values) }

type filterPage uint

func (p filterPage) set(v url.Values) {
	if p < 2 {
		return
	}
	v.Set(filterPageNumberKey, utoa(uint(p)))
}

type Arch string

func (a Arch) set(v url.Values) { v.Set(archKey, string(a)) }

type forceDelete struct{}

var ForceDelete = forceDelete{}

func (forceDelete) set(v url.Values) { v.Set(forceDeleteKey, "1") }

type DeleteSubscription string

func (s DeleteSubscription) set(v url.Values) { v.Set(deleteSubscriptionKey, string(s)) }

const (
	ServerOn    ServerState = "On"
	ServerOff   ServerState = "Off"
	ServerOther ServerState = "Other"
)

// ServerState is a state that a server can be in.
type ServerState string

func (s ServerState) set(v url.Values) { v.Set(filterStateKey, string(s)) }

const (
	FilterTypeDistro        FilterType = "distro"
	FilterTypeSaltContainer FilterType = "salt-container"
)

type FilterType string

func (t FilterType) set(v url.Values) { v.Set(filterTypeKey, string(t)) }

type FilterName string

func (n FilterName) set(v url.Values) { v.Set(filterNameKey, string(n)) }

type FilterLocation string

func (l FilterLocation) set(v url.Values) { v.Set(filterLocationKey, string(l)) }

const (
	FilterProductLinux            FilterProductType = "LINVPS"
	FilterProductWindows          FilterProductType = "WINVPS"
	FilterProductCloudHosting     FilterProductType = "CLHOST"
	FilterProductCloudContainer   FilterProductType = "CLDCON"
	FilterProductVirtualDedicated FilterProductType = "VDSERV"
	FilterProductDedicated        FilterProductType = "SERVER"
	FilterProductColocated        FilterProductType = "COLO"
	FilterProductEnterpriseCloud  FilterProductType = "CLENT"
	FilterProductPrivateCloud     FilterProductType = "PCLOUD"
)

type FilterProductType string

func (p FilterProductType) set(v url.Values) { v.Set(filterProductTypeKey, string(p)) }

type FilterProductCode string

func (p FilterProductCode) set(v url.Values) { v.Set(filterProductCodeKey, string(p)) }

const (
	FilterOSLinux   FilterOS = "linux"
	FilterOSWindows FilterOS = "windows"
)

type FilterOS string

func (o FilterOS) set(v url.Values) { v.Set(filterOSKey, string(o)) }

type ParamName string

func (n ParamName) set(v url.Values) { v.Set(paramNameKey, string(n)) }

// ParamIP configures ipv4 and ipv6 addresses.
//
//	ipv4 maps to params[ipv4]
//	ipv6 maps to params[ipv6]
type ParamIP net.IP

func (i ParamIP) set(v url.Values) {
	if v.Get(paramIPv4Key) == "auto" {
		v.Del(paramIPv4Key)
	}
	if v.Get(paramIPv6Key) == "auto" {
		v.Del(paramIPv6Key)
	}
	if i := net.IP(i).To4(); i != nil {
		v.Add(paramIPv4Key, i.String())
		return
	}
	if i := net.IP(i).To16(); i != nil {
		v.Add(paramIPv6Key, i.String())
		return
	}
	v.Set(paramIPv4Key, "auto")
	v.Set(paramIPv6Key, "auto")
}

type ParamSSHKey string

func (k ParamSSHKey) set(v url.Values) { v.Add(paramSSHKey, string(k)) }

type ParamContactID uint

func (c ParamContactID) set(v url.Values) { v.Set(paramContactIDKey, utoa(uint(c))) }

var ParamBackup = paramBackup{}

type paramBackup struct{}

func (paramBackup) set(v url.Values) { v.Set(paramBackupKey, "1") }

var ParamSendEmail = paramSendEmail{}

type paramSendEmail struct{}

func (paramSendEmail) set(v url.Values) { v.Set(paramSendEmailKey, "1") }

type UpdateLabel string

func (l UpdateLabel) set(v url.Values) { v.Set(updateLabelKey, string(l)) }

type UpdateNote string

func (n UpdateNote) set(v url.Values) { v.Set(updateNoteKey, string(n)) }

type VNCScreen struct{ X, Y uint }

func (c VNCScreen) set(v url.Values) { v.Set(updateVNCKey, utoa(c.X)+"x"+utoa(c.Y)) }

func (s *VNCScreen) UnmarshalJSON(b []byte) error {
	s0, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	s1 := strings.Split(s0, "x")
	if len(s1) != 2 {
		return errors.New("not NxM format: " + s0)
	}
	x, err := atou(s1[0])
	if err != nil {
		return err
	}
	y, err := atou(s1[1])
	if err != nil {
		return err
	}
	s.X, s.Y = x, y
	return nil
}

type UpdateKernel string

func (k UpdateKernel) set(v url.Values) { v.Set(updateKernelKey, string(k)) }

type UpgradeCores uint

func (c UpgradeCores) set(v url.Values) { v.Set(upgradeCoresKey, utoa(uint(c))) }

type UpgradeRAM uint

func (r UpgradeRAM) set(v url.Values) { v.Set(upgradeRAMKey, utoa(uint(r))) }
