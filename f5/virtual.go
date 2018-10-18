package f5

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/rabbitt/f5er/mergo"
)

type LBVirtualPolicy struct {
	Name      string `json:"name,omitempty"`
	Partition string `json:"partition,omitempty"`
	FullPath  string `json:"fullPath,omitempty"`
}

type LBVirtualProfile struct {
	Name      string `json:"name,omitempty"`
	Partition string `json:"partition,omitempty"`
	FullPath  string `json:"fullPath,omitempty"`
	Context   string `json:"context,omitempty"`
}

type LBVirtualPersist struct {
	Name      string `json:"name,omitempty"`
	Partition string `json:"partition,omitempty"`
	TmDefault string `json:"tmDefault,omitempty"`
}

type LBVirtualSNAT struct {
	Type string `json:"type,omitempty"`
}

type LBVirtual struct {
	AddressStatus              string             `json:"addressStatus,omitempty"`
	AutoLastHop                string             `json:"autoLasthop,omitempty"`
	CmpEnabled                 string             `json:"cmpEnabled,omitempty"`
	ConnectionLimit            int                `json:"connectionLimit,omitempty"`
	Destination                string             `json:"destination,omitempty"`
	Enabled                    bool               `json:"enabled,omitempty"`
	FullPath                   string             `json:"fullPath,omitempty"`
	IpProtocol                 string             `json:"ipProtocol,omitempty"`
	IpIntelligencePolicy       string             `json:"ipIntelligencePolicy,omitempty"`
	Mask                       string             `json:"mask,omitempty"`
	Mirror                     string             `json:"mirror,omitempty"`
	MobileApptunnel            string             `json:"mobileAppTunnel,omitempty"`
	Name                       string             `json:"name,omitempty"`
	Nat64                      string             `json:"nat64,omitempty"`
	Partition                  string             `json:"partition,omitempty"`
	Persist                    []LBVirtualPersist `json:"persist,omitempty"`
	Policies                   []LBVirtualPolicy  `json:"policies,omitempty"`
	Pool                       string             `json:"pool,omitempty"`
	Profiles                   []LBVirtualProfile `json:"profiles,omitempty"`
	RateLimitDstMask           int                `json:"rateLimitDstMask,omitempty"`
	RateLimitMode              string             `json:"rateLimitMode,omitempty"`
	RateLimitSrcMask           int                `json:"rateLimitSrcMask,omitempty"`
	RateLimit                  string             `json:"rateLimit,omitempty"`
	Rules                      []string           `json:"rules,omitempty"`
	SecurityLogProfiles        []string           `json:"securityLogProfiles,omitempty"`
	ServiceDownImmediateAction string             `json:"serviceDownImmediateAction,omitempty"`
	Source                     string             `json:"source,omitempty"`
	SourcePort                 string             `json:"sourcePort,omitempty"`
	SynCookieStatus            string             `json:"synCookieStatus,omitempty"`
	TranslateAddress           string             `json:"translateAddress,omitempty"`
	TranslatePort              string             `json:"translatePort,omitempty"`
	VlansDisabled              bool               `json:"vlansDisabled,omitempty"`
	SourceAddressTranslation   *LBVirtualSNAT     `json:"sourceAddressTranslation,omitempty"`
}

func (target *LBVirtual) Merge(source *LBVirtual, opts ...func(*mergo.Config)) error {
	return mergo.Merge(target, source, opts...)
}

type LBVirtualPoliciesRef struct {
	Items []LBVirtualPolicy `json:"items"`
}

type LBVirtualProfileRef struct {
	Items []LBVirtualProfile `json:"items"`
}

func (vip *LBVirtual) UnmarshalJSON(data []byte) error {
	// Strip out the Policies and Profiles Reference entries, converting them
	// so simple policies and profiles arrays.
	type Alias LBVirtual
	aux := &struct {
		PoliciesReference LBVirtualPoliciesRef `json:"policiesReference"`
		ProfilesReference LBVirtualProfileRef  `json:"profilesReference"`
		*Alias
	}{
		Alias: (*Alias)(vip),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if len(aux.Policies) > 0 && len(aux.PoliciesReference.Items) > 0 {
		return fmt.Errorf("error: policies and policiesReference cannot both be set for: %s", aux.Name)
	} else {
		if len(aux.PoliciesReference.Items) > 0 {
			aux.Policies = aux.PoliciesReference.Items
		}
	}

	if len(aux.Profiles) > 0 && len(aux.ProfilesReference.Items) > 0 {
		return fmt.Errorf("error: profiles and profilesReference cannot both be set for: %s", aux.Name)
	} else {
		if len(aux.ProfilesReference.Items) > 0 {
			aux.Profiles = aux.ProfilesReference.Items
		}
	}

	vip = (*LBVirtual)(aux.Alias)
	return nil
}

type LBVirtuals struct {
	Items []LBVirtual
}

type LBVirtualStatsDescription struct {
	Description string `json:"description"`
}

type LBVirtualStatsValue struct {
	Value int `json:"value"`
}

type LBVirtualStatsInnerEntries struct {
	Clientside_bitsIn             LBStatsValue              `json:"clientside.bitsIn"`
	Clientside_bitsOut            LBStatsValue              `json:"clientside.bitsOut"`
	Clientside_curConns           LBStatsValue              `json:"clientside.curConns"`
	Clientside_evictedConns       LBStatsValue              `json:"clientside.evictedConns"`
	Clientside_maxConns           LBStatsValue              `json:"clientside.maxConns"`
	Clientside_pktsIn             LBStatsValue              `json:"clientside.pktsIn"`
	Clientside_pktsOut            LBStatsValue              `json:"clientside.pktsOut"`
	Clientside_slowKilled         LBStatsValue              `json:"clientside.slowKilled"`
	Clientside_totConns           LBStatsValue              `json:"clientside.totConns"`
	CmpEnableMode                 LBVirtualStatsDescription `json:"cmpEnableMode"`
	CmpEnabled                    LBVirtualStatsDescription `json:"cmpEnabled"`
	CsMaxConnDur                  LBStatsValue              `json:"csMaxConnDur"`
	CsMeanConnDur                 LBStatsValue              `json:"csMeanConnDur"`
	CsMinConnDur                  LBStatsValue              `json:"csMinConnDur"`
	Destination                   LBVirtualStatsDescription `json:"destination"`
	Ephemeral_bitsIn              LBStatsValue              `json:"ephemeral.bitsIn"`
	Ephemeral_bitsOut             LBStatsValue              `json:"ephemeral.bitsOut"`
	Ephemeral_curConns            LBStatsValue              `json:"ephemeral.curConns"`
	Ephemeral_evictedConns        LBStatsValue              `json:"ephemeral.evictedConns"`
	Ephemeral_maxConns            LBStatsValue              `json:"ephemeral.maxConns"`
	Ephemeral_pktsIn              LBStatsValue              `json:"ephemeral.pktsIn"`
	Ephemeral_pktsOut             LBStatsValue              `json:"ephemeral.pktsOut"`
	Ephemeral_slowKilled          LBStatsValue              `json:"ephemeral.slowKilled"`
	Ephemeral_totConns            LBStatsValue              `json:"ephemeral.totConns"`
	FiveMinAvgUsageRatio          LBStatsValue              `json:"fiveMinAvgUsageRatio"`
	FiveSecAvgUsageRatio          LBStatsValue              `json:"fiveSecAvgUsageRatio"`
	TmName                        LBVirtualStatsDescription `json:"tmName"`
	OneMinAvgUsageRatio           LBStatsValue              `json:"oneMinAvgUsageRatio"`
	Status_availabilityState      LBVirtualStatsDescription `json:"status.availabilityState"`
	Status_enabledState           LBVirtualStatsDescription `json:"status.enabledState"`
	Status_statusReason           LBVirtualStatsDescription `json:"status.statusReason"`
	SyncookieStatus               LBVirtualStatsDescription `json:"syncookieStatus"`
	Syncookie_accepts             LBStatsValue              `json:"syncookie.accepts"`
	Syncookie_hwAccepts           LBStatsValue              `json:"syncookie.hwAccepts"`
	Syncookie_hwSyncookies        LBStatsValue              `json:"syncookie.hwSyncookies"`
	Syncookie_hwsyncookieInstance LBStatsValue              `json:"syncookie.hwsyncookieInstance"`
	Syncookie_rejects             LBStatsValue              `json:"syncookie.rejects"`
	Syncookie_swsyncookieInstance LBStatsValue              `json:"syncookie.swsyncookieInstance"`
	Syncookie_syncacheCurr        LBStatsValue              `json:"syncookie.syncacheCurr"`
	Syncookie_syncacheOver        LBStatsValue              `json:"syncookie.syncacheOver"`
	Syncookie_syncookies          LBStatsValue              `json:"syncookie.syncookies"`
	TotRequests                   LBStatsValue              `json:"totRequests"`
}

type LBVirtualStatsNestedStats struct {
	Kind     string                     `json:"kind"`
	SelfLink string                     `json:"selfLink"`
	Entries  LBVirtualStatsInnerEntries `json:"entries"`
}

type LBVirtualURLKey struct {
	NestedStats LBVirtualStatsNestedStats `json:"nestedStats"`
}
type LBVirtualStatsOuterEntries map[string]LBVirtualURLKey

type LBVirtualStats struct {
	Kind       string                     `json:"kind"`
	Generation int                        `json:"generation"`
	SelfLink   string                     `json:"selfLink"`
	Entries    LBVirtualStatsOuterEntries `json:"entries"`
}

func (f *Device) ShowVirtuals() (error, *LBVirtuals) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/virtual"
	res := LBVirtuals{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) ShowVirtual(vname string) (error, *LBVirtual) {

	vname = strings.Replace(vname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/virtual/" + vname + "?expandSubcollections=true"
	res := LBVirtual{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) ShowVirtualStats(vname string) (error, *LBObjectStats) {

	vname = strings.Replace(vname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/virtual/" + vname + "/stats"
	res := LBObjectStats{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}
}

func (f *Device) ShowAllVirtualStats() (error, *LBVirtualStats) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/virtual/stats"
	res := LBVirtualStats{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}
}

func (f *Device) AddVirtual(virt *json.RawMessage) (error, *LBVirtual) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/virtual"
	res := LBVirtual{}

	// post the request
	err, _ := f.sendRequest(u, POST, virt, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) UpdateVirtual(vname string, body *json.RawMessage) (error, *LBVirtual) {

	vname = strings.Replace(vname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/virtual/" + vname
	res := LBVirtual{}

	// put the request
	err, _ := f.sendRequest(u, PUT, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) PatchVirtual(vname string, patch *LBVirtual) (error, *LBVirtual) {
	vname = strings.Replace(vname, "/", "~", -1)
	url := fmt.Sprintf("%s://%s/mgmt/tm/ltm/virtual/%s", f.Proto, f.Hostname, vname)
	existing := &LBVirtual{}
	var err error

	// Unless we're overwriting, grab the original and merge the patch with
	// the existing record's data  so that existing settings aren't overwritten,
	// but instead added to.
	if f.MergeStrategy() >= mergo.AppendAdditive {
		err, existing = f.ShowVirtual(vname)
		if err != nil {
			return err, nil
		}

		// merge existing fields into patch so we don't lose settings
		patch.Merge(existing, f.MergeConfig())
	}

	// merge the patch with our existing resource settings so we can see if
	// the patch is already applied or not
	new := &LBVirtual{}
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

func (f *Device) DeleteVirtual(vname string) (error, *Response) {

	vname = strings.Replace(vname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/virtual/" + vname
	res := json.RawMessage{}

	err, resp := f.sendRequest(u, DELETE, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, resp
	}

}
