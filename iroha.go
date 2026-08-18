// Package iroha registers "k6/x/iroha", a stateless k6 module bridging
// onto irohasign for signing Iroha transactions from k6 scripts. Scripts
// supply key material and data directly; goja can't build or sign Iroha
// protobuf itself, so this work happens in Go, in-process. Signed bytes
// are returned manually base64-encoded, since a raw []byte returned to
// goja doesn't auto-marshal like encoding/json does.
package iroha

import (
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/iroha", New())
}

// RootModule implements modules.Module.
type RootModule struct{}

// New returns a new RootModule — the entry point xk6/k6 calls to register
// this extension.
func New() *RootModule {
	return &RootModule{}
}

// NewModuleInstance implements modules.Module. k6 calls this once per VU.
func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &Instance{vu: vu}
}

// Instance implements modules.Instance. Stateless — holds nothing beyond
// the VU reference k6's module lifecycle requires.
type Instance struct {
	vu modules.VU
}

// Exports implements modules.Instance.
func (i *Instance) Exports() modules.Exports {
	return modules.Exports{
		Named: map[string]any{
			"signTransfer":        i.SignTransfer,
			"signGrantPermission": i.SignGrantPermission,
			"reSignWithBatchMeta": i.ReSignWithBatchMeta,
		},
	}
}
