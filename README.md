# envopts
![Tests](https://github.com/fxlib/envopts/actions/workflows/tests.yml/badge.svg)

Provides a code generator for turning [env](github.com/caarlos0/env) structure into functional options: https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis

```
// FooEnv would describe an environment variable struct parsed using: github.com/caarlos0/env
type FooEnv struct {
	Home         string        `env:"HOME"`
	Port         int           `env:"PORT" envDefault:"3000"`
	Password     string        `env:"PASSWORD,unset"`
	IsProduction bool          `env:"PRODUCTION"`
	Hosts        []string      `env:"HOSTS" envSeparator:":"`
	Duration     time.Duration `env:"DURATION"`
	TempFolder   string        `env:"TEMP_FOLDER" envDefault:"${HOME}/tmp" envExpand:"true"`
}

//go:generate go run github.com/fxlib/envopts -type=FooEnv
// The generate above will generate
```
