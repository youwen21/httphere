package conf

import (
	_ "embed"
)

//go:embed httphere.yaml
var embedConfig []byte
