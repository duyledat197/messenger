package certs

import (
	"embed"
)

//go:embed *.pem
var CertFiles embed.FS
