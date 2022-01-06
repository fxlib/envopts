// Code generated by "stringer -type=FooEnv"; DO NOT EDIT.
package example

import (
	"time"

	"github.com/caarlos0/env/v6"
)

type Option func(*FooEnv)

func WithHome(v string) Option            { return func(o *FooEnv) { o.Home = v } }
func WithPort(v int) Option               { return func(o *FooEnv) { o.Port = v } }
func WithPassword(v string) Option        { return func(o *FooEnv) { o.Password = v } }
func WithIsProduction(v bool) Option      { return func(o *FooEnv) { o.IsProduction = v } }
func WithHosts(v []string) Option         { return func(o *FooEnv) { o.Hosts = v } }
func WithDuration(v time.Duration) Option { return func(o *FooEnv) { o.Duration = v } }
func WithTempFolder(v string) Option      { return func(o *FooEnv) { o.TempFolder = v } }
func FromFooEnv(v *FooEnv) Option {
	return func(c *FooEnv) { *c = *v }
}
func ParseEnv(eo env.Options) (opts []Option, err error) {
	var o FooEnv
	opts = append(opts, FromFooEnv(&o))
	return opts, env.Parse(&o, eo)
}
func ApplyOptions(opts ...Option) (res FooEnv) {
	env.Parse(&res, env.Options{Environment: make(map[string]string)})
	for _, o := range opts {
		o(&res)
	}
	return
}
