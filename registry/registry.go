package registry

import (
	"github.com/bayupermadi/dbcheck"
)

var (
	dialerRegistry = make(map[string]dbcheck.Dialer)
)

func Register(name string, dialer dbcheck.Dialer) {
	dialerRegistry[name] = dialer
}

func Dialers(db string) dbcheck.Dialer {
	return dialerRegistry[db]
}
