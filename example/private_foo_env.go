// Code generated by "stringer -type=fooEnv -optionType=PrivOption -output=private_foo_env.go"; DO NOT EDIT.
package example

import (
	"time"

	"github.com/caarlos0/env/v6"
)

type PrivOption func(*fooEnv)

func WithHostsPrivOption(v []string) PrivOption         { return func(o *fooEnv) { o.Hosts = v } }
func WithDurationPrivOption(v time.Duration) PrivOption { return func(o *fooEnv) { o.Duration = v } }
func fromfooEnv(v *fooEnv) PrivOption {
	return func(c *fooEnv) { *c = *v }
}
func PrivOptionParseEnv(eo env.Options) (opts []PrivOption, err error) {
	var o fooEnv
	opts = append(opts, fromfooEnv(&o))
	return opts, env.Parse(&o, eo)
}
func applyPrivOptions(opts ...PrivOption) (res fooEnv) {
	env.Parse(&res, env.Options{Environment: make(map[string]string)})
	for _, o := range opts {
		o(&res)
	}
	return
}
