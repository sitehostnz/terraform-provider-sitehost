package api

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"
)

// nolint:golint,stylecheck
var JobClient_GetTests = []struct {
	body string
	want Job
	err  bool
}{
	{body: `{"msg":"Error: The specified job does not exist or you do not have access to it.","status":false}`, err: true},
	{body: `{"return":{"created":"2020-08-05 19:23:28","started":"2020-08-05 19:23:28","completed":"2020-08-05 19:23:54","message":"VPS(): VPS example1 Provisioned Successfully","state":"Completed","logs":[{"date":"2020-08-05 19:23:29","level":"3","message":"[example1] Preparing to provision"},{"date":"2020-08-05 19:23:29","level":"3","message":"[example1] Creating disks"},{"date":"2020-08-05 19:23:33","level":"3","message":"[example1] Installing operating system"},{"date":"2020-08-05 19:23:48","level":"3","message":"[example1] Configuring system"},{"date":"2020-08-05 19:23:51","level":"3","message":"[example1] Finished provisioning"},{"date":"2020-08-05 19:23:51","level":"3","message":"[example1] Updating server configuration"},{"date":"2020-08-05 19:23:52","level":"3","message":"[example1] Booting up"},{"date":"2020-08-05 19:23:54","level":"3","message":"[example1] Server is now on"},{"date":"2020-08-05 19:23:54","level":"3","message":"[example1] Finished provisioning server"}]},"msg":"Successful.","status":true}`, want: Job{utc(1596655408), utc(1596655408), utc(1596655434), "VPS(): VPS example1 Provisioned Successfully", "Completed", []struct {
		Date    time.Time
		Level   uint
		Message string
	}{{utc(1596655409), 3, "[example1] Preparing to provision"}, {utc(1596655409), 3, "[example1] Creating disks"}, {utc(1596655413), 3, "[example1] Installing operating system"}, {utc(1596655428), 3, "[example1] Configuring system"}, {utc(1596655431), 3, "[example1] Finished provisioning"}, {utc(1596655431), 3, "[example1] Updating server configuration"}, {utc(1596655432), 3, "[example1] Booting up"}, {utc(1596655434), 3, "[example1] Server is now on"}, {utc(1596655434), 3, "[example1] Finished provisioning server"}}}},
}

func TestJobClient_Get(t *testing.T) {
	for i, tt := range JobClient_GetTests {
		tt := tt
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := (*JobClient)(client(tt.body)).Get(ctx, 0, JobType(""))
			switch {
			case tt.err && err == nil:
				t.Errorf("err should not be: <nil>")
			case !tt.err && err != nil:
				t.Errorf("err should be <nil>: %v", err)
			case !reflect.DeepEqual(tt.want, got):
				fmt.Printf("DEBUG>%#v %#v\n",
					tt.want.Created,
					got.Created,
				)
				fmt.Println("DEBUG>",
					tt.want.Created == got.Created,
					tt.want.Started == got.Started,
					tt.want.Completed == got.Completed,
					tt.want.Message == got.Message,
					tt.want.State == got.State,
				)
				t.Errorf("want: %v, got: %v", tt.want, got)
			}
		})
	}
}
