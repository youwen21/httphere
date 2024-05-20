package conf

import (
	"regexp"
	"strings"
)

type VueHistoryRouters []string

func (p VueHistoryRouters) IsContain(path string) bool {
	for _, v := range p {
		if strings.Contains(v, "\\") {
			ok, err := regexp.MatchString(v, path)
			if ok && err == nil {
				return true
			}
		} else {
			if v == path {
				return true
			}
		}
	}
	return false
}

type FastCGICfg struct {
	Proto   string `json:"proto" toml:"proto" yaml:"proto" mapstructure:"proto"`
	Address string `json:"address" toml:"address" yaml:"address" mapstructure:"address"`
	Root    string `json:"root" toml:"root" yaml:"root" mapstructure:"root"`
}

type TlsCfg struct {
	CertFile string `json:"cert_file" toml:"cert_file" yaml:"cert_file" mapstructure:"cert_file"`
	KeyFile  string `json:"key_file" toml:"key_file" yaml:"key_file" mapstructure:"key_file"`
}

type BaseCfg struct {
	ListenHost string `json:"listen_host"  toml:"listen_host" yaml:"listen_host" mapstructure:"listen_host"`
	ListenPort string `json:"listen_port" toml:"listen_port" yaml:"listen_port" mapstructure:"listen_port"`

	DumpRequest string `json:"dump_request" toml:"dump_request" yaml:"dump_request" mapstructure:"dump_request"`

	StaticServer   string            `json:"static_server" toml:"static_server" yaml:"static_server" mapstructure:"static_server"`
	StaticRoot     string            `json:"static_root" toml:"static_root" yaml:"static_root" mapstructure:"static_root"`
	HistoryRouters VueHistoryRouters `json:"history_routers" toml:"history_routers" yaml:"history_routers"  mapstructure:"history_routers"`
}

type HostCfg struct {
	Host        string            `json:"host" toml:"host" yaml:"host" mapstructure:"host"`
	ReverseType string            `json:"reverse_type" toml:"reverse_type" yaml:"reverse_type" mapstructure:"reverse_type"`
	Paths       map[string]string `json:"paths" toml:"paths" yaml:"paths" mapstructure:"paths"`
	Rewrite     map[string]string `json:"rewrite" toml:"rewrite" yaml:"rewrite" mapstructure:"rewrite"`

	FastCGI FastCGICfg `json:"fast_cgi" toml:"fast_cgi" yaml:"fast_cgi" mapstructure:"fast_cgi"`
}

type HereConf struct {
	Base  BaseCfg   `json:"base" toml:"base" yaml:"base" mapstructure:"base"`
	Tls   TlsCfg    `json:"tls" toml:"tls" yaml:"tls" mapstructure:"tls"`
	Hosts []HostCfg `json:"hosts" toml:"hosts" yaml:"hosts" mapstructure:"hosts"`
}

func (hc *HereConf) GetHostRewrite(host string) map[string]string {
	for _, v := range Here.Hosts {
		if v.Host == host {
			return v.Rewrite
		}
	}

	return nil
}
