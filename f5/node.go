package f5

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/rabbitt/f5er/mergo"
)

type LBNodeFQDN struct {
	AddressFamily string `json:"addressFamily,omitempty"`
	AutoPopulate  string `json:"autopopulate,omitempty"`
	DownInterval  int    `json:"downInterval,omitempty"`
	// hack - ref issue https://github.com/pr8kerl/f5er/issues/9
	// BIG-IP v12.0 returns a string, whereas v11 returns an int
	// if you use this field, you'll have to convert it explicitly before use :(
	Interval interface{} `json:"interval,omitempty"`
	TmName   string      `json:"tmName,omitempty"`
}

type LBNode struct {
	Name            string      `json:"name,omitempty"`
	Partition       string      `json:"partition,omitempty"`
	FullPath        string      `json:"fullPath,omitempty"`
	Generation      int         `json:"generation,omitempty"`
	Address         string      `json:"address,omitEmpty,omitempty"`
	ConnectionLimit int         `json:"connectionLimit,omitempty"`
	DynamicRatio    int         `json:"dynamicRatio,omitempty"`
	Ephemeral       string      `json:"ephemeral,omitempty"`
	Fqdn            *LBNodeFQDN `json:"fqdn,omitempty"`
	Logging         string      `json:"logging,omitempty"`
	Monitor         string      `json:"monitor,omitempty"`
	RateLimit       string      `json:"rateLimit,omitempty"`
	Ratio           int         `json:"ratio,omitempty"`
	Session         string      `json:"session,omitEmpty,omitempty"`
	State           string      `json:"state,omitEmpty,omitempty"`
}

func (target *LBNode) Merge(source *LBNode, opts ...func(*mergo.Config)) (err error) {
	return mergo.Merge(target, source, opts...)
}

type LBNodeRef struct {
	Link  string   `json:"selfLink,omitempty"`
	Items []LBNode `json:"items,omitempty"`
}

type LBNodes struct {
	Items []LBNode `json:"items,omitempty"`
}

type LBNodeFQDNUpdate struct {
	DownInterval int `json:"downInterval,omitempty"`
	Interval     int `json:"interval,omitempty"`
}

type LBNodeUpdate struct {
	Name            string           `json:"name,omitempty"`
	Partition       string           `json:"partition,omitempty"`
	FullPath        string           `json:"fullPath,omitempty"`
	Generation      int              `json:"generation,omitempty"`
	ConnectionLimit int              `json:"connectionLimit,omitempty"`
	Fqdn            LBNodeFQDNUpdate `json:"fqdn,omitempty"`
	Logging         string           `json:"logging,omitempty"`
	Monitor         string           `json:"monitor,omitempty"`
	RateLimit       string           `json:"rateLimit,omitempty"`
}

type LBNodeStatsDescription struct {
	Description string `json:"description,omitempty"`
}

type LBNodeStatsValue struct {
	Value int `json:"value,omitempty"`
}

type LBNodeStatsInnerEntries struct {
	Addr                     LBNodeStatsDescription `json:"addr,omitempty"`
	CurSessions              LBStatsValue           `json:"curSessions,omitempty"`
	MonitorRule              LBNodeStatsDescription `json:"monitorRule,omitempty"`
	MonitorStatus            LBNodeStatsDescription `json:"monitorStatus,omitempty"`
	TmName                   LBNodeStatsDescription `json:"tmName,omitempty"`
	Serverside_bitsIn        LBStatsValue           `json:"serverside.bitsIn,omitempty"`
	Serverside_bitsOut       LBStatsValue           `json:"serverside.bitsOut,omitempty"`
	Serverside_curConns      LBStatsValue           `json:"serverside.curConns,omitempty"`
	Serverside_maxConns      LBStatsValue           `json:"serverside.maxConns,omitempty"`
	Serverside_pktsIn        LBStatsValue           `json:"serverside.pktsIn,omitempty"`
	Serverside_pktsOut       LBStatsValue           `json:"serverside.pktsOut,omitempty"`
	Serverside_totConns      LBStatsValue           `json:"serverside.totConns,omitempty"`
	SessionStatus            LBNodeStatsDescription `json:"sessionStatus,omitempty"`
	Status_availabilityState LBNodeStatsDescription `json:"status.availabilityState,omitempty"`
	Status_enabledState      LBNodeStatsDescription `json:"status.enabledState,omitempty"`
	Status_statusReason      LBNodeStatsDescription `json:"status.statusReason,omitempty"`
	TotRequests              LBStatsValue           `json:"totRequests,omitempty"`
}

type LBNodeStatsNestedStats struct {
	Kind     string                  `json:"kind,omitempty"`
	SelfLink string                  `json:"selfLink,omitempty"`
	Entries  LBNodeStatsInnerEntries `json:"entries,omitempty"`
}

type LBNodeURLKey struct {
	NestedStats LBNodeStatsNestedStats `json:"nestedStats,omitempty"`
}
type LBNodeStatsOuterEntries map[string]LBNodeURLKey

type LBNodeStats struct {
	Kind       string                  `json:"kind,omitempty"`
	Generation int                     `json:"generation,omitempty"`
	SelfLink   string                  `json:"selfLink,omitempty"`
	Entries    LBNodeStatsOuterEntries `json:"entries,omitempty"`
}

func (f *Device) ShowNodes() (error, *LBNodes) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/node"
	res := LBNodes{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) ShowNode(nname string) (error, *LBNode) {

	//u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/pool/~" + partition + "~" + pname + "?expandSubcollections=true"
	node := strings.Replace(nname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/node/" + node
	res := LBNode{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) ShowNodeStats(nname string) (error, *LBObjectStats) {

	node := strings.Replace(nname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/node/" + node + "/stats"
	res := LBObjectStats{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) ShowAllNodeStats() (error, *LBNodeStats) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/node/stats"
	res := LBNodeStats{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) AddNode(body *json.RawMessage) (error, *LBNode) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/node"
	res := LBNode{}

	// post the request
	err, _ := f.sendRequest(u, POST, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) UpdateNode(nname string, body *json.RawMessage) (error, *LBNode) {

	node := strings.Replace(nname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/node/" + node
	res := LBNode{}

	// put the request
	err, _ := f.sendRequest(u, PUT, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}
}

func (f *Device) PatchNode(name string, patch *LBNode) (error, *LBNode) {
	name = strings.Replace(name, "/", "~", -1)
	url := fmt.Sprintf("%s://%s/mgmt/tm/ltm/node/%s", f.Proto, f.Hostname, name)
	existing := &LBNode{}
	var err error

	// Unless we're overwriting, grab the original and merge the patch with
	// the existing record's data  so that existing settings aren't overwritten,
	// but instead added to.
	if f.MergeStrategy() >= mergo.AppendAdditive {
		err, existing = f.ShowNode(name)
		if err != nil {
			return err, nil
		}

		// merge existing fields into patch so we don't lose settings
		patch.Merge(existing, f.MergeConfig())
	}

	// merge the patch with our existing resource settings so we can see if
	// the patch is already applied or not
	new := &LBNode{}
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

func (f *Device) DeleteNode(nname string) (error, *Response) {

	node := strings.Replace(nname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/node/" + node
	res := json.RawMessage{}

	err, resp := f.sendRequest(u, DELETE, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, resp
	}

}
