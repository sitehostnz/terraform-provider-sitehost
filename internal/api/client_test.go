package api

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

var ctx = context.Background()

func ExampleClient() {
	c := &Client{ID: 970180, Token: os.Getenv("SITEHOST_TOKEN")}
	// i, err := c.API().Info(ctx)
	// fmt.Println(i, err)
	// s, _ := c.Server().Get(ctx, "ch-test2")
	// fmt.Println(s, err)
	// json.NewEncoder(os.Stdout).Encode(s)
	ss, _ := c.Server().List(ctx)
	fmt.Println()
	json.NewEncoder(os.Stdout).Encode(ss)
	// fmt.Println(ss, err)
	// l, err := c.Server().ListLocations(ctx)
	// fmt.Println(len(l), err)
	// r, err := c.Server().ListResources(ctx)
	// fmt.Println(r, err)
	// ii, err := c.Server().ListImages(ctx, FilterTypeDistro)
	// fmt.Println(len(ii), err)
	// j, err := c.Job().Get(ctx, 0, SchedulerJob)
	// fmt.Println(j, err)

	// Output: [Server] <nil>
	// {amd64 [{true pvgrub  pvgrub-1.0-amd64 }] [] 0 25 0 2020-07-31 17:34:42 +0000 UTC 5 ubuntu-xenial 0 1  5 [{4 xenbr0 223.165.66.255 223.165.66.3 45066 223.165.66.82 1 00:0c:29:b9:1c:ea 255.255.255.0 223.165.66.0 377 32 true rdns.223.165.66.82.sth.nz 41176} {6 xenbr0 2403:7000:8000:bff:ffff:ffff:ffff:ffff 2403:7000:8000:b00::3 45067 2403:7000:8000:b00::7d 1 00:0c:29:b9:1c:ea ffff:ffff:ffff:ff00:: 2403:7000:8000:b00:: 380 128 true 2403-7000-8000-b00-0-7d.v6.sitehost.co.nz 41176}] pvgrub-1.0-amd64 test AKLCITY NZ - Auckland Central false 0001-01-01 00:00:00 +0000 UTC 0001-01-01 00:00:00 +0000 UTC false false  ch-test2  linux [{1.0 0 ch-test2-data 4896 336 0 ext4 45017 655360 586 /data xvda3 0 5 SSD}] CLDCON1 Cloud Container - 1 Core CLDCON 1024 false xvda2 On {CLDCON1 Cloud Container - 1 Core 35} VPS 6022 800x600} <nil>
	// [{amd64 [] [] 970180  0 2020-07-31 17:34:42 +0000 UTC 5 ubuntu-xenial     []  test  NZ - Auckland Central false 0001-01-01 00:00:00 +0000 UTC 0001-01-01 00:00:00 +0000 UTC false false  ch-test2  linux [] CLDCON1 Cloud Container - 1 Core CLDCON 1 false  On {  } VPS  }] <nil>
	// 10 <nil>
	// [{1 NZ - Linux Servers [{1 VPS Disk Space GB 0 17 17 0 [ch-test2]} {2 VPS Memory MB 0 1024 1024 0 [ch-test2]}]}] <nil>
	// 17 <nil>
	// {0001-01-01 00:00:00 +0000 UTC 0001-01-01 00:00:00 +0000 UTC 0001-01-01 00:00:00 +0000 UTC   []} Unauthorised. You do not have access to this Api method (101019) <nil>
}
