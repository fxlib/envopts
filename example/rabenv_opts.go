// Code generated by "envopts -test.testlogfile=/var/folders/wx/kfpf_v7s2m5fvwnm0scmzdzm0000gn/T/go-build337994121/b001/testlog.txt -test.paniconexit0 -test.timeout=10m0s"; DO NOT EDIT.
package example

import (
	"github.com/caarlos0/env/v6"
)

type RabEnvOption func(*RabEnv)

func FromRabEnv(v *RabEnv) RabEnvOption {
	return func(c *RabEnv) { *c = *v }
}
func RabEnvOptionParseEnv(eo env.Options) (opts []RabEnvOption, err error) {
	var o RabEnv
	opts = append(opts, FromRabEnv(&o))
	return opts, env.Parse(&o, eo)
}
func ApplyRabEnvOptions(opts ...RabEnvOption) (res RabEnv) {
	env.Parse(&res, env.Options{Environment: make(map[string]string)})
	for _, o := range opts {
		o(&res)
	}
	return
}
func OnHomeRabEnvOption(v string) RabEnvOption { return func(o *RabEnv) { o.Home = v } }
