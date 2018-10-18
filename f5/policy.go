package f5

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/rabbitt/f5er/mergo"
)

type LBPolicyConditions struct {
	Name            string   `json:"name,omitempty"`
	FullPath        string   `json:"fullPath,omitempty"`
	Generation      int      `json:"generation,omitempty"`
	All             bool     `json:"all,omitempty"`
	CaseInsensitive bool     `json:"caseInsensitive,omitempty"`
	ClientAccepted  bool     `json:"clientAccepted,omitempty"`
	Equals          bool     `json:"equals,omitempty"`
	External        bool     `json:"external,omitempty"`
	Host            bool     `json:"host,omitempty"`
	HttpHost        bool     `json:"httpHost,omitempty"`
	HttpMethod      bool     `json:"httpMethod,omitempty"`
	HttpUri         bool     `json:"httpUri,omitempty"`
	Index           int      `json:"index,omitempty"`
	Local           bool     `json:"local,omitempty"`
	Not             bool     `json:"not,omitempty"`
	Normalized      bool     `json:"normalized,omitempty"`
	Path            bool     `json:"path,omitempty"`
	Port            bool     `json:"port,omitempty"`
	Present         bool     `json:"present,omitempty"`
	Remote          bool     `json:"remote,omitempty"`
	Request         bool     `json:"request,omitempty"`
	StartsWith      bool     `json:"startsWith,omitempty"`
	Tcp             bool     `json:"tcp,omitempty"`
	Values          []string `json:"values,omitempty"`
}

type LBPolicyActions struct {
	Name           string `json:"name,omitempty"`
	FullPath       string `json:"fullPath,omitempty"`
	Generation     int    `json:"generation,omitempty"`
	Asm            bool   `json:"asm,omitempty"`
	ClientAccepted bool   `json:"clientAccepted,omitempty"`
	ClientSsl      bool   `json:"clientSsl,omitempty"`
	Code           int    `json:"code,omitempty"`
	Connection     bool   `json:"connection,omitempty"`
	Enable         bool   `json:"enable,omitempty"`
	Disable        bool   `json:"disable,omitempty"`
	ExpirySecs     int    `json:"expirySecs,omitempty"`
	Facility       string `json:"facility,omitempty"`
	Forward        bool   `json:"forward,omitempty"`
	HTTPHost       bool   `json:"httpHost,omitempty"`
	HTTPHeader     bool   `json:"httpHeader,omitempty"`
	HTTPReply      bool   `json:"httpReply,omitempty"`
	Insert         bool   `json:"insert,omitempty"`
	Length         int    `json:"length,omitempty"`
	Location       string `json:"location,omitempty"`
	Message        string `json:"message,omitempty"`
	Offset         int    `json:"offset,omitempty"`
	Pool           string `json:"pool,omitempty"`
	Policy         string `json:"policy,omitempty"`
	Priority       string `json:"priority,omitempty"`
	Port           int    `json:"port,omitempty"`
	Redirect       bool   `json:"redirect,omitempty"`
	Replace        bool   `json:"replace,omitempty"`
	Request        bool   `json:"request,omitempty"`
	Select         bool   `json:"select,omitempty"`
	Shutdown       bool   `json:"shutdown,omitempty"`
	Status         int    `json:"status,omitempty"`
	Timeout        int    `json:"timeout,omitempty"`
	TmName         string `json:"tmName,omitempty"`
	Value          string `json:"value,omitempty"`
	VlanId         int    `json:"vlanId,omitempty"`
	Write          bool   `json:"write,omitempty"`
}

type LBPolicyConditionsRef struct {
	Items []LBPolicyConditions `json:"items,omitempty"`
}

type LBPolicyActionsRef struct {
	Items []LBPolicyActions `json:"items,omitempty"`
}

type LBPolicyRules struct {
	Name       string               `json:"name,omitempty"`
	FullPath   string               `json:"fullPath,omitempty"`
	Generation int                  `json:"generation,omitempty"`
	Ordinal    int                  `json:"ordinal,omitempty"`
	Actions    []LBPolicyActions    `json:"actions,omitempty"`
	Conditions []LBPolicyConditions `json:"conditions,omitempty"`
}

func (rules *LBPolicyRules) UnmarshalJSON(data []byte) error {
	// Strip out the Policies and Profiles Reference entries, converting them
	// so simple policies and profiles arrays.
	type Alias LBPolicyRules
	aux := &struct {
		ActionsReference    LBPolicyActionsRef    `json:"actionsReference,omitempty"`
		ConditionsReference LBPolicyConditionsRef `json:"conditionsReference,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(rules),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if len(aux.Actions) > 0 && len(aux.ActionsReference.Items) > 0 {
		return fmt.Errorf("error: actions and actionsReference can not both be set for: %s", aux.Name)
	} else {
		if len(aux.ActionsReference.Items) > 0 {
			aux.Actions = aux.ActionsReference.Items
		}
	}

	if len(aux.Conditions) > 0 && len(aux.ConditionsReference.Items) > 0 {
		return fmt.Errorf("error: conditions and conditionsReference can not both be set for: %s", aux.Name)
	} else {
		if len(aux.ConditionsReference.Items) > 0 {
			aux.Conditions = aux.ConditionsReference.Items
		}
	}

	rules = (*LBPolicyRules)(aux.Alias)
	return nil
}

type LBPolicyRulesRef struct {
	Items []LBPolicyRules `json:"items,omitempty"`
}

type LBPolicy struct {
	Name       string          `json:"name,omitempty"`
	Partition  string          `json:"partition,omitempty"`
	FullPath   string          `json:"fullPath,omitempty"`
	Generation int             `json:"generation,omitempty"`
	Controls   []string        `json:"controls,omitempty"`
	Requires   []string        `json:"requires,omitempty"`
	Strategy   string          `json:"strategy,omitempty"`
	Rules      []LBPolicyRules `json:"rules,omitempty"`
}

func (target *LBPolicy) Merge(source *LBPolicy, opts ...func(*mergo.Config)) (err error) {
	return mergo.Merge(target, source, opts...)
}

func (policy *LBPolicy) UnmarshalJSON(data []byte) error {
	// Strip out the Policies and Profiles Reference entries, converting them
	// so simple policies and profiles arrays.
	type Alias LBPolicy
	aux := &struct {
		RulesReference LBPolicyRulesRef `json:"rulesReference,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(policy),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if len(aux.Rules) > 0 && len(aux.RulesReference.Items) > 0 {
		return fmt.Errorf("error: rules and rulesReference can not both be set for: %s", aux.Name)
	} else {
		if len(aux.RulesReference.Items) > 0 {
			aux.Rules = aux.RulesReference.Items
		}
	}

	policy = (*LBPolicy)(aux.Alias)
	return nil
}

type LBPolicies struct {
	Items []LBPolicy `json:"items,omitempty"`
}

func (f *Device) ShowPolicies() (error, *LBPolicies) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/policy"
	res := LBPolicies{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) ShowPolicy(pname string) (error, *LBPolicy) {

	policy := strings.Replace(pname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/policy/" + policy + "?expandSubcollections=true"
	res := LBPolicy{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) AddPolicy(body *json.RawMessage) (error, *LBPolicy) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/policy"
	res := LBPolicy{}

	// post the request
	err, _ := f.sendRequest(u, POST, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) UpdatePolicy(pname string, body *json.RawMessage) (error, *LBPolicy) {

	policy := strings.Replace(pname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/policy/" + policy
	res := LBPolicy{}

	// put the request
	err, _ := f.sendRequest(u, PUT, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}
}

func (f *Device) PatchPolicy(name string, patch *LBPolicy) (error, *LBPolicy) {
	name = strings.Replace(name, "/", "~", -1)
	url := fmt.Sprintf("%s://%s/mgmt/tm/ltm/policy/%s", f.Proto, f.Hostname, name)
	existing := &LBPolicy{}
	var err error

	// Unless we're overwriting, grab the original and merge the patch with
	// the existing record's data  so that existing settings aren't overwritten,
	// but instead added to.
	if f.MergeStrategy() >= mergo.AppendAdditive {
		err, existing = f.ShowPolicy(name)
		if err != nil {
			return err, nil
		}

		// merge existing fields into patch so we don't lose settings
		patch.Merge(existing, f.MergeConfig())
	}

	// merge the patch with our existing resource settings so we can see if
	// the patch is already applied or not
	new := &LBPolicy{}
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

func (f *Device) DeletePolicy(pname string) (error, *Response) {

	//u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/policy/~" + partition + "~" + pname + "?expandSubcollections=true"
	policy := strings.Replace(pname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/policy/" + policy
	res := json.RawMessage{}

	err, resp := f.sendRequest(u, DELETE, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, resp
	}

}
