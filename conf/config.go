package conf

type FastCGIConf struct {
	Proto   string `json:"proto" toml:"proto" yaml:"proto" mapstructure:"proto"`
	Address string `json:"address" toml:"address" yaml:"address" mapstructure:"address"`
	Root    string `json:"root" toml:"root" yaml:"root" mapstructure:"root"`
}

type HereConf struct {
	Base  BaseConf   `json:"base" toml:"base" yaml:"base" mapstructure:"base"`
	Hosts []HostConf `json:"hosts" toml:"hosts" yaml:"hosts" mapstructure:"hosts"`
}

type BaseConf struct {
	ListenHost string `json:"listen_host"  toml:"listen_host" yaml:"listen_host" mapstructure:"listen_host"`
	ListenPort string `json:"listen_port" toml:"listen_port" yaml:"listen_port" mapstructure:"listen_port"`

	StaticServer string `json:"static_server" toml:"static_server" yaml:"static_server" mapstructure:"static_server"`
	StaticRoot   string `json:"static_root" toml:"static_root" yaml:"static_root" mapstructure:"static_root"`

	DumpRequest string `json:"dump_request" toml:"dump_request" yaml:"dump_request" mapstructure:"dump_request"`
}

type HostConf struct {
	Host        string            `json:"host" toml:"host" yaml:"host" mapstructure:"host"`
	ReverseType string            `json:"reverse_type" toml:"reverse_type" yaml:"reverse_type" mapstructure:"reverse_type"`
	Paths       map[string]string `json:"paths" toml:"paths" yaml:"paths" mapstructure:"paths"`
	Rewrite     map[string]string `json:"rewrite" toml:"rewrite" yaml:"rewrite" mapstructure:"rewrite"`

	FastCGI FastCGIConf `json:"fast_cgi" toml:"fast_cgi" yaml:"fast_cgi" mapstructure:"fast_cgi"`
}
