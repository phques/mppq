// Package go_provider is an autogenerated binder stub for package provider.
//   gobind -lang=go github.com/phques/mppq/examples/libprovider/provider
//
// File is generated by gobind. Do not edit.
package go_provider

import (
	"github.com/phques/mppq/examples/libprovider/provider"
	"golang.org/x/mobile/bind/seq"
)

func proxy_InitAppFilesDir(out, in *seq.Buffer) {
	param_appFilesDir_ := in.ReadUTF16()
	provider.InitAppFilesDir(param_appFilesDir_)
}

func proxy_Register(out, in *seq.Buffer) {
	param_serviceName := in.ReadUTF16()
	provider.Register(param_serviceName)
}

func proxy_Start(out, in *seq.Buffer) {
	provider.Start()
}

func init() {
	seq.Register("provider", 1, proxy_InitAppFilesDir)
	seq.Register("provider", 2, proxy_Register)
	seq.Register("provider", 3, proxy_Start)
}
