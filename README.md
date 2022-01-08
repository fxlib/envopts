# envopts

![Tests](https://github.com/fxlib/envopts/actions/workflows/tests.yml/badge.svg)

Provides a code generator for turning [env](github.com/caarlos0/env) structure into functional options: https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis. Check out the examples directory for general usage.

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
// The generate above will generate functional options for each exported struct member with the 'env' tag
```

## backlog

- [ ] Write some more documentation to get the point of this project across. Tell about features:
  - Comment handling
  - EnvDefautl handling
- [ ] We could read the 'required' tag and error when calling ApplyOptions when this is not provided. But it
      would required to return an extra `err` value and required options should be passed as separate arguments
      anyway.
- [ ] Figure out if we need to take care of nested structs of env options
- [ ] Instead of depending on `goimport` being present to clean up unused or used imports it would be nice if
      we could do everything from our own binary
- [x] Add comments to the generated code so developers that read it can follow what is happening. Also for godoc
- [x] Clean up the codebase, proper error handing in walking
- [x] Write proper unit tests instead of lazy smoke tests that call the Go command
