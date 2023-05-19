package conf

import (
	_ "embed"
)

//go:embed httphere.toml
var embedConfig []byte
