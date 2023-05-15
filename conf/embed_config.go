package conf

import (
	_ "embed"
)

//go:embed config.toml
var embedConfig []byte
