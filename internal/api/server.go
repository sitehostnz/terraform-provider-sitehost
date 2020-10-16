package api

import (
	"context"
	"errors"
	"net"
	"net/url"
	"time"
)

// ServerClient is the container with methods to call /server/.
type ServerClient Client

// Server vends a client for making calls to /server/.
func (c *Client) Server() *ServerClient { return (*ServerClient)(c) }

// get is a wrapper to prevent inline type conversions.
//
// See: Client.get.
func (s *ServerClient) get(ctx context.Context, path string, data url.Values, keys []string, v interface{}) error {
	return (*Client)(s).get(ctx, path, data, keys, v)
}

// postForm is a wrapper to prevent inline type conversions.
//
// See: Client.postForm.
func (s *ServerClient) postForm(ctx context.Context, path string, data url.Values, keys []string, v interface{}) error {
	return (*Client)(s).postForm(ctx, path, data, keys, v)
}

// ServerDeleteOption is an optional parameter for ServerClient.Delete.
//
// This includes:
//	- ForceDelete
//	- DeleteSubscription
type ServerDeleteOption interface {
	serverDelete()
	callOption
}

func (forceDelete) serverDelete()        {}
func (DeleteSubscription) serverDelete() {}

// Delete a server.
func (s *ServerClient) Delete(ctx context.Context, name string, opts ...ServerDeleteOption) (job uint, err error) {
	data := url.Values{nameKey: []string{name}}
	for _, o := range opts {
		o.set(data)
	}
	var r jobIDResponse
	err = s.postForm(ctx, "/server/delete.json", data, []string{clientIDKey, nameKey, deleteSubscriptionKey, forceDeleteKey}, &r)
	if err != nil {
		return 0, err
	}
	if !r.Status {
		return 0, errors.New(r.Message)
	}
	return r.Return.JobID, nil
}

// Server is metadata about a provisioned sitehost server.
type Server struct {
	Arch             string
	AvailableKernels []struct {
		Default    bool
		Hypervisor string
		Initrd     string
		Kernel     string
		Modules    string
	}
	Backup        bool
	BackupProduct string
	BackupTypes   []interface{}
	ClientID      uint
	Core          uint
	Cores         uint
	Created       time.Time
	Disk          uint
	DiskNew       uint
	Distro        string
	EmailLogs     bool
	GroupID       uint
	ID            uint
	Initrd        string
	Interfaces    []string
	IPAddrLimit   uint
	IPs           []struct {
		AddrFamily uint
		Bridge     string
		Broadcast  net.IP
		Gateway    net.IP
		ID         uint
		IPAddr     net.IP
		IPType     uint
		MACAddr    net.HardwareAddr
		Netmask    net.IP
		Network    net.IP
		NetworkID  uint
		Prefix     uint
		Primary    bool
		RDNS       string
		ServerID   uint
	}
	Kernel  string
	Label   string
	LastJob struct {
		ID    uint
		State JobState
		Type  JobType
	}
	Location     string
	LocationName string
	LocationNode string
	Locked       bool
	MaintDate    time.Time
	MaintDateEnd time.Time
	Managed      bool
	Mirror       bool
	Modules      string
	Name         string
	Notes        string
	OS           string
	Partitions   []struct {
		AlertThreshold float64
		Backup         bool
		Device         string
		DiskTotal      uint
		DiskUsed       uint
		DRBD           bool
		FSType         string
		ID             uint
		InodesTotal    uint
		InodesUsed     uint
		Mountpoint     string
		Name           string
		NewSize        uint
		Size           uint
		Type           string
	}
	ProductCode  string
	ProductName  string
	ProductType  string
	RAM          uint
	Rescue       bool
	Root         string
	State        ServerState
	Subscription struct {
		Code  string
		Name  string
		Price uint
	}
	Type      string
	VNCPort   uint
	VNCScreen VNCScreen
}

type getServerReturn struct {
	Arch             string `json:"arch"`
	AvailableKernels []struct {
		Default    bool   `json:"default"`
		Hypervisor string `json:"hypervisor"`
		Initrd     string `json:"initrd"`
		Kernel     string `json:"kernel"`
		Modules    string `json:"modules"`
	} `json:"available_kernels"`
	BackupTypes   []interface{} `json:"backup_types"`
	Backup        bool          `json:"backups_enabled"`
	BackupProduct string        `json:"backups_product"`
	ClientID      uint          `json:"client_id"`
	Core          uint          `json:"core,string"`
	Cores         uint          `json:"cores,string"`
	Created       stime         `json:"created"`
	Disk          uint          `json:"disk"`
	Distro        string        `json:"distro"`
	EmailLogs     sbool         `json:"email_logs"`
	GroupID       uint          `json:"group_id,string"`
	Initrd        string        `json:"initrd"`
	Interfaces    []string      `json:"interfaces"`
	IPAddrLimit   uint          `json:"ip_addr_limit,string"`
	IPs           []struct {
		AddrFamily uint   `json:"addr_family"`
		Bridge     string `json:"bridge"`
		Broadcast  net.IP `json:"broadcast"`
		Gateway    net.IP `json:"gateway"`
		ID         uint   `json:"id,string"`
		IPAddr     net.IP `json:"ip_addr"`
		IPType     uint   `json:"ip_type,string"`
		MACAddr    smac   `json:"mac_addr"`
		Netmask    net.IP `json:"netmask"`
		Network    net.IP `json:"network"`
		NetworkID  uint   `json:"network_id,string"`
		Prefix     uint   `json:"prefix,string"`
		Primary    bool   `json:"primary"`
		RDNS       string `json:"rdns"`
		ServerID   uint   `json:"server_id,string"`
	} `json:"ips"`
	Kernel  string `json:"kernel"`
	Label   string `json:"label"`
	LastJob struct {
		ID    uint     `json:"id,string"`
		State JobState `json:"state"`
		Type  JobType  `json:"type"`
	} `json:"last_job"`
	Location     string `json:"location_code"`
	LocationName string `json:"location_name"`
	LocationNode string `json:"location_node"`
	Locked       sbool  `json:"locked"`
	MaintDate    stime  `json:"maint_date"`
	MaintDateEnd stime  `json:"maint_date_end"`
	Managed      sbool  `json:"managed"`
	Mirror       sbool  `json:"mirror"`
	Modules      string `json:"modules"`
	Name         string `json:"name"`
	Notes        string `json:"notes"`
	OS           string `json:"os"`
	Partitions   []struct {
		AlertThreshold float64 `json:"alert_threshold,string"`
		Backup         sbool   `json:"backup"`
		Device         string  `json:"device"`
		DiskTotal      uint    `json:"disk_total,string"`
		DiskUsed       uint    `json:"disk_used,string"`
		DRBD           sbool   `json:"drbd"`
		FSType         string  `json:"fstype"`
		ID             uint    `json:"id,string"`
		InodesTotal    uint    `json:"inodes_total,string"`
		InodesUsed     uint    `json:"inodes_used,string"`
		Mountpoint     string  `json:"mountpoint"`
		Name           string  `json:"name"`
		NewSize        uint    `json:"new_size,string"`
		Size           uint    `json:"size,string"` // GB?
		Type           string  `json:"type"`
	} `json:"partitions"`
	ProductCode  string      `json:"product_code"`
	ProductName  string      `json:"product_name"`
	ProductType  string      `json:"product_type"`
	RAM          uint        `json:"ram,string"` // MB
	Rescue       sbool       `json:"rescue"`
	Root         string      `json:"root"`
	State        ServerState `json:"state"`
	Subscription struct {
		Code  string `json:"code"`
		Name  string `json:"name"`
		Price uint   `json:"price,string"`
	} `json:"subscription"`
	Type      string    `json:"type"`
	VNCPort   uint      `json:"vnc_port,string"`
	VNCScreen VNCScreen `json:"vnc_screen"`
}

func (r getServerReturn) into() Server {
	s := Server{
		Arch:          r.Arch,
		BackupTypes:   r.BackupTypes,
		Backup:        r.Backup,
		BackupProduct: r.BackupProduct,
		ClientID:      r.ClientID,
		Core:          r.Core,
		Cores:         r.Cores,
		Created:       time.Time(r.Created),
		Disk:          r.Disk,
		Distro:        r.Distro,
		EmailLogs:     bool(r.EmailLogs),
		GroupID:       r.GroupID,
		Initrd:        r.Initrd,
		Interfaces:    r.Interfaces,
		IPAddrLimit:   r.IPAddrLimit,
		Kernel:        r.Kernel,
		Label:         r.Label,
		Location:      r.Location,
		LocationName:  r.LocationName,
		LocationNode:  r.LocationNode,
		Locked:        bool(r.Locked),
		MaintDate:     time.Time(r.MaintDate),
		MaintDateEnd:  time.Time(r.MaintDateEnd),
		Managed:       bool(r.Managed),
		Mirror:        bool(r.Mirror),
		Modules:       r.Modules,
		Name:          r.Name,
		Notes:         r.Notes,
		OS:            r.OS,
		ProductCode:   r.ProductCode,
		ProductName:   r.ProductName,
		ProductType:   r.ProductType,
		RAM:           r.RAM,
		Rescue:        bool(r.Rescue),
		Root:          r.Root,
		State:         r.State,
		Type:          r.Type,
		VNCPort:       r.VNCPort,
		VNCScreen:     r.VNCScreen,
	}
	s.LastJob.ID = r.LastJob.ID
	s.LastJob.State = r.LastJob.State
	s.LastJob.Type = r.LastJob.Type
	s.Subscription.Code = r.Subscription.Code
	s.Subscription.Name = r.Subscription.Name
	s.Subscription.Price = r.Subscription.Price
	for _, v := range r.AvailableKernels {
		s.AvailableKernels = append(s.AvailableKernels, (struct {
			Default    bool
			Hypervisor string
			Initrd     string
			Kernel     string
			Modules    string
		})(v))
	}
	for _, v := range r.IPs {
		s.IPs = append(s.IPs, struct {
			AddrFamily uint
			Bridge     string
			Broadcast  net.IP
			Gateway    net.IP
			ID         uint
			IPAddr     net.IP
			IPType     uint
			MACAddr    net.HardwareAddr
			Netmask    net.IP
			Network    net.IP
			NetworkID  uint
			Prefix     uint
			Primary    bool
			RDNS       string
			ServerID   uint
		}{
			v.AddrFamily,
			v.Bridge,
			v.Broadcast,
			v.Gateway,
			v.ID,
			v.IPAddr,
			v.IPType,
			net.HardwareAddr(v.MACAddr),
			v.Netmask,
			v.Network,
			v.NetworkID,
			v.Prefix,
			v.Primary,
			v.RDNS,
			v.ServerID,
		})
	}
	for _, v := range r.Partitions {
		s.Partitions = append(s.Partitions, struct {
			AlertThreshold float64
			Backup         bool
			Device         string
			DiskTotal      uint
			DiskUsed       uint
			DRBD           bool
			FSType         string
			ID             uint
			InodesTotal    uint
			InodesUsed     uint
			Mountpoint     string
			Name           string
			NewSize        uint
			Size           uint
			Type           string
		}{
			v.AlertThreshold,
			bool(v.Backup),
			v.Device,
			v.DiskTotal,
			v.DiskUsed,
			bool(v.DRBD),
			v.FSType,
			v.ID,
			v.InodesTotal,
			v.InodesUsed,
			v.Mountpoint,
			v.Name,
			v.NewSize,
			v.Size,
			v.Type,
		})
	}
	return s
}

// The requested resource does not exist.
var ErrNotFound = errors.New("not found")

// Retrieves the details for a server.
func (s *ServerClient) Get(ctx context.Context, name string) (Server, error) {
	if name == "" {
		return Server{}, errors.New("The server name is missing.")
	}
	var r struct {
		Return  getServerReturn `json:"return"`
		Message string          `json:"msg"`
		Status  bool            `json:"status"`
	}
	switch err := s.get(ctx, "/server/get_server.json", url.Values{nameKey: []string{name}}, []string{clientIDKey, nameKey}, &r); {
	case err != nil:
		return Server{}, err
	case r.Message == "Error: Not Found":
		return Server{}, ErrNotFound
	case !r.Status:
		return Server{}, errors.New(r.Message)
	default:
		return r.Return.into(), nil
	}
}

type listServersData []struct {
	Arch          string `json:"arch"`
	Backup        bool   `json:"backups_enabled"`
	BackupProduct string `json:"backups_product"`
	ClientID      uint   `json:"client_id,string"`
	Cores         uint   `json:"cores"`
	Created       stime  `json:"created"`
	Disk          uint   `json:"disk,string"`
	DiskNew       uint   `json:"disk_new,string"`
	Distro        string `json:"distro"`
	ID            uint   `json:"server_id,string"`
	Label         string `json:"label"`
	Location      string `json:"location"`
	LocationName  string `json:"location_name"`
	Locked        sbool  `json:"locked"`
	MaintDate     stime  `json:"maint_date"`
	MaintDateEnd  stime  `json:"maint_date_end"`
	Managed       sbool  `json:"managed"`
	Mirror        sbool  `json:"mirror"`
	Name          string `json:"name"`
	OS            string `json:"os"`
	Pending       sbool  `json:"pending"`
	PrimaryIPs    []struct {
		IPAddr net.IP `json:"ip_addr"`
		Prefix uint   `json:"prefix,string"`
	} `json:"primary_ips"`
	ProductCode string      `json:"product_code"`
	ProductName string      `json:"product_name"`
	ProductType string      `json:"product_type"`
	RAM         float64     `json:"ram,string"`
	Rescue      sbool       `json:"rescue"`
	State       ServerState `json:"state"`
	Type        string      `json:"type"`
}

func (d listServersData) into() []Server {
	var ss []Server
	for _, v := range d {
		s := Server{
			Arch:          v.Arch,
			Backup:        v.Backup,
			BackupProduct: v.BackupProduct,
			ClientID:      v.ClientID,
			Cores:         v.Cores,
			Created:       time.Time(v.Created),
			Disk:          v.Disk,
			DiskNew:       v.DiskNew,
			Distro:        v.Distro,
			ID:            v.ID,
			Label:         v.Label,
			Location:      v.Location,
			LocationName:  v.LocationName,
			Locked:        bool(v.Locked),
			MaintDate:     time.Time(v.MaintDate),
			MaintDateEnd:  time.Time(v.MaintDateEnd),
			Managed:       bool(v.Managed),
			Mirror:        bool(v.Mirror),
			Name:          v.Name,
			OS:            v.OS,
			ProductCode:   v.ProductCode,
			ProductName:   v.ProductName,
			ProductType:   v.ProductType,
			RAM:           uint(v.RAM * 1024),
			Rescue:        bool(v.Rescue),
			State:         v.State,
			Type:          v.Type,
		}
		if v.Pending {
			s.LastJob.State = PendingJob
		}
		for _, v := range v.PrimaryIPs {
			s.IPs = append(s.IPs, struct {
				AddrFamily uint
				Bridge     string
				Broadcast  net.IP
				Gateway    net.IP
				ID         uint
				IPAddr     net.IP
				IPType     uint
				MACAddr    net.HardwareAddr
				Netmask    net.IP
				Network    net.IP
				NetworkID  uint
				Prefix     uint
				Primary    bool
				RDNS       string
				ServerID   uint
			}{
				IPAddr:  v.IPAddr,
				Prefix:  v.Prefix,
				Primary: true,
			})
		}
		ss = append(ss, s)
	}
	return ss
}

// ServerListOption is an optional parameter for ServerClient.List.
//
// This includes:
//	- ServerState
//	- FilterType
//	- FilterName
//	- FilterLocation
//	- FilterProductType
//	- FilterProductCode
type ServerListOption interface {
	serverList()
	callOption
}

func (ServerState) serverList()       {}
func (FilterType) serverList()        {}
func (FilterName) serverList()        {}
func (FilterLocation) serverList()    {}
func (FilterProductType) serverList() {}
func (FilterProductCode) serverList() {}

// Lists all servers for a client.
func (s *ServerClient) List(ctx context.Context, opts ...ServerListOption) ([]Server, error) {
	data := url.Values{}
	for _, o := range opts {
		o.set(data)
	}
	var ss []Server
	err := iter(ctx, func(ctx context.Context, opts ...callOption) (uint, error) {
		for _, o := range opts {
			o.set(data)
		}
		var r struct {
			Return struct {
				TotalPages uint            `json:"total_pages"`
				Data       listServersData `json:"data"`
			} `json:"return"`
			Message string `json:"msg"`
			Status  bool   `json:"status"`
		}
		err := s.get(ctx, "/server/list_servers.json", data, []string{clientIDKey, filterStateKey, filterTypeKey, filterNameKey, filterLocationKey, filterProductTypeKey, filterSortByKey, filterSortDirKey, filterPageSizeKey, filterPageNumberKey}, &r)
		if err != nil {
			return 0, err
		}
		if !r.Status {
			return 0, errors.New(r.Message)
		}
		ss = append(ss, r.Return.Data.into()...)
		return r.Return.TotalPages, nil
	})
	return ss, err
}

// ServerProvisionOption is an optional parameter for ServerClient.List.
//
// This includes:
//	- ParamName
//	- ParamIP
//	- ParamSSHKey
//	- ParamContactID
//	- ParamBackup
//	- ParamSendEmail
type ServerProvisionOption interface {
	serverProvision()
	callOption
}

func (ParamName) serverProvision()      {}
func (ParamIP) serverProvision()        {}
func (ParamSSHKey) serverProvision()    {}
func (ParamContactID) serverProvision() {}
func (paramBackup) serverProvision()    {}
func (paramSendEmail) serverProvision() {}

// Provisions a new server.
func (s *ServerClient) Provision(ctx context.Context, label, location, product, image string, opts ...ServerProvisionOption) (id, job uint, name, password string, addr []net.IP, err error) {
	data := url.Values{
		labelKey:       []string{label},
		locationKey:    []string{location},
		productCodeKey: []string{product},
		imageKey:       []string{image},
	}
	for _, o := range opts {
		o.set(data)
	}
	var r struct {
		Return struct {
			JobID    uint     `json:"job_id,string"`
			Name     string   `json:"name"`
			Password string   `json:"password"`
			IPs      []net.IP `json:"ips"`
			ServerID uint     `json:"server_id,string"`
		} `json:"return"`
		Message string `json:"msg"`
		Status  bool   `json:"status"`
	}
	err = s.postForm(ctx, "/server/provision.json", data, []string{clientIDKey, labelKey, locationKey, productCodeKey, imageKey, paramNameKey, paramIPv4Key, paramIPv6Key, paramSSHKey, paramContactIDKey, paramBackupKey, paramSendEmailKey}, &r)
	if err != nil {
		return 0, 0, "", "", nil, err
	}
	if !r.Status {
		return 0, 0, "", "", nil, errors.New(r.Message)
	}
	return r.Return.ServerID, r.Return.JobID, r.Return.Name, r.Return.Password, r.Return.IPs, nil
}
