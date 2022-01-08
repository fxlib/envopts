# envopts

![Tests](https://github.com/fxlib/envopts/actions/workflows/tests.yml/badge.svg)

Provides a code generator to turn structs annotated for the popular [env](github.com/caarlos0/env) library into functional options. [Functional options](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis) are a common pattern in Go but a lot of boiler place is required to use them
when the values come from environmnet variables. This code generator aims to solve this problem.

Furthermore it also automatically takes into account any default values specified through the "envDefault" tags. And usefull documentation is generated
when a comment is specified for each struct field.

## example:
Given the following struct:
```Go
//go:generate go run github.com/fxlib/envopts -type=FooEnv

// FooEnv would describe an environment variable struct parsed using: github.com/caarlos0/env
type FooEnv struct {
	Home     string        `env:"HOME"`
	Hosts    []string      `env:"HOSTS" envSeparator:":"`
	Duration time.Duration `env:"DURATION"` // Duration of the timeout
	Foo, Dar []env.Options `env:"FOO"`
}
```
Running `go generate` will create the following code to facilitate the functional options:
```Go
// Option is a functional option to configure FooEnv
type Option func(*FooEnv)

// FromFooEnv takes fully configured FooEnv and returns it as an option. Can be used to parse environment
// variables manually and provide the result in places where an option argument is expected.
func FromFooEnv(v *FooEnv) Option {
	return func(c *FooEnv) { *c = *v }
}

// ParseEnv will parse environment variables into a slice of options. Any options for parsing the
// environment can be supplied, for example to parse under a prefix.
func ParseEnv(eo env.Options) (opts []Option, err error) {
	var o FooEnv
	opts = append(opts, FromFooEnv(&o))
	return opts, env.Parse(&o, eo)
}

// ApplyOptions will merge all options into the resulting FooEnv while also ensuring default values are
// always set.
func ApplyOptions(opts ...Option) (res FooEnv) {
	env.Parse(&res, env.Options{Environment: make(map[string]string)})
	for _, o := range opts {
		o(&res)
	}
	return
}

// WithHome configures FooEnv
func WithHome(v string) Option { return func(o *FooEnv) { o.Home = v } }

// WithHosts configures FooEnv
func WithHosts(v []string) Option { return func(o *FooEnv) { o.Hosts = v } }

// WithDuration configures: Duration of the timeout
func WithDuration(v time.Duration) Option { return func(o *FooEnv) { o.Duration = v } }

// WithFoo configures FooEnv
func WithFoo(v []env.Options) Option { return func(o *FooEnv) { o.Foo = v } }

// WithDar configures FooEnv
func WithDar(v []env.Options) Option { return func(o *FooEnv) { o.Dar = v } }

```

## backlog

- [ ] We could read the 'required' tag and error when calling ApplyOptions when this is not provided. But it
      would required to return an extra `err` value and required options should be passed as separate arguments
      anyway.
- [ ] Figure out if we need to take care of nested structs of env options
- [ ] Instead of depending on `goimport` being present to clean up unused or used imports it would be nice if
      we could do everything from our own binary
- [x] Add comments to the generated code so developers that read it can follow what is happening. Also for godoc
- [x] Clean up the codebase, proper error handing in walking
- [x] Write proper unit tests instead of lazy smoke tests that call the Go command
- [x] Write some more documentation to get the point of this project across. Tell about features:
  - Comment handling
  - EnvDefautl handling
