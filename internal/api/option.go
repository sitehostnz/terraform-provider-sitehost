package api

import (
	"net"
	"net/url"
)

type callOption interface{ set(url.Values) }

type filterPage uint

func (p filterPage) set(v url.Values) { v.Set(filterPageNumberKey, string(p)) }

type Arch string

func (a Arch) set(v url.Values) { v.Set(archKey, string(a)) }

type forceDelete struct{}

var ForceDelete = forceDelete{}

func (forceDelete) set(v url.Values) { v.Set(forceDeleteKey, "1") }

type DeleteSubscription string

func (s DeleteSubscription) set(v url.Values) { v.Set(deleteSubscriptionKey, string(s)) }

const (
	FilterOn    FilterState = "On"
	FilterOff   FilterState = "Off"
	FilterOther FilterState = "Other"
)

type FilterState string

func (s FilterState) set(v url.Values) { v.Set(filterStateKey, string(s)) }

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

type ParamIPv4 net.IP

func (i ParamIPv4) set(v url.Values) {
	if ii := net.IP(i).To4(); ii != nil {
		v.Add(paramIPv4Key, ii.String())
	}
}

type ParamIPv6 net.IP

func (i ParamIPv6) set(v url.Values) {
	if ii := net.IP(i).To16(); ii != nil {
		v.Add(paramIPv6Key, ii.String())
	}
}

type ParamSSHKey string

func (k ParamSSHKey) set(v url.Values) { v.Add(paramSSHKey, string(k)) }

type ParamContactID uint

func (c ParamContactID) set(v url.Values) { v.Set(paramContactIDKey, string(c)) }

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

type UpdateVNC struct {
	X uint
	Y uint
}

func (c UpdateVNC) set(v url.Values) { v.Set(updateVNCKey, utoa(c.X)+"x"+utoa(c.Y)) }

type UpdateKernel string

func (k UpdateKernel) set(v url.Values) { v.Set(updateKernelKey, string(k)) }

/*
type UpdatePartitionThreshold struct {
	Name string
	Size uint
}

func (t UpdatePartitionThreshold) set(v url.Values) {
	v.Set("updates[partitions]["+t.Name+"][threshold]", string(t.Size))
}
*/

type UpgradeCores uint

func (c UpgradeCores) set(v url.Values) { v.Set(upgradeCoresKey, utoa(uint(c))) }

type UpgradeRAM uint

func (r UpgradeRAM) set(v url.Values) { v.Set(upgradeRAMKey, utoa(uint(r))) }

/*
type UpgradeDisk struct {
	Name string
	Size uint
}

func (d UpgradeDisk) set(v url.Values) { v.Set("upgrade[disk]["+d.Name+"]", string(d.Size)) }
*/
