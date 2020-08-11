package api

import (
	"net"
	"reflect"
	"strconv"
	"testing"
)

// nolint:golint,stylecheck
var ServerClient_ProvisionTests = []struct {
	body         string
	id, job      uint
	name, passwd string
	addr         []net.IP
	err          bool
}{
	{body: `{"msg":"Unauthorised. The Api key is missing."}`, err: true},
	{body: `{"return":{"job_id":"2568588","name":"example","password":"cvSOBrC0Bzth","ips":["223.165.66.169","2403:7000:8000:b00::4b"],"server_id":"41186"},"msg":"Successful","status":true}`, id: 41186, job: 2568588, name: "example", passwd: "cvSOBrC0Bzth", addr: []net.IP{{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 223, 165, 66, 169}, {36, 3, 112, 0, 128, 0, 11, 0, 0, 0, 0, 0, 0, 0, 0, 75}}},
}

func TestServerClient_Provision(t *testing.T) {
	for i, tt := range ServerClient_ProvisionTests {
		tt := tt
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			id, job, name, passwd, addr, err := (*ServerClient)(client(tt.body)).Provision(ctx, "", "", "", "")
			switch {
			case tt.err && err == nil:
				t.Errorf("err should not be: <nil>")
			case !tt.err && err != nil:
				t.Errorf("err should be <nil>: %v", err)
			case tt.id != id:
				t.Errorf("want: %v, got: %v", tt.id, id)
			case tt.job != job:
				t.Errorf("want: %v, got: %v", tt.job, job)
			case tt.name != name:
				t.Errorf("want: %v, got: %v", tt.name, name)
			case tt.passwd != passwd:
				t.Errorf("want: %v, got: %v", tt.passwd, passwd)
			case !reflect.DeepEqual(tt.addr, addr):
				t.Errorf("want: %v, got: %v", tt.addr, addr)
			}
		})
	}
}

// nolint:golint,stylecheck
var ServerClient_ListTests = []struct {
	body []string
	want []Server
	err  bool
}{
	{body: []string{`{"msg":"Unauthorised. The Api key is missing."}`}, err: true},
	{body: []string{`{"return":{"total_items":0,"current_items":0,"current_page":1,"total_pages":1,"data":[]},"msg":"Successful","status":true}`}},
	{body: []string{
		`{"return":{"total_items":10,"current_items":0,"current_page":1,"total_pages":5,"data":[]},"msg":"Successful","status":true}`,
		`{"return":{"total_items":10,"current_items":1,"current_page":2,"total_pages":5,"data":[{}]},"msg":"Successful","status":true}`,
		`{"return":{"total_items":10,"current_items":2,"current_page":3,"total_pages":5,"data":[{},{}]},"msg":"Successful","status":true}`,
		`{"return":{"total_items":10,"current_items":3,"current_page":4,"total_pages":5,"data":[{},{},{}]},"msg":"Successful","status":true}`,
		`{"return":{"total_items":10,"current_items":4,"current_page":5,"total_pages":5,"data":[{},{},{},{}]},"msg":"Successful","status":true}`,
	}, want: make([]Server, 10)},
	// {body: `{"return":{"total_items":1,"current_items":1,"current_page":1,"total_pages":1,"data":[{"location":"AKLCITY","location_name":"NZ - Auckland Central","product_code":"XENLIT","product_name":"VPS Hosting - Linux 1.5GB","product_type":"LINVPS","client_id":"970180","created":"2020-08-06 19:09:34","name":"example","label":"example","type":"VPS","state":"On","locked":"0","rescue":"0","cores":1,"disk":"15","ram":"1.5","distro":"ubuntu-xenial","arch":"amd64","os":"linux","managed":"0","maint_date":"0000-00-00 00:00:00","maint_date_end":"0000-00-00 00:00:00","pending":"0","server_id":"41210","disk_new":"0","mirror":"0","primary_ips":[{"prefix":"32","ip_addr":"120.138.18.109"},{"prefix":"128","ip_addr":"2403:7000:8000:300::3e"}],"backups_enabled":false,"backups_product":"SRVC-OFFSITE"}]},"msg":"Successful","status":true}`, want: []Server{{ [{amd64 [] false SRVC-OFFSITE [] 970180 0 1 2020-08-06 19:09:34 +0000 UTC 15 0 ubuntu-xenial false 0 41210  [] 0 [{0  <nil> <nil> 0 120.138.18.109 0  <nil> <nil> 0 32 true  0} {0  <nil> <nil> 0 2403:7000:8000:300::3e 0  <nil> <nil> 0 128 true  0}]  example {0 Completed } AKLCITY NZ - Auckland Central  false 0001-01-01 00:00:00 +0000 UTC 0001-01-01 00:00:00 +0000 UTC false false  example  linux [] XENLIT VPS Hosting - Linux 1.5GB LINVPS 1536 false  On {  0} VPS 0 {0 0}}]}
}

func TestServerClient_List(t *testing.T) {
	for i, tt := range ServerClient_ListTests {
		tt := tt
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := (*ServerClient)(client(tt.body...)).List(ctx)
			switch {
			case tt.err && err == nil:
				t.Errorf("err should not be: <nil>")
			case !tt.err && err != nil:
				t.Errorf("err should be <nil>: %v", err)
			case !reflect.DeepEqual(tt.want, got):
				t.Errorf("want: %v, got: %v", tt.want, got)
			}
		})
	}
}

// nolint:golint,stylecheck
var ServerClient_GetTests = []struct {
	name string
	body string
	want Server
	err  bool
}{
	{err: true},
	{name: "example", body: `{"msg":"Unauthorised. The Api key is missing."}`, err: true},
	// {name: "example", body: `{"return":{"name":"example","label":"example","client_id":970180,"created":"2020-08-05 21:04:45","type":"VPS","ram":"1536","root":"xvda2","disk":15,"cores":"0","core":"36","arch":"amd64","kernel":"pvgrub-1.0-amd64","initrd":"","modules":"","distro":"ubuntu-xenial","os":"linux","rescue":"0","managed":"0","locked":"0","state":"On","maint_date":"0000-00-00 00:00:00","maint_date_end":"0000-00-00 00:00:00","email_logs":"0","vnc_port":"6025","vnc_screen":"800x600","ip_addr_limit":"5","notes":"","ips":[{"id":"45107","server_id":"41197","bridge":"xenbr0","mac_addr":"00:0c:29:cc:47:57","ip_type":"1","prefix":"32","primary":true,"addr_family":4,"network_id":"6","rdns":"rdns.120.138.18.109.sth.nz","ip_addr":"120.138.18.109","network":"120.138.18.0","gateway":"120.138.18.3","netmask":"255.255.255.0","broadcast":"120.138.18.255"},{"id":"45108","server_id":"41197","bridge":"xenbr0","mac_addr":"00:0c:29:cc:47:57","ip_type":"1","prefix":"128","primary":true,"addr_family":6,"network_id":"298","rdns":"2403-7000-8000-300-0-3e.v6.sitehost.co.nz","ip_addr":"2403:7000:8000:300::3e","network":"2403:7000:8000:300::","gateway":"2403:7000:8000:300::3","netmask":"ffff:ffff:ffff:ff00::","broadcast":"2403:7000:8000:3ff:ffff:ffff:ffff:ffff"}],"interfaces":["xenbr0"],"group_id":"1","partitions":[{"id":"45062","name":"xvda1","device":"example-swap","mountpoint":"swap","size":"1","new_size":"0","fstype":"swap","drbd":"0","backup":"0","disk_total":"0","disk_used":"0","inodes_total":"0","inodes_used":"0","alert_threshold":"1.0","type":"SSD"},{"id":"45063","name":"xvda2","device":"example-disk","mountpoint":"\/","size":"15","new_size":"0","fstype":"ext4","drbd":"0","backup":"0","disk_total":"0","disk_used":"0","inodes_total":"0","inodes_used":"0","alert_threshold":"1.0","type":"SSD"}],"backup_types":[],"available_kernels":[{"default":true,"kernel":"pvgrub-1.0-amd64","initrd":"","modules":"","hypervisor":"pvgrub"}],"product_code":"XENLIT","product_type":"LINVPS","product_name":"VPS Hosting - Linux 1.5GB","subscription":{"code":"XENLIT","name":"VPS Hosting - Linux 1.5GB","price":"30"},"location_name":"NZ - Auckland Central","location_code":"AKLCITY","mirror":"0","last_job":{"id":"2568631","type":"daemon","state":"Complete"},"backups_enabled":false,"backups_product":"SRVC-OFFSITE","location_node":"NZ-WAL1-6G2QEQ"},"msg":"Successful","status":true}`, want: Server{{amd64 [{true pvgrub  pvgrub-1.0-amd64 }] false SRVC-OFFSITE [] 970180 36 0 2020-08-05 21:04:45 +0000 UTC 15 0 ubuntu-xenial false 1 0  [xenbr0] 5 [{4 xenbr0 120.138.18.255 120.138.18.3 45107 120.138.18.109 1 00:0c:29:cc:47:57 255.255.255.0 120.138.18.0 6 32 true rdns.120.138.18.109.sth.nz 41197} {6 xenbr0 2403:7000:8000:3ff:ffff:ffff:ffff:ffff 2403:7000:8000:300::3 45108 2403:7000:8000:300::3e 1 00:0c:29:cc:47:57 ffff:ffff:ffff:ff00:: 2403:7000:8000:300:: 298 128 true 2403-7000-8000-300-0-3e.v6.sitehost.co.nz 41197}] pvgrub-1.0-amd64 example {2568631 Complete daemon} AKLCITY NZ - Auckland Central NZ-WAL1-6G2QEQ false 0001-01-01 00:00:00 +0000 UTC 0001-01-01 00:00:00 +0000 UTC false false  example  linux [{1 false example-swap 0 0 false swap 45062 0 0 swap xvda1 0 1 SSD} {1 false example-disk 0 0 false ext4 45063 0 0 / xvda2 0 15 SSD}] XENLIT VPS Hosting - Linux 1.5GB LINVPS 1536 false xvda2 On {XENLIT VPS Hosting - Linux 1.5GB 30} VPS 6025 {800 600}}}},
}

func TestServerClient_Get(t *testing.T) {
	for i, tt := range ServerClient_GetTests {
		tt := tt
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := (*ServerClient)(client(tt.body)).Get(ctx, tt.name)
			switch {
			case tt.err && err == nil:
				t.Errorf("err should not be: <nil>")
			case !tt.err && err != nil:
				t.Errorf("err should be <nil>: %v", err)
			case !reflect.DeepEqual(tt.want, got):
				t.Errorf("want: %v, got: %v", tt.want, got)
			}
		})
	}
}

// nolint:golint,stylecheck
var ServerClient_DeleteTests = []struct {
	body string
	want uint
	err  bool
}{
	{body: `{"msg":"Unauthorised. The Api key is missing."}`, err: true},
	{body: `{"return":{"job_id":"2568630"},"msg":"Successful","status":true}`, want: 2568630},
}

func TestServerClient_Delete(t *testing.T) {
	for i, tt := range ServerClient_DeleteTests {
		tt := tt
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := (*ServerClient)(client(tt.body)).Delete(ctx, "")
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
