package conf

import (
	"regexp"
	"strings"
)

type FastCGIConf struct {
	Proto   string `json:"proto" toml:"proto" yaml:"proto" mapstructure:"proto"`
	Address string `json:"address" toml:"address" yaml:"address" mapstructure:"address"`
	Root    string `json:"root" toml:"root" yaml:"root" mapstructure:"root"`
}

type HereConf struct {
	Base  BaseConf   `json:"base" toml:"base" yaml:"base" mapstructure:"base"`
	Tls   TlsConf    `json:"tls" toml:"tls" yaml:"tls" mapstructure:"tls"`
	Hosts []HostConf `json:"hosts" toml:"hosts" yaml:"hosts" mapstructure:"hosts"`
}

type TlsConf struct {
	CertFile string `json:"cert_file" toml:"cert_file" yaml:"cert_file" mapstructure:"cert_file"`
	KeyFile  string `json:"key_file" toml:"key_file" yaml:"key_file" mapstructure:"key_file"`
}

type BaseConf struct {
	ListenHost string `json:"listen_host"  toml:"listen_host" yaml:"listen_host" mapstructure:"listen_host"`
	ListenPort string `json:"listen_port" toml:"listen_port" yaml:"listen_port" mapstructure:"listen_port"`

	StaticServer string `json:"static_server" toml:"static_server" yaml:"static_server" mapstructure:"static_server"`
	StaticRoot   string `json:"static_root" toml:"static_root" yaml:"static_root" mapstructure:"static_root"`

	DumpRequest string `json:"dump_request" toml:"dump_request" yaml:"dump_request" mapstructure:"dump_request"`
}

type HistoryRouters []string

func (p HistoryRouters) IsContain(path string) bool {
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

type HostConf struct {
	Host        string            `json:"host" toml:"host" yaml:"host" mapstructure:"host"`
	ReverseType string            `json:"reverse_type" toml:"reverse_type" yaml:"reverse_type" mapstructure:"reverse_type"`
	Paths       map[string]string `json:"paths" toml:"paths" yaml:"paths" mapstructure:"paths"`
	Rewrite     map[string]string `json:"rewrite" toml:"rewrite" yaml:"rewrite" mapstructure:"rewrite"`

	FastCGI FastCGIConf `json:"fast_cgi" toml:"fast_cgi" yaml:"fast_cgi" mapstructure:"fast_cgi"`
}

func (hc *HereConf) GetHostRewrite(host string) map[string]string {
	for _, v := range Here.Hosts {
		if v.Host == host {
			return v.Rewrite
		}
	}

	return nil
}
