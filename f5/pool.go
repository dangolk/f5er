package f5

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/rabbitt/f5er/mergo"
)

type LBPoolMemberFQDN struct {
	Autopopulate string `json:"autopopulate,omitempty"`
}

// a pool member
type LBPoolMember struct {
	Name            string            `json:"name,omitempty"`
	Partition       string            `json:"partition,omitempty"`
	FullPath        string            `json:"fullPath,omitempty"`
	Address         string            `json:"address,omitempty"`
	ConnectionLimit int               `json:"connectionLimit,omitempty"`
	DynamicRatio    int               `json:"dynamicRatio,omitempty"`
	Ephemeral       string            `json:"ephemeral,omitempty"`
	Fqdn            *LBPoolMemberFQDN `json:"fqdn,omitempty"`
	InheritProfile  string            `json:"inheritProfile,omitempty"`
	Logging         string            `json:"logging,omitempty"`
	Monitor         string            `json:"monitor,omitempty"`
	PriorityGroup   int               `json:"priorityGroup,omitempty"`
	RateLimit       string            `json:"rateLimit,omitempty"`
	Ratio           int               `json:"ratio,omitempty"`
	Session         string            `json:"session,omitempty"`
	State           string            `json:"state,omitempty"`
}

// a pool member reference - just a link and an array of pool members
type LBPoolMemberRef struct {
	Link  string         `json:"link,omitempty"`
	Items []LBPoolMember `json:"items,omitempty"`
}

type LBPoolMembers struct {
	Link  string         `json:"selfLink,omitempty"`
	Items []LBPoolMember `json:"items,omitempty"`
}

// used by online/offline
type LBPoolMemberState struct {
	State   string `json:"state,omitempty"`
	Session string `json:"session,omitempty"`
}

type LBPool struct {
	Name                   string         `json:"name,omitempty"`
	Partition              string         `json:"partition,omitempty"`
	FullPath               string         `json:"fullPath,omitempty"`
	Generation             int            `json:"generation,omitempty"`
	AllowNat               string         `json:"allowNat,omitempty"`
	AllowSnat              string         `json:"allowSnat,omitempty"`
	IgnorePersistedWeight  string         `json:"ignorePersistedWeight,omitempty"`
	IpTosToClient          string         `json:"ipTosToClient,omitempty"`
	IpTosToServer          string         `json:"ipTosToServer,omitempty"`
	LinkQosToClient        string         `json:"linkQosToClient,omitempty"`
	LinkQosToServer        string         `json:"linkQosToServer,omitempty"`
	LoadBalancingMode      string         `json:"loadBalancingMode,omitempty"`
	MinActiveMembers       int            `json:"minActiveMembers,omitempty"`
	MinUpMembers           int            `json:"minUpMembers,omitempty"`
	MinUpMembersAction     string         `json:"minUpMembersAction,omitempty"`
	MinUpMembersChecking   string         `json:"minUpMembersChecking,omitempty"`
	Monitor                string         `json:"monitor,omitempty"`
	QueueDepthLimit        int            `json:"queueDepthLimit,omitempty"`
	QueueOnConnectionLimit string         `json:"queueOnConnectionLimit,omitempty"`
	QueueTimeLimit         int            `json:"queueTimeLimit,omitempty"`
	ReselectTries          int            `json:"reselectTries,omitempty"`
	ServiceDownAction      string         `json:"serviceDownAction,omitempty"`
	SlowRampTime           int            `json:"slowRampTime,omitempty"`
	Members                []LBPoolMember `json:"members,omitempty"`
}

func (target *LBPool) Merge(source *LBPool, opts ...func(*mergo.Config)) (err error) {
	return mergo.Merge(target, source, opts...)
}

func (pool *LBPool) UnmarshalJSON(data []byte) error {
	// Strip out the Policies and Profiles Reference entries, converting them
	// so simple policies and profiles arrays.
	type Alias LBPool
	aux := &struct {
		MemberReference LBPoolMemberRef `json:"membersReference"`
		*Alias
	}{
		Alias: (*Alias)(pool),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if len(aux.Members) > 0 && len(aux.MemberReference.Items) > 0 {
		return fmt.Errorf("error: members and memberReference can not both be set for: %s", aux.Name)
	} else {
		if len(aux.MemberReference.Items) > 0 {
			aux.Members = aux.MemberReference.Items
		}
	}

	pool = (*LBPool)(aux.Alias)
	return nil
}

type LBPools struct {
	Items []LBPool `json:"items"`
}

type LBPoolStatsDescription struct {
	Description string `json:"description"`
}

type LBPoolStatsValue struct {
	Value int `json:"value"`
}

type LBPoolStatsInnerEntries struct {
	ActiveMemberCnt          LBStatsValue           `json:"activeMemberCnt"`
	ConnqAll_ageEdm          LBStatsValue           `json:"connqAll.ageEdm"`
	ConnqAll_ageEma          LBStatsValue           `json:"connqAll.ageEma"`
	ConnqAll_ageHead         LBStatsValue           `json:"connqAll.ageHead"`
	ConnqAll_ageMax          LBStatsValue           `json:"connqAll.ageMax"`
	ConnqAll_depth           LBStatsValue           `json:"connqAll.depth"`
	ConnqAll_serviced        LBStatsValue           `json:"connqAll.serviced"`
	Connq_ageEdm             LBStatsValue           `json:"connq.ageEdm"`
	Connq_ageEma             LBStatsValue           `json:"connq.ageEma"`
	Connq_ageHead            LBStatsValue           `json:"connq.ageHead"`
	Connq_ageMax             LBStatsValue           `json:"connq.ageMax"`
	Connq_depth              LBStatsValue           `json:"connq.depth"`
	Connq_serviced           LBStatsValue           `json:"connq.serviced"`
	CurSessions              LBStatsValue           `json:"curSessions"`
	MinActiveMembers         LBStatsValue           `json:"minActiveMembers"`
	MonitorRule              LBPoolStatsDescription `json:"monitorRule"`
	TmName                   LBPoolStatsDescription `json:"tmName"`
	Serverside_bitsIn        LBStatsValue           `json:"serverside.bitsIn"`
	Serverside_bitsOut       LBStatsValue           `json:"serverside.bitsOut"`
	Serverside_curConns      LBStatsValue           `json:"serverside.curConns"`
	Serverside_maxConns      LBStatsValue           `json:"serverside.maxConns"`
	Serverside_pktsIn        LBStatsValue           `json:"serverside.pktsIn"`
	Serverside_pktsOut       LBStatsValue           `json:"serverside.pktsOut"`
	Serverside_totConns      LBStatsValue           `json:"serverside.totConns"`
	Status_availabilityState LBPoolStatsDescription `json:"status.availabilityState"`
	Status_enabledState      LBPoolStatsDescription `json:"status.enabledState"`
	Status_statusReason      LBPoolStatsDescription `json:"status.statusReason"`
	TotRequests              LBStatsValue           `json:"totRequests"`
}

type LBPoolStatsNestedStats struct {
	Kind     string                  `json:"kind"`
	SelfLink string                  `json:"selfLink"`
	Entries  LBPoolStatsInnerEntries `json:"entries"`
}

type LBPoolURLKey struct {
	NestedStats LBPoolStatsNestedStats `json:"nestedStats"`
}
type LBPoolStatsOuterEntries map[string]LBPoolURLKey

type LBPoolStats struct {
	Kind       string                  `json:"kind"`
	Generation int                     `json:"generation"`
	SelfLink   string                  `json:"selfLink"`
	Entries    LBPoolStatsOuterEntries `json:"entries"`
}

func (f *Device) ShowPools() (error, *LBPools) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/pool"
	res := LBPools{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) ShowPool(pname string) (error, *LBPool) {

	pool := strings.Replace(pname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/pool/" + pool + "?expandSubcollections=true"
	res := LBPool{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) ShowPoolStats(pname string) (error, *LBObjectStats) {

	pool := strings.Replace(pname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/pool/" + pool + "/stats"
	res := LBObjectStats{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}
}

func (f *Device) ShowAllPoolStats() (error, *LBPoolStats) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/pool/stats"
	res := LBPoolStats{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}
}

func (f *Device) AddPool(body *json.RawMessage) (error, *LBPool) {
	// we use json.RawMessage so we can modify the input file without using a struct
	// use of a struct will send all available fields, some of which can't be modified

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/pool"
	res := LBPool{}

	// post the request
	err, _ := f.sendRequest(u, POST, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) UpdatePool(pname string, body *json.RawMessage) (error, *LBPool) {

	pool := strings.Replace(pname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/pool/" + pool
	res := LBPool{}

	// put the request
	err, _ := f.sendRequest(u, PUT, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}
}

func (f *Device) PatchPool(name string, patch *LBPool) (error, *LBPool) {
	name = strings.Replace(name, "/", "~", -1)
	url := fmt.Sprintf("%s://%s/mgmt/tm/ltm/pool/%s", f.Proto, f.Hostname, name)
	existing := &LBPool{}
	var err error

	// Unless we're overwriting, grab the original and merge the patch with
	// the existing record's data  so that existing settings aren't overwritten,
	// but instead added to.
	if f.MergeStrategy() >= mergo.AppendAdditive {
		err, existing = f.ShowPool(name)
		if err != nil {
			return err, nil
		}

		// merge existing fields into patch so we don't lose settings
		patch.Merge(existing, f.MergeConfig())
	}

	// merge the patch with our existing resource settings so we can see if
	// the patch is already applied or not
	new := &LBPool{}
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

func (f *Device) DeletePool(pname string) (error, *Response) {

	pool := strings.Replace(pname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/pool/" + pool
	res := json.RawMessage{}

	err, resp := f.sendRequest(u, DELETE, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, resp
	}

}

func (f *Device) ShowPoolMembers(pname string) (error, *LBPoolMembers) {

	pool := strings.Replace(pname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/pool/" + pool + "/members"
	res := LBPoolMembers{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) ShowPoolMembersStats(pname string) (error, *LBPoolStats) {

	pool := strings.Replace(pname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/pool/" + pool + "/members/stats"
	res := LBPoolStats{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}
}

func (f *Device) ShowAllPoolMembersStats() (error, *LBPoolStats) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/pool/members/stats"
	res := LBPoolStats{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}
}

func (f *Device) AddPoolMembers(pname string, body *json.RawMessage) (error, *LBPoolMembers) {

	pool := strings.Replace(pname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/pool/" + pool + "/members"
	res := LBPoolMembers{}

	// post the request
	err, _ := f.sendRequest(u, POST, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) UpdatePoolMembers(pname string, body *json.RawMessage) (error, *LBPoolMembers) {

	pool := strings.Replace(pname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/pool/" + pool + "/members"
	res := LBPoolMembers{}

	// put the request
	err, _ := f.sendRequest(u, PUT, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) DeletePoolMembers(pname string) (error, *Response) {

	pool := strings.Replace(pname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/pool/" + pool + "/members"
	res := json.RawMessage{}

	err, resp := f.sendRequest(u, DELETE, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, resp
	}

}

func (f *Device) OnlinePoolMember(pname string, mname string) (error, *Response) {

	pmember := strings.Replace(mname, "/", "~", -1)
	pool := strings.Replace(pname, "/", "~", -1)

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/pool/" + pool + "/members/" + pmember
	res := json.RawMessage{}

	/*
	   {"state": "user-down", "session": "user-disabled"} (Member Forced Offline in GUI)
	   {"state": "user-up", "session": "user-disabled"} (Member Disabled in GUI)
	   {"state": "user-up", "session": "user-enabled"}  (Member Enabled in GUI)
	*/
	body := LBPoolMemberState{"user-up", "user-enabled"}

	// put the request
	err, resp := f.sendRequest(u, PUT, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, resp
	}

}

func (f *Device) OfflinePoolMember(pname string, mname string) (error, *Response) {

	pmember := strings.Replace(mname, "/", "~", -1)
	pool := strings.Replace(pname, "/", "~", -1)

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/pool/" + pool + "/members/" + pmember
	res := json.RawMessage{}

	/*
	   {"state": "user-down", "session": "user-disabled"} (Member Forced Offline in GUI)
	   {"state": "user-up", "session": "user-disabled"} (Member Disabled in GUI)
	   {"state": "user-up", "session": "user-enabled"}  (Member Enabled in GUI)
	*/
	body := LBPoolMemberState{"user-up", "user-disabled"}

	// put the request
	err, resp := f.sendRequest(u, PUT, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, resp
	}

}
func (f *Device) OfflinePoolMemberForced(pname string, mname string) (error, *Response) {

	pmember := strings.Replace(mname, "/", "~", -1)
	pool := strings.Replace(pname, "/", "~", -1)

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/pool/" + pool + "/members/" + pmember
	res := json.RawMessage{}

	/*
	   {"state": "user-down", "session": "user-disabled"} (Member Forced Offline in GUI)
	   {"state": "user-up", "session": "user-disabled"} (Member Disabled in GUI)
	   {"state": "user-up", "session": "user-enabled"}  (Member Enabled in GUI)
	*/
	body := LBPoolMemberState{"user-down", "user-disabled"}

	// put the request
	err, resp := f.sendRequest(u, PUT, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, resp
	}

}
