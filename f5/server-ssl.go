package f5

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/rabbitt/f5er/mergo"
)

type LBServerSsl struct {
	Name                         string   `json:"name,omitempty"`
	Partition                    string   `json:"partition,omitempty"`
	FullPath                     string   `json:"fullPath,omitempty"`
	Generation                   int      `json:"generation,omitempty"`
	SelfLink                     string   `json:"selfLink,omitempty"`
	AlertTimeout                 string   `json:"alertTimeout,omitempty"`
	AllowExpiredCrl              string   `json:"allowExpiredCrl,omitempty"`
	AppService                   string   `json:"appService,omitempty"`
	Authenticate                 string   `json:"authenticate,omitempty"`
	AuthenticateDepth            int      `json:"authenticateDepth,omitempty"`
	AuthenticateName             string   `json:"authenticateName,omitempty"`
	BypassOnClientCertFail       string   `json:"bypassOnClientCertFail,omitempty"`
	BypassOnHandshakeAlert       string   `json:"bypassOnHandshakeAlert,omitempty"`
	C3dCaCert                    string   `json:"c3dCaCert,omitempty"`
	C3dCaKey                     string   `json:"c3dCaKey,omitempty"`
	C3dCertExtensionCustomOids   []string `json:"c3dCertExtensionCustomOids,omitempty"`
	C3dCertExtensionIncludes     []string `json:"c3dCertExtensionIncludes,omitempty"`
	C3dCertLifespan              int      `json:"c3dCertLifespan,omitempty"`
	CaFile                       string   `json:"caFile,omitempty"`
	CacheSize                    int      `json:"cacheSize,omitempty"`
	CacheTimeout                 int      `json:"cacheTimeout,omitempty"`
	Cert                         string   `json:"cert,omitempty"`
	Chain                        string   `json:"chain,omitempty"`
	CipherGroup                  string   `json:"cipherGroup,omitempty"`
	Ciphers                      string   `json:"ciphers,omitempty"`
	CrlFile                      string   `json:"crlFile,omitempty"`
	DefaultsFrom                 string   `json:"defaultsFrom,omitempty"`
	Description                  string   `json:"description,omitempty"`
	ExpireCertResponseControl    string   `json:"expireCertResponseControl,omitempty"`
	GenericAlert                 string   `json:"genericAlert,omitempty"`
	HandshakeTimeout             string   `json:"handshakeTimeout,omitempty"`
	Key                          string   `json:"key,omitempty"`
	MaxActiveHandshakes          string   `json:"maxActiveHandshakes,omitempty"`
	ModSslMethods                string   `json:"modSslMethods,omitempty"`
	Mode                         string   `json:"mode,omitempty"`
	Ocsp                         string   `json:"ocsp,omitempty"`
	TmOptions                    []string `json:"tmOptions,omitempty"`
	PeerCertMode                 string   `json:"peerCertMode,omitempty"`
	ProxySsl                     string   `json:"proxySsl,omitempty"`
	ProxySslPassthrough          string   `json:"proxySslPassthrough,omitempty"`
	RenegotiatePeriod            string   `json:"renegotiatePeriod,omitempty"`
	RenegotiateSize              string   `json:"renegotiateSize,omitempty"`
	Renegotiation                string   `json:"renegotiation,omitempty"`
	RetainCertificate            string   `json:"retainCertificate,omitempty"`
	SecureRenegotiation          string   `json:"secureRenegotiation,omitempty"`
	ServerName                   string   `json:"serverName,omitempty"`
	SessionMirroring             string   `json:"sessionMirroring,omitempty"`
	SessionTicket                string   `json:"sessionTicket,omitempty"`
	SniDefault                   string   `json:"sniDefault,omitempty"`
	SniRequire                   string   `json:"sniRequire,omitempty"`
	SslC3d                       string   `json:"sslC3d,omitempty"`
	SslForwardProxy              string   `json:"sslForwardProxy,omitempty"`
	SslForwardProxyBypass        string   `json:"sslForwardProxyBypass,omitempty"`
	SslSignHash                  string   `json:"sslSignHash,omitempty"`
	StrictResume                 string   `json:"strictResume,omitempty"`
	UncleanShutdown              string   `json:"uncleanShutdown,omitempty"`
	UntrustedCertResponseControl string   `json:"untrustedCertResponseControl,omitempty"`
}

func (target *LBServerSsl) Merge(source *LBServerSsl, opts ...func(*mergo.Config)) (err error) {
	return mergo.Merge(target, source, opts...)
}

type LBServerSsls struct {
	Items []LBServerSsl `json:"items"`
}

func (f *Device) ShowServerSsls() (error, *LBServerSsls) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/profile/server-ssl"
	res := LBServerSsls{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) ShowServerSsl(sname string) (error, *LBServerSsl) {

	server := strings.Replace(sname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/profile/server-ssl/" + server
	res := LBServerSsl{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) AddServerSsl(body *json.RawMessage) (error, *LBServerSsl) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/profile/server-ssl"
	res := LBServerSsl{}

	// post the request
	err, _ := f.sendRequest(u, POST, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) UpdateServerSsl(sname string, body *json.RawMessage) (error, *LBServerSsl) {

	server := strings.Replace(sname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/profile/server-ssl/" + server
	res := LBServerSsl{}

	// put the request
	err, _ := f.sendRequest(u, PUT, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) PatchServerSsl(name string, patch *LBServerSsl) (error, *LBServerSsl) {
	name = strings.Replace(name, "/", "~", -1)
	url := fmt.Sprintf("%s://%s/mgmt/tm/ltm/profile/server-ssl/%s", f.Proto, f.Hostname, name)
	existing := &LBServerSsl{}
	var err error

	// Unless we're overwriting, grab the original and merge the patch with
	// the existing record's data  so that existing settings aren't overwritten,
	// but instead added to.
	if f.MergeStrategy() >= mergo.AppendAdditive {
		err, existing = f.ShowServerSsl(name)
		if err != nil {
			return err, nil
		}

		// merge existing fields into patch so we don't lose settings
		patch.Merge(existing, f.MergeConfig())
	}

	// merge the patch with our existing resource settings so we can see if
	// the patch is already applied or not
	new := &LBServerSsl{}
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

func (f *Device) DeleteServerSsl(sname string) (error, *Response) {

	server := strings.Replace(sname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/profile/server-ssl/" + server
	res := json.RawMessage{}

	err, resp := f.sendRequest(u, DELETE, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, resp
	}

}
