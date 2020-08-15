package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"
)

type ServerService Client

func (c *Client) Server() *ServerService { return (*ServerService)(c) }

func (s *ServerService) get(ctx context.Context, path string, data url.Values, keys []string, v interface{}) error {
	return (*Client)(s).get(ctx, path, data, keys, v)
}

func (s *ServerService) postForm(ctx context.Context, path string, data url.Values, keys []string, v interface{}) error {
	return (*Client)(s).postForm(ctx, path, data, keys, v)
}

// AddIP to a server.
func (s *ServerService) AddIP(ctx context.Context, name string, param net.IP) (net.IP, error) {
	var r struct {
		Return struct {
			Addr net.IP `json:"ip_addr"`
		} `json:"return"`
		Message string `json:"msg"`
		Status  bool   `json:"status"`
	}
	err := s.postForm(ctx, "/server/add_ip.json", url.Values{nameKey: []string{name}, paramKey: []string{param.String()}}, []string{clientIDKey, nameKey, paramKey}, &r)
	if err != nil {
		return nil, err
	}
	if !r.Status {
		return nil, errors.New(r.Message)
	}
	return r.Return.Addr, nil
}

type ServerCanProvisionOption interface {
	serverCanProvision()
	callOption
}

func (Arch) serverCanProvision() {}

// CanProvision checks if there are available resources to provision a specific server product.
func (s *ServerService) CanProvision(ctx context.Context, product, location, distro string, opts ...ServerProvisionOption) error {
	data := url.Values{productKey: []string{product}, locationKey: []string{location}, distroKey: []string{distro}}
	for _, o := range opts {
		o.set(data)
	}
	var r voidResponse
	err := s.postForm(ctx, "/server/can_provision.json", data, []string{productKey, locationKey, distroKey, archKey, clientIDKey}, &r)
	if err != nil {
		return err
	}
	if !r.Status {
		return errors.New(r.Message)
	}
	return nil
}

type ServerState string

const (
	PowerOn   ServerState = "power_on"
	PowerOff  ServerState = "power_off"
	RescueOn  ServerState = "rescue_on"
	RescueOff ServerState = "rescue_off"
	Reboot    ServerState = "reboot"
)

// ChangesState of the named server.
func (s *ServerService) ChangeState(ctx context.Context, name string, t ServerState) (job uint, err error) {
	var r jobIDResponse
	err = s.postForm(ctx, "/server/change_state.json", url.Values{nameKey: []string{name}, "state": []string{string(t)}}, []string{clientIDKey, nameKey, stateKey}, &r)
	if err != nil {
		return 0, err
	}
	if !r.Status {
		return 0, errors.New(r.Message)
	}
	return r.Return.JobID, nil
}

/*
// Applies upgrade changes to disks for a server.
func (s *ServerService) CommitDiskChanges(headers, queryParams map[string]interface{}) (types.ServerCommit_disk_changes_JsonPostRespBody, error) {
	resp, err := s.client().doReqWithBody("POST", s.client().BaseURI+"/server/commit_disk_changes.json", nil, headers, queryParams)
}
*/

type ServerDeleteOption interface {
	serverDelete()
	callOption
}

func (forceDelete) serverDelete()        {}
func (DeleteSubscription) serverDelete() {}

// Delete a server.
func (s *ServerService) Delete(ctx context.Context, name string, opts ...ServerDeleteOption) (job uint, err error) {
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

/*
// Generates the network configuration for a server.
func (s *ServerService) GenerateNetworkConfiguration(headers, queryParams map[string]interface{}) (types.ServerGenerate_network_config_JsonGetRespBody, error) {
	resp, err := s.client().doReqNoBody("GET", s.client().BaseURI+"/server/generate_network_config.json", headers, queryParams)
}

// Generates an access token to connect to the VNC console for a server.
func (s *ServerService) GenerateVNCToken(headers, queryParams map[string]interface{}) (types.ServerGenerate_vnc_token_JsonGetRespBody, error) {
	resp, err := s.client().doReqNoBody("GET", s.client().BaseURI+"/server/generate_vnc_token.json", headers, queryParams)
}
*/

type Server struct {
	Arch             string `json:"arch"`
	AvailableKernels []struct {
		Default    bool   `json:"default"`
		Hypervisor string `json:"hypervisor"`
		Initrd     string `json:"initrd"`
		Kernel     string `json:"kernel"`
		Modules    string `json:"modules"`
	} `json:"available_kernels"`
	BackupTypes []interface{} `json:"backup_types"`
	ClientID    uint          `json:"client_id"`
	Core        uint          `json:"core"`
	Cores       uint          `json:"cores"`
	Created     time.Time     `json:"created"`
	Disk        uint          `json:"disk"`
	Distro      string        `json:"distro"`
	EmailLogs   bool          `json:"email_logs"`
	GroupID     uint          `json:"group_id"`
	Initrd      string        `json:"initrd"`
	IPAddrLimit uint          `json:"ip_addr_limit"`
	IPs         []struct {
		AddrFamily uint             `json:"addr_family"`
		Bridge     string           `json:"bridge"`
		Broadcast  net.IP           `json:"broadcast"`
		Gateway    net.IP           `json:"gateway"`
		ID         uint             `json:"id"`
		IPAddr     net.IP           `json:"ip_addr"`
		IPType     uint             `json:"ip_type"`
		MACAddr    net.HardwareAddr `json:"mac_addr"`
		Netmask    net.IP           `json:"netmask"`
		Network    net.IP           `json:"network"`
		NetworkID  uint             `json:"network_id"`
		Prefix     uint             `json:"prefix"`
		Primary    bool             `json:"primary"`
		RDNS       string           `json:"rdns"`
		ServerID   uint             `json:"server_id"`
	} `json:"ips"`
	Kernel       string    `json:"kernel"`
	Label        string    `json:"label"`
	LocationCode string    `json:"location_code"`
	LocationName string    `json:"location_name"`
	Locked       bool      `json:"locked"`
	MaintDate    time.Time `json:"maint_date"`
	MaintDateEnd time.Time `json:"maint_date_end"`
	Managed      bool      `json:"managed"`
	Mirror       bool      `json:"mirror"`
	Modules      string    `json:"modules"`
	Name         string    `json:"name"`
	Notes        string    `json:"notes"`
	OS           string    `json:"os"`
	Partitions   []struct {
		AlertThreshold float64 `json:"alert_threshold"`
		Backup         bool    `json:"backup"`
		Device         string  `json:"device"`
		DiskTotal      uint    `json:"disk_total"`
		DiskUsed       uint    `json:"disk_used"`
		DRBD           bool    `json:"drbd"`
		FSType         string  `json:"fstype"`
		ID             uint    `json:"id"`
		InodesTotal    uint    `json:"inodes_total"`
		InodesUsed     uint    `json:"inodes_used"`
		Mountpoint     string  `json:"mountpoint"`
		Name           string  `json:"name"`
		NewSize        uint    `json:"new_size"`
		Size           uint    `json:"size"`
		Type           string  `json:"type"`
	} `json:"partitions"`
	ProductCode  string `json:"product_code"`
	ProductName  string `json:"product_name"`
	ProductType  string `json:"product_type"`
	RAM          uint   `json:"ram"`
	Rescue       bool   `json:"rescue"`
	Root         string `json:"root"`
	State        string `json:"state"`
	Subscription struct {
		Code  string `json:"code"`
		Name  string `json:"name"`
		Price uint   `json:"price"`
	} `json:"subscription"`
	Type      string `json:"type"`
	VNCPort   uint   `json:"vnc_port"`
	VNCScreen string `json:"vnc_screen"`
}

func (s *Server) UnmarshalJSON(b []byte) error {
	var v struct {
		Arch             *string `json:"arch"`
		AvailableKernels *[]struct {
			Default    bool   `json:"default"`
			Hypervisor string `json:"hypervisor"`
			Initrd     string `json:"initrd"`
			Kernel     string `json:"kernel"`
			Modules    string `json:"modules"`
		} `json:"available_kernels"`
		BackupTypes *[]interface{} `json:"backup_types"`
		ClientID    interface{}    `json:"client_id"`
		Core        *uint          `json:"core,string"`
		Cores       interface{}    `json:"cores"`
		Created     string         `json:"created"`
		Disk        *uint          `json:"disk"`
		Distro      *string        `json:"distro"`
		EmailLogs   uint           `json:"email_logs,string"`
		GroupID     *uint          `json:"group_id,string"`
		Initrd      *string        `json:"initrd"`
		IPAddrLimit *uint          `json:"ip_addr_limit,string"`
		IPs         *[]struct {
			AddrFamily uint             `json:"addr_family"`
			Bridge     string           `json:"bridge"`
			Broadcast  net.IP           `json:"broadcast"`
			Gateway    net.IP           `json:"gateway"`
			ID         uint             `json:"id"`
			IPAddr     net.IP           `json:"ip_addr"`
			IPType     uint             `json:"ip_type"`
			MACAddr    net.HardwareAddr `json:"mac_addr"`
			Netmask    net.IP           `json:"netmask"`
			Network    net.IP           `json:"network"`
			NetworkID  uint             `json:"network_id"`
			Prefix     uint             `json:"prefix"`
			Primary    bool             `json:"primary"`
			RDNS       string           `json:"rdns"`
			ServerID   uint             `json:"server_id"`
		} `json:"ips"`
		Kernel       *string     `json:"kernel"`
		Label        *string     `json:"label"`
		LocationCode *string     `json:"location_code"`
		LocationName *string     `json:"location_name"`
		Locked       uint        `json:"locked,string"`
		MaintDate    string      `json:"maint_date"`
		MaintDateEnd string      `json:"maint_date_end"`
		Managed      uint        `json:"managed,string"`
		Mirror       uint        `json:"mirror,string"`
		Modules      *string     `json:"modules"`
		Name         *string     `json:"name"`
		Notes        *string     `json:"notes"`
		OS           *string     `json:"os"`
		Partitions   interface{} `json:"partitions"`
		ProductCode  *string     `json:"product_code"`
		ProductName  *string     `json:"product_name"`
		ProductType  *string     `json:"product_type"`
		RAM          *uint       `json:"ram,string"`
		Rescue       uint        `json:"rescue,string"`
		Root         *string     `json:"root"`
		State        *string     `json:"state"`
		Subscription *struct {
			Code  string `json:"code"`
			Name  string `json:"name"`
			Price uint   `json:"price"`
		} `json:"subscription"`
		Type      *string `json:"type"`
		VNCPort   *uint   `json:"vnc_port,string"`
		VNCScreen *string `json:"vnc_screen"`
	}
	v.Arch = &s.Arch
	v.AvailableKernels = &s.AvailableKernels
	v.BackupTypes = &s.BackupTypes
	v.Core = &s.Core
	v.Disk = &s.Disk
	v.Distro = &s.Distro
	v.GroupID = &s.GroupID
	v.Initrd = &s.Initrd
	v.IPAddrLimit = &s.IPAddrLimit
	v.IPs = &s.IPs
	v.Kernel = &s.Kernel
	v.Label = &s.Label
	v.LocationCode = &s.LocationCode
	v.LocationName = &s.LocationName
	v.Modules = &s.Modules
	v.Name = &s.Name
	v.Notes = &s.Notes
	v.OS = &s.OS
	v.ProductCode = &s.ProductCode
	v.ProductName = &s.ProductName
	v.ProductType = &s.ProductType
	v.RAM = &s.RAM
	v.Root = &s.Root
	v.State = &s.State
	v.Subscription = &s.Subscription
	v.Type = &s.Type
	v.VNCPort = &s.VNCPort
	v.VNCScreen = &s.VNCScreen
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	panic(string(b))
	for i, v := range s.IPs {
		s.IPs[i].MACAddr, err = net.ParseMAC(string(v.MACAddr))
		if err != nil {
			return err
		}
	}
	for s, v := range map[*time.Time]string{&s.Created: v.Created, &s.MaintDate: v.MaintDate, &s.MaintDateEnd: v.MaintDateEnd} {
		if v == "0000-00-00 00:00:00" {
			continue
		}
		*s, err = time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return err
		}
	}
	for s, v := range map[*uint]interface{}{&s.ClientID: v.ClientID, &s.Cores: v.Cores} {
		switch t := v.(type) {
		case string:
			u, err := strconv.ParseUint(t, 10, 0)
			*s = uint(u)
			if err != nil {
				return err
			}
		case int64:
			*s = uint(t)
		}
	}
	for s, v := range map[*bool]uint{&s.Locked: v.Locked, &s.Managed: v.Managed, &s.Mirror: v.Mirror, &s.Rescue: v.Rescue, &s.EmailLogs: v.EmailLogs} {
		*s = v == 1
	}
	var pp []map[string]interface{}
	switch t := v.Partitions.(type) {
	case map[string]interface{}:
		for _, v := range t {
			pp = append(pp, v.(map[string]interface{}))
		}
	case []interface{}:
		for _, v := range t {
			pp = append(pp, v.(map[string]interface{}))
		}
	}
	for _, i := range pp {
		m := map[string]string{}
		for k, v := range i {
			m[k] = v.(string)
		}
		p := struct {
			AlertThreshold float64 `json:"alert_threshold"`
			Backup         bool    `json:"backup"`
			Device         string  `json:"device"`
			DiskTotal      uint    `json:"disk_total"`
			DiskUsed       uint    `json:"disk_used"`
			DRBD           bool    `json:"drbd"`
			FSType         string  `json:"fstype"`
			ID             uint    `json:"id"`
			InodesTotal    uint    `json:"inodes_total"`
			InodesUsed     uint    `json:"inodes_used"`
			Mountpoint     string  `json:"mountpoint"`
			Name           string  `json:"name"`
			NewSize        uint    `json:"new_size"`
			Size           uint    `json:"size"`
			Type           string  `json:"type"`
		}{
			Device:     m["device"],
			FSType:     m["fstype"],
			Mountpoint: m["mountpoint"],
			Name:       m["name"],
			Type:       m["type"],
		}
		p.AlertThreshold, err = strconv.ParseFloat(m["alert_threshold"], 0)
		if err != nil {
			return err
		}
		for k, p := range map[string]*bool{"backup": &p.Backup, "drbd": &p.DRBD} {
			*p = m[k] == "1"
		}
		for k, p := range map[string]*uint{"disk_total": &p.DiskTotal, "disk_used": &p.DiskUsed, "id": &p.ID, "inodes_total": &p.InodesTotal, "inodes_used": &p.InodesUsed, "new_size": &p.NewSize, "size": &p.Size} {
			u, err := strconv.ParseUint(m[k], 10, 0)
			*p = uint(u)
			if err != nil {
				return err
			}
		}
		s.Partitions = append(s.Partitions, p)
	}
	return nil
}

// Retrieves the details for a server.
func (s *ServerService) Get(ctx context.Context, name string) (Server, error) {
	if name == "" {
		return Server{}, errors.New("The server name is missing.")
	}
	var r struct {
		Return  Server `json:"return"`
		Message string `json:"msg"`
		Status  bool   `json:"status"`
	}
	err := s.get(ctx, "/server/get_server.json", url.Values{nameKey: []string{name}}, []string{clientIDKey, nameKey}, &r)
	if err != nil {
		return Server{}, err
	}
	if !r.Status {
		return Server{}, errors.New(r.Message)
	}
	return r.Return, nil
}

// Retrieves the state of a server.
func (s *ServerService) GetState(ctx context.Context, name string) (state string, rescue bool, err error) {
	var r struct {
		Return struct {
			State  string `json:"state"`
			Rescue bool   `json:"rescue"`
		} `json:"return"`
		Message string `json:"msg"`
		Status  bool   `json:"status"`
	}
	err = s.get(ctx, "/server/get_state.json", url.Values{nameKey: []string{name}}, []string{clientIDKey, nameKey}, &r)
	if err != nil {
		return "", false, err
	}
	if !r.Status {
		return "", false, errors.New(r.Message)
	}
	return r.Return.State, r.Return.Rescue, nil
}

/*
// Retrieves statistics for a server.
func (s *ServerService) GetStatistics(headers, queryParams map[string]interface{}) (types.ServerGet_statistics_JsonGetRespBody, error) {
	resp, err := s.client().doReqNoBody("GET", s.client().BaseURI+"/server/get_statistics.json", headers, queryParams)
}

// Lists IP addresses allocated to a client.
func (s *ServerService) ListAllocatedIPAddresses(headers, queryParams map[string]interface{}) (types.ServerList_allocated_ips_JsonGetRespBody, error) {
	resp, err := s.client().doReqNoBody("GET", s.client().BaseURI+"/server/list_allocated_ips.json", headers, queryParams)
}
*/

type ServerListImagesOption interface {
	listImages()
	callOption
}

func (FilterType) listImages() {}
func (FilterOS) listImages()   {}

type Image struct {
	Arch   string `json:"arch"`
	Code   string `json:"code"`
	Distro string `json:"distro"`
	Name   string `json:"name"`
	OS     string `json:"os"`
	Type   string `json:"type"`
}

// Lists images available to provision.
func (s *ServerService) ListImages(ctx context.Context, opts ...ServerListImagesOption) ([]Image, error) {
	data := url.Values{}
	for _, o := range opts {
		o.set(data)
	}
	var r struct {
		Return  []Image `json:"return"`
		Message string  `json:"msg"`
		Status  bool    `json:"status"`
	}
	err := s.get(ctx, "/server/list_images.json", data, []string{filterTypeKey, filterOSKey, filterSortByKey, filterSortDirKey, filterPageSizeKey, filterPageNumberKey}, &r)
	if err != nil {
		return nil, err
	}
	if !r.Status {
		return nil, errors.New(r.Message)
	}
	return r.Return, nil
}

/*
// Lists IP addresses assigned to a client.
func (s *ServerService) ListIPAddresses(headers, queryParams map[string]interface{}) (types.ServerList_ips_JsonGetRespBody, error) {
	resp, err := s.client().doReqNoBody("GET", s.client().BaseURI+"/server/list_ips.json", headers, queryParams)
}
*/

type Location struct {
	AvailableIps uint     `json:"available_ips"`
	Code         string   `json:"code"`
	Datacenter   string   `json:"datacenter"`
	IPv6         bool     `json:"ipv6"`
	Label        string   `json:"label"`
	OS           []string `json:"os"`
	ProductTypes []string `json:"product_types"`
	Public       uint     `json:"public,string"`
}

func (l *Location) UnmarshalJSON(b []byte) error {
	v := struct {
		AvailableIps *uint       `json:"available_ips"`
		Code         *string     `json:"code"`
		Datacenter   *string     `json:"datacenter"`
		IPv6         *bool       `json:"ipv6"`
		Label        *string     `json:"label"`
		OS           *[]string   `json:"os"`
		ProductTypes interface{} `json:"product_types"`
		Public       *uint       `json:"public,string"`
	}{
		&l.AvailableIps,
		&l.Code,
		&l.Datacenter,
		&l.IPv6,
		&l.Label,
		&l.OS,
		nil,
		&l.Public,
	}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch t := v.ProductTypes.(type) {
	case []interface{}:
		for _, v := range t {
			l.ProductTypes = append(l.ProductTypes, v.(string))
		}
	case map[string]interface{}:
		for _, v := range t {
			l.ProductTypes = append(l.ProductTypes, v.(string))
		}
	default:
		return fmt.Errorf("cannot unmarshal product_types: %t", t)
	}
	return nil
}

// Lists locations available to provision to.
func (s *ServerService) ListLocations(ctx context.Context) ([]Location, error) {
	var r struct {
		Return  []Location `json:"return"`
		Message string     `json:"msg"`
		Status  bool       `json:"status"`
	}
	err := s.get(ctx, "/server/list_locations.json", nil, []string{clientIDKey}, &r)
	if err != nil {
		return nil, err
	}
	if !r.Status {
		return nil, errors.New(r.Message)
	}
	return r.Return, nil
}

type ResourceGroup struct {
	ID     uint   `json:"group_id,string"`
	Name   string `json:"group_name"`
	Quotas []struct {
		ID      uint     `json:"attribute_id,string"`
		Name    string   `json:"attribute_name"`
		Unit    string   `json:"attribute_unit"`
		Type    uint     `json:"attribute_type,string"`
		Total   uint     `json:"total_units,string"`
		Used    uint     `json:"used_units,string"`
		Avail   uint     `json:"available_units"`
		Objects []string `json:"objects"`
	} `json:"quotas"`
}

// Lists resources available to a client.
func (s *ServerService) ListResources(ctx context.Context) ([]ResourceGroup, error) {
	var r struct {
		Return  []ResourceGroup `json:"return"`
		Message string          `json:"msg"`
		Status  bool            `json:"status"`
	}
	err := s.get(ctx, "/server/list_resources.json", nil, []string{clientIDKey}, &r)
	if err != nil {
		return nil, err
	}
	if !r.Status {
		return nil, errors.New(r.Message)
	}
	return r.Return, nil
}

type ServerListOption interface {
	serverList()
	callOption
}

// func (direction) serverList()  {}.
func (FilterState) serverList()       {}
func (FilterType) serverList()        {}
func (FilterName) serverList()        {}
func (FilterLocation) serverList()    {}
func (FilterProductType) serverList() {}
func (FilterProductCode) serverList() {}

// Lists all servers for a client.
func (s *ServerService) List(ctx context.Context, opts ...ServerListOption) ([]Server, error) {
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
				TotalPages uint     `json:"total_pages"`
				Data       []Server `json:"data"`
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
		ss = append(ss, r.Return.Data...)
		return r.Return.TotalPages, nil
	})
	return ss, err
}

/*
// Lists statistic types available for a server.
func (s *ServerService) ListStatisticTypes(headers, queryParams map[string]interface{}) (types.ServerList_statistic_types_JsonGetRespBody, error) {
	resp, err := s.client().doReqNoBody("GET", s.client().BaseURI+"/server/list_statistic_types.json", headers, queryParams)
}

// Lists upgrades available for a server.
func (s *ServerService) ListUpgrades(headers, queryParams map[string]interface{}) (types.ServerList_upgrades_JsonGetRespBody, error) {
	resp, err := s.client().doReqNoBody("GET", s.client().BaseURI+"/server/list_upgrades.json", headers, queryParams)
}
*/

type ServerProvisionOption interface {
	serverProvision()
	callOption
}

func (ParamName) serverProvision()      {}
func (ParamIPv4) serverProvision()      {}
func (ParamIPv6) serverProvision()      {}
func (ParamSSHKey) serverProvision()    {}
func (ParamContactID) serverProvision() {}
func (paramBackup) serverProvision()    {}
func (paramSendEmail) serverProvision() {}

// Provisions a new server.
func (s *ServerService) Provision(ctx context.Context, label, location, product, image string, opts ...ServerProvisionOption) (id, job uint, name, password string, addr []net.IP, err error) {
	var r struct {
		Return struct {
			JobID    uint     `json:"json_id,string"`
			Name     string   `json:"name"`
			Password string   `json:"password"`
			IPs      []net.IP `json:"ips"`
			ServerID uint     `json:"server_id,string"`
		} `json:"return"`
		Message string `json:"msg"`
		Status  bool   `json:"status"`
	}
	err = s.postForm(ctx, "/server/provision.json", url.Values{
		labelKey:       []string{label},
		locationKey:    []string{location},
		productCodeKey: []string{product},
		imageKey:       []string{image},
	}, []string{clientIDKey, labelKey, locationKey, productCodeKey, imageKey, paramNameKey, paramIPv4Key, paramIPv6Key, paramSSHKey, paramContactIDKey, paramBackupKey, paramSendEmailKey}, &r)
	if err != nil {
		return 0, 0, "", "", nil, err
	}
	if !r.Status {
		return 0, 0, "", "", nil, errors.New(r.Message)
	}
	return r.Return.ServerID, r.Return.JobID, r.Return.Name, r.Return.Password, r.Return.IPs, nil
}

/*
// Removes an IP address from a server.
func (s *ServerService) RemoveIPAddress(headers, queryParams map[string]interface{}) (types.ServerRemove_ip_JsonPostRespBody, error) {
	resp, err := s.client().doReqWithBody("POST", s.client().BaseURI+"/server/remove_ip.json", nil, headers, queryParams)
}

// Sets the primary IP address for a server.
func (s *ServerService) SetPrimaryIPAddress(headers, queryParams map[string]interface{}) (types.ServerSet_primary_ip_JsonPostRespBody, error) {
	resp, err := s.client().doReqWithBody("POST", s.client().BaseURI+"/server/set_primary_ip.json", nil, headers, queryParams)
}

// Creates a snapshot from a disk
func (s *ServerService) CreateSnapshot(headers, queryParams map[string]interface{}) (types.ServerSnapshot_JsonCreatePostRespBody, error) {
	resp, err := s.client().doReqWithBody("POST", s.client().BaseURI+"/server/snapshot.json/create", nil, headers, queryParams)
}

// Deletes a snapshot
func (s *ServerService) DeleteSnapshot(headers, queryParams map[string]interface{}) (types.ServerSnapshot_JsonDeletePostRespBody, error) {
	resp, err := s.client().doReqWithBody("POST", s.client().BaseURI+"/server/snapshot.json/delete", nil, headers, queryParams)
}

// Lists all active snapshots for a server
func (s *ServerService) ListServerSnapshots(headers, queryParams map[string]interface{}) (types.ServerSnapshot_JsonList_allGetRespBody, error) {
	resp, err := s.client().doReqNoBody("GET", s.client().BaseURI+"/server/snapshot.json/list_all", headers, queryParams)
}

// Restores a server from a snapshot
func (s *ServerService) RestoreSnapshot(headers, queryParams map[string]interface{}) (types.ServerSnapshot_JsonRestorePostRespBody, error) {
	resp, err := s.client().doReqWithBody("POST", s.client().BaseURI+"/server/snapshot.json/restore", nil, headers, queryParams)
}

// Extends the lifetime of a snapshot.
func (s *ServerService) SetLifetime(headers, queryParams map[string]interface{}) (types.ServerSnapshot_JsonSet_lifetimePostRespBody, error) {
	resp, err := s.client().doReqWithBody("POST", s.client().BaseURI+"/server/snapshot.json/set_lifetime", nil, headers, queryParams)
}
*/

type ServerUpdateOption interface {
	serverUpdate()
	callOption
}

func (UpdateLabel) serverUpdate()  {}
func (UpdateNote) serverUpdate()   {}
func (UpdateVNC) serverUpdate()    {}
func (UpdateKernel) serverUpdate() {}

// func (UpdatePartitionThreshold) serverUpdate() {}

// Updates various configuration options of a server.
func (s *ServerService) Update(ctx context.Context, name string, opts ...ServerUpdateOption) error {
	data := url.Values{nameKey: []string{name}}
	for _, o := range opts {
		o.set(data)
	}
	var r voidResponse
	err := s.postForm(ctx, "/server/update.json", data, []string{clientIDKey, nameKey, updateLabelKey, updateNoteKey, updateVNCKey, updateKernelKey /*, updatePartitionThreshold*/}, &r)
	if err != nil {
		return err
	}
	if !r.Status {
		return errors.New(r.Message)
	}
	return nil
}

type ServerUpgradeOption interface {
	serverUpgrade()
	callOption
}

func (UpgradeCores) serverUpgrade() {}
func (UpgradeRAM) serverUpgrade()   {}

// func (UpgradeDisk) serverUpgrade()  {}

// Upgrades specific components of a server.
func (s *ServerService) Upgrade(ctx context.Context, name string, opts ...ServerUpgradeOption) (cores bool, disk map[string]bool, err error) {
	data := url.Values{nameKey: []string{name}}
	for _, o := range opts {
		o.set(data)
	}
	var r struct {
		Return struct {
			Cores bool            `json:""`
			Disk  map[string]bool `json:"disk"`
		} `json:"return"`
		Message string `json:"msg"`
		Status  bool   `json:"status"`
	}
	err = s.postForm(ctx, "/server/upgrade.json", data, []string{clientIDKey, nameKey, upgradeCoresKey, upgradeRAMKey /*, upgradeDiskKey*/}, &r)
	if err != nil {
		return false, nil, err
	}
	if !r.Status {
		return false, nil, errors.New(r.Message)
	}
	return r.Return.Cores, r.Return.Disk, nil
}

// UpgradePlan of the named server.
func (s *ServerService) UpgradePlan(ctx context.Context, name, plan string) error {
	var r voidResponse
	err := s.postForm(ctx, "/server/upgrade_plan.json", url.Values{nameKey: []string{name}, planKey: []string{plan}}, []string{clientIDKey, nameKey, planKey}, &r)
	if err != nil {
		return err
	}
	if !r.Status {
		return errors.New(r.Message)
	}
	return nil
}
