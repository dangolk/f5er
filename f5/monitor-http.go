package f5

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/rabbitt/f5er/mergo"
)

type LBMonitorHttp struct {
	Name                     string `json:"name"`
	Partition                string `json:"partition"`
	FullPath                 string `json:"fullPath"`
	Adaptive                 string `json:"adaptive"`
	AdaptiveDivergenceType   string `json:"adaptiveDivergenceType"`
	AdaptiveDivergenceValue  int    `json:"adaptiveDivergenceValue"`
	AdaptiveLimit            int    `json:"adaptiveLimit"`
	AdaptiveSamplingTimespan int    `json:"adaptiveSamplingTimespan"`
	DefaultsFrom             string `json:"defaultsFrom"`
	Destination              string `json:"destination"`
	Interval                 int    `json:"interval"`
	IpDscp                   int    `json:"ipDscp"`
	ManualResume             string `json:"manualResume"`
	Recv                     string `json:"recv"`
	Reverse                  string `json:"reverse"`
	Send                     string `json:"send"`
	TimeUntilUp              int    `json:"timeUntilUp"`
	Timeout                  int    `json:"timeout"`
	Transparent              string `json:"transparent"`
	UpInterval               int    `json:"upInterval"`
}

func (target *LBMonitorHttp) Merge(source *LBMonitorHttp, opts ...func(*mergo.Config)) (err error) {
	return mergo.Merge(target, source, opts...)
}

type LBMonitorHttpRef struct {
	Items []LBMonitorHttp `json:"items"`
}

func (f *Device) ShowMonitorsHttp() (error, *LBMonitorHttpRef) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/monitor/http"
	res := LBMonitorHttpRef{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) ShowMonitorHttp(vname string) (error, *LBMonitorHttp) {

	vname = strings.Replace(vname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/monitor/http/" + vname + "?expandSubcollections=true"
	res := LBMonitorHttp{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) AddMonitorHttp(body *json.RawMessage) (error, *LBMonitorHttp) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/monitor/http"
	res := LBMonitorHttp{}

	// post the request
	err, _ := f.sendRequest(u, POST, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}
}

func (f *Device) UpdateMonitorHttp(vname string, body *json.RawMessage) (error, *LBMonitorHttp) {

	vname = strings.Replace(vname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/monitor/http/" + vname
	res := LBMonitorHttp{}

	// put the request
	err, _ := f.sendRequest(u, PUT, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) PatchMonitorHttp(name string, patch *LBMonitorHttp) (error, *LBMonitorHttp) {
	name = strings.Replace(name, "/", "~", -1)
	url := fmt.Sprintf("%s://%s/mgmt/tm/ltm/monitor/http/%s", f.Proto, f.Hostname, name)
	existing := &LBMonitorHttp{}
	var err error

	// Unless we're overwriting, grab the original and merge the patch with
	// the existing record's data  so that existing settings aren't overwritten,
	// but instead added to.
	if f.MergeStrategy() >= mergo.AppendAdditive {
		err, existing = f.ShowMonitorHttp(name)
		if err != nil {
			return err, nil
		}

		// merge existing fields into patch so we don't lose settings
		patch.Merge(existing, f.MergeConfig())
	}

	// merge the patch with our existing resource settings so we can see if
	// the patch is already applied or not
	new := &LBMonitorHttp{}
	new.Merge(patch, f.MergeConfig(), func(c *mergo.Config) { c.SkipEmptyFields = false })
	new.Merge(existing, f.MergeConfig(), func(c *mergo.Config) { c.SkipEmptyFields = false })

	if f.DryRun() {
		fmt.Printf("Patching: %s\nPatch Diff:\n%s\nPatch Data (merge strategy: %s):\n",
			url, cmp.Diff(existing, new, cmpopts.EquateEmpty()), f.MergeStrategy())
		return nil, patch
	} else {
		if cmp.Equal(new, existing, cmpopts.EquateEmpty()) {
			return nil, existing
		} else {
			err, _ = f.sendRequest(url, PATCH, patch, existing)
			if err != nil {
				return err, nil
			} else {
				return nil, existing
			}
		}
	}
}
func (f *Device) DeleteMonitorHttp(vname string) (error, *Response) {

	vname = strings.Replace(vname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/monitor/http/" + vname
	res := json.RawMessage{}

	err, resp := f.sendRequest(u, DELETE, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, resp
	}

}
