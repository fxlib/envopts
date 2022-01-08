// Code generated by "envopts -test.testlogfile=/var/folders/wx/kfpf_v7s2m5fvwnm0scmzdzm0000gn/T/go-build337994121/b001/testlog.txt -test.paniconexit0 -test.timeout=10m0s"; DO NOT EDIT.
package example

import (
	"time"

	"github.com/caarlos0/env/v6"
)

type FooOption func(*fooEnv)

func fromFooEnv(v *fooEnv) FooOption {
	return func(c *fooEnv) { *c = *v }
}
func FooOptionParseEnv(eo env.Options) (opts []FooOption, err error) {
	var o fooEnv
	opts = append(opts, fromFooEnv(&o))
	return opts, env.Parse(&o, eo)
}
func applyFooOptions(opts ...FooOption) (res fooEnv) {
	env.Parse(&res, env.Options{Environment: make(map[string]string)})
	for _, o := range opts {
		o(&res)
	}
	return
}
func WithHostsFooOption(v []string) FooOption         { return func(o *fooEnv) { o.Hosts = v } }
func WithDurationFooOption(v time.Duration) FooOption { return func(o *fooEnv) { o.Duration = v } }
